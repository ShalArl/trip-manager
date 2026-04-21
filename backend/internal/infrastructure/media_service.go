package infrastructure

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/ShalArl/trip-manager/internal/storage"
)

type MediaType string

const (
	MediaTypeAvatar   MediaType = "avatars"
	MediaTypeTrip     MediaType = "trips"
	MediaTypeLocation MediaType = "locations"
	MediaTypeActivity MediaType = "activity"
)

func ParseMediaType(s string) (MediaType, error) {
	switch s {
	case "avatar":
		return MediaTypeAvatar, nil
	case "trip":
		return MediaTypeTrip, nil
	case "location":
		return MediaTypeLocation, nil
	case "activity":
		return MediaTypeActivity, nil
	default:
		return "", fmt.Errorf("invalid media type: %q", s)
	}
}

type MediaService struct {
	storage storage.Storage
}

func NewMediaService(stor storage.Storage) *MediaService {
	return &MediaService{storage: stor}
}

type UploadTicket struct {
	UploadURL string // Short-lived PUT-URL
	Key       string // Object-Key for DB
}

// PrepareUpload generates Object-Key + Upload-URL.
func (ms *MediaService) PrepareUpload(ctx context.Context, userID string, mediaType MediaType, fileName string) (UploadTicket, error) {
	key := buildKey(mediaType, userID, fileName)

	url, err := ms.storage.GetUploadURL(ctx, key)
	if err != nil {
		return UploadTicket{}, fmt.Errorf("upload url: %w", err)
	}

	return UploadTicket{UploadURL: url, Key: key}, nil
}

// GetDownloadURL returns short-lived GET-URL.
func (ms *MediaService) GetDownloadURL(ctx context.Context, key string) (string, error) {
	return ms.storage.GetDownloadURL(ctx, key)
}

// ConfirmUpload verifies successful upload.
func (ms *MediaService) ConfirmUpload(ctx context.Context, key string) (bool, error) {
	return ms.storage.Exists(ctx, key)
}

// Delete removes an object.
func (ms *MediaService) Delete(ctx context.Context, key string) error {
	return ms.storage.Delete(ctx, key)
}

// buildKey collision free object-key.
func buildKey(mType MediaType, userID string, originalFileName string) string {
	ext := getFileExtension(originalFileName)

	if mType == MediaTypeAvatar {
		// Avatar pro User ist singulär — User-ID als Key reicht.
		return fmt.Sprintf("avatars/%s%s", userID, ext)
	}

	// Alle anderen Media-Types: UUID pro Upload, verhindert Kollisionen.
	return fmt.Sprintf("%s/%s/%s%s", mType, userID, uuid.NewString(), ext)
}

func getFileExtension(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".heic":
		return ext
	default:
		return ".jpg"
	}
}
