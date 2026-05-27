package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/presigner/config"
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

// PresignProvider defines the interface for presigning URLs
type PresignProvider interface {
	// GetUploadURL returns a presigned URL for uploading
	GetUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)

	// GetDownloadURL returns a presigned URL for downloading
	GetDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)

	// Close closes any resources used by the provider
	Close() error
}
