package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/internal/storage"
)

// PresignedURLOptions contains options for generating presigned URLs
type PresignedURLOptions struct {
	MediaType MediaType
	UserID    string
	FileName  string
	ExpiresIn time.Duration // URL expiration time
}

// GeneratePresignedURL generates a presigned URL for direct file upload to S3/MinIO
func (ms *MediaService) GeneratePresignedURL(ctx context.Context, opts PresignedURLOptions) (string, error) {
	if opts.UserID == "" {
		return "", fmt.Errorf("user ID is required")
	}

	if opts.FileName == "" {
		return "", fmt.Errorf("file name is required")
	}

	if opts.ExpiresIn == 0 {
		opts.ExpiresIn = 15 * time.Minute // Default: 15 minutes
	}

	// Generate storage path (same logic as UploadImage)
	storagePath, err := ms.generateStoragePath(opts.MediaType, opts.UserID, opts.FileName)
	if err != nil {
		return "", err
	}

	// Check if storage supports presigned URLs
	s3Storage, ok := ms.storage.(*storage.S3Storage)
	if !ok {
		return "", fmt.Errorf("presigned URLs are only supported for S3 storage")
	}

	// Generate presigned URL
	presignedURL, err := s3Storage.GeneratePresignedURL(ctx, storagePath, opts.ExpiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL, nil
}

// generateStoragePath generates the storage path for a file (DRY - same as UploadImage)
func (ms *MediaService) generateStoragePath(mediaType MediaType, userID string, fileName string) (string, error) {
	validTypes := map[MediaType]bool{
		MediaTypeAvatar:   true,
		MediaTypeTrip:     true,
		MediaTypeLocation: true,
		MediaTypeActivity: true,
	}

	if !validTypes[mediaType] {
		return "", fmt.Errorf("invalid media type: %s", mediaType)
	}

	var storagePath string
	switch mediaType {
	case MediaTypeAvatar:
		// Avatar: one per user, overwrite previous
		ext := getFileExtension(fileName)
		storagePath = fmt.Sprintf("avatars/%s%s", userID, ext)
	default:
		// Other types: use filename with media type and user ID folder
		safeFileName := sanitizeFileName(fileName)
		storagePath = fmt.Sprintf("%s/%s/%s", mediaType, userID, safeFileName)
	}

	return storagePath, nil
}

// getFileExtension extracts the file extension from a filename
func getFileExtension(fileName string) string {
	for i := len(fileName) - 1; i >= 0 && fileName[i] != '/'; i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ".jpg" // Default
}

