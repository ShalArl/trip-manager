package provider

import (
	"context"
	"io"
)

type Storage interface {
	// GetUploadURL Generate a signed URL for uploading a file. fileName is the Object-Key in the storage.
	GetUploadURL(ctx context.Context, fileName string) (string, error)
	// GetDownloadURL Generate a signed URL for downloading a file. fileName is the Object-Key in the storage.
	GetDownloadURL(ctx context.Context, fileName string) (string, error)
	// ReadFile INTERNAL
	ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error)
	// WriteFile Optional
	WriteFile(ctx context.Context, fileName string, file io.Reader) error
	// Delete Optional: Delete
	Delete(ctx context.Context, fileName string) error
	// Exists Optional: Check if file exists
	Exists(ctx context.Context, fileName string) (bool, error)
}
