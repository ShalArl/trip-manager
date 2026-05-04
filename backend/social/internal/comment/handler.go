package comment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/social/internal/shared"
)

// ListTripCommentsHandler handles GET /trips/{tripId}/comments (optional authclient)
func ListTripCommentsHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		resp, err := svc.ListComments(r.Context(), tripID)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusOK, resp)
	}
}

// ListRepliesHandler handles GET /comments/{commentId}/replies (optional authclient)
func ListRepliesHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID := r.PathValue("commentId")
		if commentID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Comment ID is required")
			return
		}

		resp, err := svc.ListComments(r.Context(), commentID)
		if err != nil {
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusOK, resp)
	}
}

// CreateTripCommentHandler handles POST /trips/{tripId}/comments (authclient required)
func CreateTripCommentHandler(svc Service) http.HandlerFunc {
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

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		resp, err := svc.CreateComment(r.Context(), userID, tripID, req.Text)
		if err != nil {
			if errors.Is(err, shared.ErrInvalidInput) {
				shared.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusCreated, resp)
	}
}

// CreateReplyHandler handles POST /comments/{commentId}/replies (authclient required)
func CreateReplyHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			shared.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := r.PathValue("commentId")
		if commentID == "" {
			shared.RespondError(w, http.StatusBadRequest, "Comment ID is required")
		}

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		resp, err := svc.CreateComment(r.Context(), userID, commentID, req.Text)
		if err != nil {
			if errors.Is(err, shared.ErrInvalidInput) {
				shared.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, shared.ErrNotFound) {
				shared.RespondError(w, http.StatusNotFound, "parent comment not found")
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusCreated, resp)

	}
}

// UpdateCommentHandler handles PUT /comments/{commentId} (authclient required)
func UpdateCommentHandler(svc Service) http.HandlerFunc {
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

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		resp, err := svc.UpdateComment(r.Context(), userID, commentID, req.Text)
		if err != nil {
			if errors.Is(err, shared.ErrInvalidInput) {
				shared.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, shared.ErrNotFound) {
				shared.RespondError(w, http.StatusNotFound, "comment not found")
				return
			}
			if errors.Is(err, shared.ErrForbidden) {
				shared.RespondError(w, http.StatusForbidden, "forbidden")
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		shared.RespondJSON(w, http.StatusOK, resp)
	}
}

// DeleteCommentHandler handles DELETE /comments/{commentId} (authclient required)
func DeleteCommentHandler(svc Service) http.HandlerFunc {
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

		err := svc.DeleteComment(r.Context(), userID, commentID)
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				shared.RespondError(w, http.StatusNotFound, "comment not found")
				return
			}
			if errors.Is(err, shared.ErrForbidden) {
				shared.RespondError(w, http.StatusForbidden, "forbidden")
				return
			}
			shared.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
