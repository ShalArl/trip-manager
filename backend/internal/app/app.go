package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/container"
	"github.com/ShalArl/trip-manager/internal/database"
	"github.com/ShalArl/trip-manager/internal/storage"
	"github.com/jmoiron/sqlx"
)

// App holds shared application dependencies
type App struct {
	DB       *sqlx.DB
	Logger   *log.Logger
	Services *container.ServiceContainer
	Config   *config.Config
}

// New initializes and returns a new App instance
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	// Initialize database with config database URL
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations automatically
	if err := database.RunEmbeddedMigrations(db); err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to run migrations: %w", err),
			db.Close(),
		)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[trip-manager] ", log.LstdFlags|log.Lshortfile)
	stor, err := storage.NewFromEnv(ctx, cfg.Storage)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to setup storage: %w", err),
			db.Close())
	}
	svcs, err := container.NewServiceContainer(&container.ServiceConfig{
		DB:      db,
		Logger:  logger,
		Config:  cfg,
		Storage: stor,
	})

	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to setup service container: %w", err), db.Close())
	}

	app := &App{
		DB:       db,
		Services: svcs,
		Logger:   logger,
		Config:   cfg,
	}

	logger.Println("Application initialized successfully")
	return app, nil
}

// Close closes the database connection
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
