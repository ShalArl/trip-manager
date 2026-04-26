package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/container"
	"github.com/ShalArl/trip-manager/internal/database"
	"github.com/ShalArl/trip-manager/internal/storage"
	"github.com/jmoiron/sqlx"
)

// App holds shared application dependencies
type App struct {
	SQLDb           *sqlx.DB
	FirestoreClient *firestore.Client
	Logger          *log.Logger
	Services        *container.ServiceContainer
	Config          *config.Config
}

// New initializes and returns a new App instance
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	// Initialize sql database with config database URL
	sqlDb, err := database.ConnectSql(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations automatically
	if err := database.RunEmbeddedMigrations(sqlDb); err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to run migrations: %w", err),
			sqlDb.Close(),
		)
	}

	// Initialize firestore database with config Firestore URL
	firestoreClient, err := database.ConnectFirestore(ctx, cfg.FirebaseConfig)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to connect to firestore: %w", err),
			sqlDb.Close())
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[trip-manager] ", log.LstdFlags|log.Lshortfile)
	stor, err := storage.NewFromEnv(ctx, cfg.Storage)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to setup storage: %w", err),
			sqlDb.Close(), firestoreClient.Close())
	}
	svcs, err := container.NewServiceContainer(&container.ServiceConfig{
		SQLDb:           sqlDb,
		FirestoreClient: firestoreClient,
		Logger:          logger,
		Config:          cfg,
		Storage:         stor,
	})

	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to setup service container: %w", err), sqlDb.Close())
	}

	app := &App{
		SQLDb:    sqlDb,
		Services: svcs,
		Logger:   logger,
		Config:   cfg,
	}

	logger.Println("Application initialized successfully")
	return app, nil
}

// Close closes the database connection
func (a *App) Close() error {
	if a.SQLDb != nil {
		return a.SQLDb.Close()
	}
	return nil
}
