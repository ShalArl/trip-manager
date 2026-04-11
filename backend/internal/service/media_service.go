package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ShalArl/trip-manager/internal/storage"
)

// MediaService handles file operations like avatar uploads
type MediaService struct {
	storage storage.Storage
}

// NewMediaService creates a new media service
func NewMediaService(stor storage.Storage) *MediaService {
	return &MediaService{
		storage: stor,
	}
}

// MediaType defines the type of media being uploaded
type MediaType string

const (
	MediaTypeAvatar    MediaType = "avatar"
	MediaTypeTrip      MediaType = "trip"
	MediaTypeLocation  MediaType = "location"
	MediaTypeActivity  MediaType = "activity"
)

// UploadImageOptions contains options for image uploads
type UploadImageOptions struct {
	MediaType MediaType
	UserID    string
	FileName  string
}

// UploadImage uploads an image file and returns the URL
func (ms *MediaService) UploadImage(ctx context.Context, file io.Reader, opts UploadImageOptions) (string, error) {
	if opts.UserID == "" {
		return "", fmt.Errorf("user ID is required")
	}

	// Validate media type
	validTypes := map[MediaType]bool{
		MediaTypeAvatar:    true,
		MediaTypeTrip:      true,
		MediaTypeLocation:  true,
		MediaTypeActivity:  true,
	}

	if !validTypes[opts.MediaType] {
		return "", fmt.Errorf("invalid media type: %s", opts.MediaType)
	}

	// Generate safe filename with media type and user ID
	ext := filepath.Ext(opts.FileName)
	if ext == "" {
		ext = ".jpg" // Default extension
	}

	var storagePath string
	switch opts.MediaType {
	case MediaTypeAvatar:
		// Avatar: one per user, overwrite previous
		storagePath = fmt.Sprintf("avatars/%s%s", opts.UserID, ext)
	default:
		// Other types: use original filename with media type and user ID folder
		safeFileName := sanitizeFileName(opts.FileName)
		storagePath = fmt.Sprintf("%s/%s/%s", opts.MediaType, opts.UserID, safeFileName)
	}

	// Upload to storage
	fileURL, err := ms.storage.Upload(ctx, storagePath, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fileURL, nil
}

// DeleteImage deletes an image file
// Note: This method is reserved for future use
func (ms *MediaService) DeleteImage(_ context.Context, _ string) error {
	// Note: This would require adding Delete method to Storage interface
	// For now, we'll just return success
	// Storage implementation can handle cleanup
	return nil
}

// sanitizeFileName removes potentially dangerous characters from filename
func sanitizeFileName(filename string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(filename)
}


