package like

import (
	"errors"
	"net/http"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/social/internal/shared"
)

// GetTripLikesHandler handles GET /trips/{tripId}/likes (optional authclient)
func GetTripLikesHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		// userID is optional for this endpoint
		userID, _ := authclient.GetUserID(r)

		resp, err := svc.GetEntityLikeInfo(r.Context(), userID, tripID, TargetTypeTrip)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusOK, resp)
	}
}

// LikeTripHandler handles POST /trips/{tripId}/likes (authclient required)
func LikeTripHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			shared.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripID := r.PathValue("tripId")
		if tripID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		err := svc.LikeEntity(r.Context(), userID, tripID, TargetTypeTrip)
		if err != nil {
			if errors.Is(err, shared.ErrConflict) {
				shared.RespondError(w, http.StatusConflict, "already liked")
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// UnlikeTripHandler handles DELETE /trips/{tripId}/likes (authclient required)
func UnlikeTripHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			shared.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripID := r.PathValue("tripId")
		if tripID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		err := svc.UnlikeEntity(r.Context(), userID, tripID, TargetTypeTrip)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetCommentLikesHandler handles GET /comments/{commentId}/likes (optional authclient)
func GetCommentLikesHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID := r.PathValue("commentId")
		if commentID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Comment ID is required")
			return
		}

		userID, _ := authclient.GetUserID(r)

		resp, err := svc.GetEntityLikeInfo(r.Context(), userID, commentID, TargetTypeComment)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusOK, resp)
	}
}

// LikeCommentHandler handles POST /comments/{commentId}/likes (authclient required)
func LikeCommentHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			shared.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := r.PathValue("commentId")
		if commentID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Comment ID is required")
			return
		}

		err := svc.LikeEntity(r.Context(), userID, commentID, TargetTypeComment)
		if err != nil {
			if errors.Is(err, shared.ErrConflict) {
				shared.RespondError(w, http.StatusConflict, "already liked")
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// UnlikeCommentHandler handles DELETE /comments/{commentId}/likes (authclient required)
func UnlikeCommentHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			shared.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := r.PathValue("commentId")
		if commentID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Comment ID is required")
			return
		}

		err := svc.UnlikeEntity(r.Context(), userID, commentID, TargetTypeComment)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
