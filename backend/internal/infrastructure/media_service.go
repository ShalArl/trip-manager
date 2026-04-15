package infrastructure

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/ShalArl/trip-manager/internal/storage"
)

type MediaType string

const (
	MediaTypeAvatar   MediaType = "avatars"
	MediaTypeTrip     MediaType = "trips"
	MediaTypeLocation MediaType = "locations"
	MediaTypeActivity MediaType = "activity"
)

type MediaService struct {
	storage storage.Storage
}

func NewMediaService(stor storage.Storage) *MediaService {
	return &MediaService{storage: stor}
}

func (ms *MediaService) UploadImage(ctx context.Context, file io.Reader, userID string, mediaType MediaType, fileName string) (string, error) {
	storagePath, err := ms.generateStoragePath(mediaType, userID, fileName)
	if err != nil {
		return "", err
	}

	_, err = ms.storage.Upload(ctx, storagePath, file)
	if err != nil {
		return "", err
	}

	return ms.storage.GetUrl(ctx, storagePath)
}

func (ms *MediaService) GeneratePresignedURL(ctx context.Context, userID string, mediaType MediaType, fileName string) (string, error) {
	storagePath, err := ms.generateStoragePath(mediaType, userID, fileName)
	if err != nil {
		return "", err
	}

	s3Storage, ok := ms.storage.(*storage.S3Storage)
	if !ok {
		return "", fmt.Errorf("storage does not support presigned URLs")
	}

	return s3Storage.GeneratePresignedURL(ctx, storagePath, 15*time.Minute)
}

func (ms *MediaService) generateStoragePath(mType MediaType, userID string, fileName string) (string, error) {
	// Nutze jetzt die Hilfsfunktion konsequent
	ext := getFileExtension(fileName)

	if mType == MediaTypeAvatar {
		// avatars/USER_ID.png
		return fmt.Sprintf("avatars/%s%s", userID, ext), nil
	}

	// trips/USER_ID/safe_name.png
	return fmt.Sprintf("%s/%s/%s", mType, userID, sanitizeFileName(fileName)), nil
}

func sanitizeFileName(filename string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filename)
}

func getFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return ".jpg"
	}
	return ext
}
