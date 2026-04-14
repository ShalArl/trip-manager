package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

// PresignedURLRequest contains the request for a presigned URL
type PresignedURLRequest struct {
	FileName  string `json:"fileName"`
	MediaType string `json:"mediaType"` // "avatar", "trip", "location", "activity"
}

// PresignedURLResponse contains the presigned URL and expiration
type PresignedURLResponse struct {
	PresignedURL string `json:"presignedUrl"`
	ExpiresIn    int    `json:"expiresIn"` // in seconds
}

// GetPresignedURLHandler handles POST /api/uploads/presigned
// Returns a presigned URL for direct uploads to MinIO/S3
func GetPresignedURLHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var req PresignedURLRequest

		// Parse JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		// Validate request
		if req.FileName == "" {
			respondError(w, http.StatusBadRequest, "fileName is required")
			return
		}
		if req.MediaType == "" {
			respondError(w, http.StatusBadRequest, "mediaType is required")
			return
		}

		app.Logger.Printf("[Handler] GetPresignedURL: userId=%s, fileName=%s, mediaType=%s", userId, req.FileName, req.MediaType)

		// Map mediaType string to infrastructure.MediaType
		var mediaType infrastructure.MediaType
		switch req.MediaType {
		case "avatar":
			mediaType = infrastructure.MediaTypeAvatar
		case "trip":
			mediaType = infrastructure.MediaTypeTrip
		case "location":
			mediaType = infrastructure.MediaTypeLocation
		case "activity":
			mediaType = infrastructure.MediaTypeActivity
		default:
			respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid mediaType: %s", req.MediaType))
			return
		}

		// Generate presigned URL via MediaService
		presignedURL, err := app.Services.Media.GeneratePresignedURL(r.Context(), infrastructure.PresignedURLOptions{
			MediaType: mediaType,
			UserID:    userId,
			FileName:  req.FileName,
			ExpiresIn: 15 * time.Minute,
		})
		if err != nil {
			app.Logger.Printf("[Handler] GetPresignedURL: Failed to generate presigned URL: %v", err)
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate presigned URL: %v", err))
			return
		}

		app.Logger.Printf("[Handler] GetPresignedURL: Successfully generated presigned URL for userId=%s", userId)

		response := PresignedURLResponse{
			PresignedURL: presignedURL,
			ExpiresIn:    15 * 60, // 15 minutes in seconds
		}

		respondJSON(w, http.StatusOK, response)
	}
}



