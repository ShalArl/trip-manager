package app

import (
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
	Config   *config.Config
	Services *container.ServiceContainer
	Storage  storage.Storage
}

// New initializes and returns a new App instance
func New(cfg *config.Config) (*App, error) {
	// Initialize database with config database URL
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations automatically
	if err := database.RunEmbeddedMigrations(db); err != nil {
		if errors.Is(err, db.Close()) {
			return nil, fmt.Errorf("failed to run migrations - FATAL could not close DB connection: %w", err)
		}
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[trip-manager] ", log.LstdFlags|log.Lshortfile)

	// Initialize storage (local or S3)
	stor, err := setupStorage(cfg, logger)
	if err != nil {
		logger.Printf("Warning: Failed to initialize storage: %v", err)
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	svcs := container.NewServiceContainer(&container.ServiceConfig{
		DB:      db,
		Logger:  logger,
		Config:  cfg,
		Storage: stor,
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
		Storage:  stor,
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

// setupStorage initializes storage based on configuration
// Supports both local storage (development) and S3 (MinIO/AWS)
func setupStorage(cfg *config.Config, logger *log.Logger) (storage.Storage, error) {
	if cfg.StorageType == "s3" {
		// Use S3 storage (MinIO in development, AWS S3 in production)
		s3Storage, err := storage.NewS3Storage(storage.S3Config{
			Bucket:    cfg.S3Bucket,
			Region:    cfg.S3Region,
			Endpoint:  cfg.S3Endpoint,
			PublicURL: cfg.S3PublicURL,
			AccessKey: cfg.S3AccessKey,
			SecretKey: cfg.S3SecretKey,
			UseSSL:    cfg.S3UseSSL,
		})
		if err != nil {
			return nil, err
		}
		logger.Println("✅ S3 Storage initialized (MinIO or AWS S3)")
		return s3Storage, nil
	}

	// Default to local storage
	localStorage, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		return nil, err
	}
	logger.Println("✅ Local Storage initialized")
	return localStorage, nil
}
