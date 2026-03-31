package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/container"
	"github.com/ShalArl/trip-manager/internal/database"
	"github.com/jmoiron/sqlx"
)

// App holds shared application dependencies
type App struct {
	DB       *sqlx.DB
	Logger   *log.Logger
	Config   *config.Config
	Services *container.ServiceContainer
}

// New initializes and returns a new App instance
func New(cfg *config.Config) (*App, error) {
	// Initialize database with config database URL
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations automatically
	if err := database.RunMigrations(db); err != nil {
		if errors.Is(err, db.Close()) {
			return nil, fmt.Errorf("failed to run migrations - FATAL could not close DB connection: %w", err)
		}
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[trip-manager] ", log.LstdFlags|log.Lshortfile)

	svcs := container.NewServiceContainer(&container.ServiceConfig{
		DB:     db,
		Logger: logger,
		Config: cfg,
	})
	if svcs == nil {
		if errors.Is(err, db.Close()) {
			return nil, fmt.Errorf("failed to create service container - FATAL could not close DB connection: %w", err)
		} // Ensure DB is closed if service container creation fails
		return nil, fmt.Errorf("failed to create service container")
	}

	app := &App{
		DB:       db,
		Logger:   logger,
		Config:   cfg,
		Services: svcs,
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
