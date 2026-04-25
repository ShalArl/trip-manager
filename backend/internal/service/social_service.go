// internal/service/social_service.go
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type SocialService interface {
	ListComments(ctx context.Context, activityID string, limit, offset int) ([]*domain.Comment, int, error)
	CreateComment(ctx context.Context, userID, activityID, text, imageKey string) (*domain.Comment, error)
	UpdateComment(ctx context.Context, userID, commentID, text, imageKey string) (*domain.Comment, error)
	DeleteComment(ctx context.Context, userID, commentID string) error

	ListReplies(ctx context.Context, commentID string, limit, offset int) ([]*domain.CommentReply, int, error)
	CreateReply(ctx context.Context, userID, commentID, text, imageKey string) (*domain.CommentReply, error)
	UpdateReply(ctx context.Context, userID, replyID, text, imageKey string) (*domain.CommentReply, error)
	DeleteReply(ctx context.Context, userID, replyID string) error

	LikeActivity(ctx context.Context, userID, activityID string) error
	UnlikeActivity(ctx context.Context, userID, activityID string) error
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error
	LikeReply(ctx context.Context, userID, replyID string) error
	UnlikeReply(ctx context.Context, userID, replyID string) error

	GetActivityCounts(ctx context.Context, activityID string) (likes, comments int, err error)
	GetCommentCounts(ctx context.Context, commentID string) (likes, replies int, err error)
	GetReplyCounts(ctx context.Context, replyID string) (likes int, err error)
}

type SocialServiceImpl struct {
	repo repository.SocialRepository
}

func NewSocialService(repo repository.SocialRepository) SocialService {
	return &SocialServiceImpl{repo: repo}
}

// ── Comments ──────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) ListComments(ctx context.Context, activityID string, limit, offset int) ([]*domain.Comment, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListComments(ctx, activityID, limit, offset)
}

func (s *SocialServiceImpl) CreateComment(ctx context.Context, userID, activityID, text, imageKey string) (*domain.Comment, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("%w: comment text required", domain.ErrInvalidInput)
	}
	if len(text) > 500 {
		return nil, fmt.Errorf("%w: comment too long (max 500)", domain.ErrInvalidInput)
	}
	if imageKey != "" && !strings.HasPrefix(imageKey, "comments/") {
		return nil, fmt.Errorf("%w: invalid image key", domain.ErrInvalidInput)
	}

	c := &domain.Comment{
		ActivityID: activityID,
		UserID:     userID,
		Text:       text,
		ImageKey:   &imageKey,
	}
	return s.repo.CreateComment(ctx, c)
}

func (s *SocialServiceImpl) UpdateComment(ctx context.Context, userID, commentID, text, imageKey string) (*domain.Comment, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("%w: comment text required", domain.ErrInvalidInput)
	}

	existing, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	existing.Text = text
	existing.ImageKey = &imageKey
	return s.repo.UpdateComment(ctx, existing)
}

func (s *SocialServiceImpl) DeleteComment(ctx context.Context, userID, commentID string) error {
	existing, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return domain.ErrUnauthorized
	}
	return s.repo.DeleteComment(ctx, commentID)
}

// ── Replies ───────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) ListReplies(ctx context.Context, commentID string, limit, offset int) ([]*domain.CommentReply, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListReplies(ctx, commentID, limit, offset)
}

func (s *SocialServiceImpl) CreateReply(ctx context.Context, userID, commentID, text, imageKey string) (*domain.CommentReply, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("%w: reply text required", domain.ErrInvalidInput)
	}
	if len(text) > 500 {
		return nil, fmt.Errorf("%w: reply too long (max 500)", domain.ErrInvalidInput)
	}

	// Parent-Comment muss existieren
	if _, err := s.repo.GetComment(ctx, commentID); err != nil {
		return nil, err
	}

	rep := &domain.CommentReply{
		CommentID: commentID,
		UserID:    userID,
		Text:      text,
		ImageKey:  &imageKey,
	}
	return s.repo.CreateReply(ctx, rep)
}

func (s *SocialServiceImpl) UpdateReply(ctx context.Context, userID, replyID, text, imageKey string) (*domain.CommentReply, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("%w: reply text required", domain.ErrInvalidInput)
	}

	existing, err := s.repo.GetReply(ctx, replyID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	existing.Text = text
	existing.ImageKey = &imageKey
	return s.repo.UpdateReply(ctx, existing)
}

func (s *SocialServiceImpl) DeleteReply(ctx context.Context, userID, replyID string) error {
	existing, err := s.repo.GetReply(ctx, replyID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return domain.ErrUnauthorized
	}
	return s.repo.DeleteReply(ctx, replyID)
}

// ── Likes ─────────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) LikeActivity(ctx context.Context, userID, activityID string) error {
	return s.repo.LikeActivity(ctx, userID, activityID)
}

func (s *SocialServiceImpl) UnlikeActivity(ctx context.Context, userID, activityID string) error {
	return s.repo.UnlikeActivity(ctx, userID, activityID)
}

func (s *SocialServiceImpl) LikeComment(ctx context.Context, userID, commentID string) error {
	return s.repo.LikeComment(ctx, userID, commentID)
}

func (s *SocialServiceImpl) UnlikeComment(ctx context.Context, userID, commentID string) error {
	return s.repo.UnlikeComment(ctx, userID, commentID)
}

func (s *SocialServiceImpl) LikeReply(ctx context.Context, userID, replyID string) error {
	return s.repo.LikeReply(ctx, userID, replyID)
}

func (s *SocialServiceImpl) UnlikeReply(ctx context.Context, userID, replyID string) error {
	return s.repo.UnlikeReply(ctx, userID, replyID)
}

// ── Counts ────────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) GetActivityCounts(ctx context.Context, activityID string) (int, int, error) {
	likes, err := s.repo.CountActivityLikes(ctx, activityID)
	if err != nil {
		return 0, 0, err
	}
	comments, err := s.repo.CountActivityComments(ctx, activityID)
	if err != nil {
		return 0, 0, err
	}
	return likes, comments, nil
}

func (s *SocialServiceImpl) GetCommentCounts(ctx context.Context, commentID string) (int, int, error) {
	likes, err := s.repo.CountCommentLikes(ctx, commentID)
	if err != nil {
		return 0, 0, err
	}
	replies, err := s.repo.CountCommentReplies(ctx, commentID)
	if err != nil {
		return 0, 0, err
	}
	return likes, replies, nil
}

func (s *SocialServiceImpl) GetReplyCounts(ctx context.Context, replyID string) (int, error) {
	return s.repo.CountReplyLikes(ctx, replyID)
}
