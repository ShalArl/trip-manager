package storage

import (
	"context"
	"io"
)

type Storage interface {
	// Upload saves a file and returns a URL
	Upload(ctx context.Context, fileName string, file io.Reader) (string, error)
	// ReadFile retrieves a file by name
	ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error)
	// GetUrl returns the public URL for a file
	GetUrl(ctx context.Context, fileName string) (string, error)
}
