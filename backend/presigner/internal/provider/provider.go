package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ShalArl/trip-manager/backend/presigner/config"
)

func NewFromEnv(ctx context.Context, cfg config.Config) (Storage, error) {
	switch strings.ToLower(cfg.Type) {
	case "gcs":
		return NewGCStorage(ctx, GCSConfig{
			Bucket:   cfg.Bucket,
			SignerSA: cfg.SignerSA,
			TTL:      cfg.TTL,
		})
	case "s3":
		return NewS3Storage(S3Config{
			Bucket:    cfg.Bucket,
			Region:    cfg.Region,
			Endpoint:  cfg.Endpoint,
			PublicURL: cfg.PublicURL,
			AccessKey: cfg.AccessKey,
			SecretKey: cfg.SecretKey,
			TTL:       cfg.TTL,
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
