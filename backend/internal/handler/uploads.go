package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

// GetPresignedURLHandler handles POST /api/uploads/presigned
// Returns a presigned URL for direct uploads to MinIO/S3 or GCS depending on config
func GetPresignedURLHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req generated.PresignedURLRequest

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

		app.Logger.Printf("[Handler] GetPresignedURL: userId=%s, fileName=%s, mediaType=%s", userID, req.FileName, req.MediaType)

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
		ticket, err := app.Services.Media.PrepareUpload(r.Context(), userID, mediaType, req.FileName)
		if err != nil {
			app.Logger.Printf("[Handler] GetPresignedURL: Failed to generate presigned URL: %v", err)
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate presigned URL: %v", err))
			return
		}

		app.Logger.Printf("[Handler] GetPresignedURL: Successfully generated presigned URL for userId=%s", userID)

		response := generated.PresignedURLResponse{
			PresignedUrl: ticket.UploadURL,
			ExpiresIn:    int(ticket.ExpiresIn.Seconds()),
			Key:          ticket.Key,
		}

		respondJSON(w, http.StatusOK, response)
	}
}
