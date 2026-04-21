package database

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
)

// ConnectSql establishes a database connection using the provided database URL
func ConnectSql(ctx context.Context, databaseURL string) (*sqlx.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("postgres database URL is required")
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database via sqlx")
	return db, nil
}

// ConnectFirestore establishes a connection to Firestore
func ConnectFirestore(ctx context.Context, cfg config.FirestoreConfig) (*firestore.Client, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("firestore project id required")
	}

	var client *firestore.Client
	var err error

	if cfg.IsLocal {
		client, err = firestore.NewClient(ctx, cfg.ProjectID, option.WithoutAuthentication(), option.WithEndpoint(cfg.Endpoint))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to firestore: %w", err)
		}
	} else {
		client, err = firestore.NewClient(ctx, cfg.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to firestore: %w", err)
		}
	}

	log.Println("Successfully connected to the firestore via firestore")
	return client, nil
}
