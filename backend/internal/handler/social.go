// internal/handler/social_handlers.go
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

func ListCommentsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityID := chi.URLParam(r, "activityId")
		if activityID == "" {
			respondError(w, http.StatusBadRequest, "activity id required")
			return
		}

		limit, offset := parseLimitOffset(r)

		comments, total, err := app.Services.Social.ListComments(r.Context(), activityID, limit, offset)
		if err != nil {
			app.Logger.Printf("[ListComments] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not load comments")
			return
		}

		// Counts pro Comment aggregieren — könnte man auch parallel machen via errgroup
		responses := make([]generated.CommentResponse, 0, len(comments))
		for _, c := range comments {
			likes, replies, err := app.Services.Social.GetCommentCounts(r.Context(), c.ID)
			if err != nil {
				app.Logger.Printf("[ListComments] count err for %s: %v", c.ID, err)
				// Counts-Fehler nicht fatal, mit 0 weitermachen
			}
			responses = append(responses, mapCommentToResponse(r.Context(), app, c, likes, replies))
		}

		respondJSON(w, http.StatusOK, generated.CommentListResponse{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Data:   responses,
		})
	}
}

func CreateCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		activityID := chi.URLParam(r, "activityId")
		if activityID == "" {
			respondError(w, http.StatusBadRequest, "activity id required")
			return
		}

		var req generated.CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid body")
			return
		}

		imageKey := ""
		if req.ImageKey != nil {
			imageKey = *req.ImageKey
		}

		comment, err := app.Services.Social.CreateComment(r.Context(), userID, activityID, req.Text, imageKey)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidInput) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			app.Logger.Printf("[CreateComment] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not create comment")
			return
		}

		respondJSON(w, http.StatusCreated, mapCommentToResponse(r.Context(), app, comment, 0, 0))
	}
}

func UpdateCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := chi.URLParam(r, "commentId")
		if commentID == "" {
			respondError(w, http.StatusBadRequest, "comment id required")
			return
		}

		var req generated.UpdateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid body")
			return
		}

		imageKey := ""
		if req.ImageKey != nil {
			imageKey = *req.ImageKey
		}

		comment, err := app.Services.Social.UpdateComment(r.Context(), userID, commentID, req.Text, imageKey)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				respondError(w, http.StatusNotFound, "comment not found")
				return
			}
			if errors.Is(err, domain.ErrForbidden) {
				respondError(w, http.StatusForbidden, "not your comment")
				return
			}
			if errors.Is(err, domain.ErrInvalidInput) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			app.Logger.Printf("[UpdateComment] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not update comment")
			return
		}

		likes, replies, _ := app.Services.Social.GetCommentCounts(r.Context(), commentID)
		respondJSON(w, http.StatusOK, mapCommentToResponse(r.Context(), app, comment, likes, replies))
	}
}

func DeleteCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		commentID := chi.URLParam(r, "commentId")

		if err := app.Services.Social.DeleteComment(r.Context(), userID, commentID); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				respondError(w, http.StatusNotFound, "comment not found")
				return
			}
			if errors.Is(err, domain.ErrForbidden) {
				respondError(w, http.StatusForbidden, "not your comment")
				return
			}
			app.Logger.Printf("[DeleteComment] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not delete comment")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ── Likes ─────────────────────────────────────────────────────────────────

func LikeActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		activityID := chi.URLParam(r, "activityId")

		if err := app.Services.Social.LikeActivity(r.Context(), userID, activityID); err != nil {
			app.Logger.Printf("[LikeActivity] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not like")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UnlikeActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		activityID := chi.URLParam(r, "activityId")

		if err := app.Services.Social.UnlikeActivity(r.Context(), userID, activityID); err != nil {
			app.Logger.Printf("[UnlikeActivity] err=%v", err)
			respondError(w, http.StatusInternalServerError, "could not unlike")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Analog für LikeCommentHandler, UnlikeCommentHandler, LikeReplyHandler, UnlikeReplyHandler
// — alle haben dieselbe Struktur, nur andere Service-Methode.

// ── Helpers ───────────────────────────────────────────────────────────────

func parseLimitOffset(r *http.Request) (int, int) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}
	return limit, offset
}

func mapCommentToResponse(ctx context.Context, app *app.App, c *domain.Comment, likes, replies int) generated.CommentResponse {
	var imageURL *string
	if *c.ImageKey != "" {
		if url, err := app.Services.Media.GetDownloadURL(ctx, *c.ImageKey); err == nil {
			imageURL = &url
		}
	}

	// User-Info aus Postgres holen für UserSummary
	user, _ := app.Services.User.GetUser(ctx, c.UserID)
	createdBy := mapUserToUserSummary(ctx, app.Services.Media, user)

	return generated.CommentResponse{
		Id:         c.ID,
		Text:       c.Text,
		ActivityId: c.ActivityID,
		ImageUrl:   imageURL,
		LikeCount:  &likes,
		ReplyCount: &replies,
		CreatedBy:  createdBy,
		CreatedAt:  &c.CreatedAt,
		UpdatedAt:  &c.UpdatedAt,
	}
}
