package storage

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/config"
)

func NewFromEnv(ctx context.Context, cfg config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "gcs":
		return NewGCStorage(ctx, GCSConfig{
			Bucket:   cfg.GCS.Bucket,
			SignerSA: cfg.GCS.SignerSA,
			TTL:      cfg.SignedURLTTL,
		})
	case "s3":
		return NewS3Storage(S3Config{
			Bucket:    cfg.S3.Bucket,
			Region:    cfg.S3.Region,
			Endpoint:  cfg.S3.Endpoint,
			PublicURL: cfg.S3.PublicURL,
			AccessKey: cfg.S3.AccessKey,
			SecretKey: cfg.S3.SecretKey,
			TTL:       cfg.SignedURLTTL,
		})
	default:
		return nil, fmt.Errorf("unknown storage type: %q", cfg.Type)
	}
}
