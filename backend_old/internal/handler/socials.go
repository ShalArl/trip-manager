package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

// GetTripLikesHandler handles GET /trips/{tripId}/likes (optional middleware)
func GetTripLikesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		// leer wenn Gast → hasLiked wird false
		userID, _ := middleware.GetUserID(r)

		resp, err := app.Services.Social.GetTripLikeInfo(r.Context(), userID, tripID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// LikeTripHandler handles POST /trips/{tripId}/likes (middleware required)
func LikeTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		err := app.Services.Social.LikeTrip(r.Context(), userID, tripID)
		if err != nil {
			if errors.Is(err, domain.ErrConflict) {
				respondError(w, http.StatusConflict, "already liked")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// UnlikeTripHandler handles DELETE /trips/{tripId}/likes (middleware required)
func UnlikeTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		err := app.Services.Social.UnlikeTrip(r.Context(), userID, tripID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ListTripCommentsHandler handles GET /trips/{tripId}/comments (optional middleware)
func ListTripCommentsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		resp, err := app.Services.Social.ListTripComments(r.Context(), tripID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}

// CreateTripCommentHandler handles POST /trips/{tripId}/comments (middleware required)
func CreateTripCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		var req generated.CreateTripCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		resp, err := app.Services.Social.CreateTripComment(r.Context(), userID, tripID, req.Text)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidInput) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, resp)
	}
}

// DeleteTripCommentHandler handles DELETE /trips/{tripId}/comments/{commentId} (middleware required)
func DeleteTripCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := r.PathValue("commentId")
		if commentID == "" {
			respondError(w, http.StatusBadRequest, "Comment ID is required")
			return
		}

		err := app.Services.Social.DeleteTripComment(r.Context(), userID, commentID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				respondError(w, http.StatusNotFound, "comment not found")
				return
			}
			if errors.Is(err, domain.ErrForbidden) {
				respondError(w, http.StatusForbidden, "forbidden")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
