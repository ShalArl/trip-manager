package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface using local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// Upload saves a file to local storage and returns the URL
func (ls *LocalStorage) Upload(ctx context.Context, fileName string, file io.Reader) (string, error) {
	// Create subdirectories if needed
	filePath := filepath.Join(ls.basePath, fileName)
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("failed to close file: ", err)
			return
		}
	}(f)

	// Copy file content
	if _, err := io.Copy(f, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return the URL
	return fmt.Sprintf("/uploads/%s", fileName), nil
}

// ReadFile retrieves a file from local storage
func (ls *LocalStorage) ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	filePath := filepath.Join(ls.basePath, fileName)

	// Security check: prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	basePath := filepath.Clean(ls.basePath)
	if !filepath.HasPrefix(cleanPath, basePath) {
		return nil, fmt.Errorf("invalid file path")
	}

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}

// GetUrl returns the public URL for a file
func (ls *LocalStorage) GetUrl(ctx context.Context, fileName string) (string, error) {
	// For local storage in development, return a path that can be served
	// In production, this would be replaced with actual S3/GCS URLs
	return fmt.Sprintf("/uploads/%s", fileName), nil
}
