package comment

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/backend/social/internal/shared"
)

type Service interface {
	CreateComment(ctx context.Context, userID, entityID, text string) (*CommentResponse, error)
	ListComments(ctx context.Context, entityID string) (*CommentListResponse, error)
	DeleteComment(ctx context.Context, userID, commentID string) error
	UpdateComment(ctx context.Context, userID, commentID, text string) (*CommentResponse, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewServiceImpl(repository Repository) Service {
	return &ServiceImpl{
		repo: repository,
	}
}

// CreateComment creates a new comment
func (s *ServiceImpl) CreateComment(ctx context.Context, userID, entityID, text string) (*CommentResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("%w: comment text cannot be empty", shared.ErrInvalidInput)
	}

	comment := &Comment{
		EntityID: entityID,
		UserID:   userID,
		Text:     text,
	}

	created, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &CommentResponse{
		ID:        created.ID,
		EntityID:  created.EntityID,
		UserID:    created.UserID,
		Text:      created.Text,
		ImageKey:  created.ImageKey,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

// ListComments lists all comments for an entity
func (s *ServiceImpl) ListComments(ctx context.Context, entityID string) (*CommentListResponse, error) {
	comments, err := s.repo.ListComments(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	var responses []*CommentResponse
	for _, c := range comments {
		responses = append(responses, &CommentResponse{
			ID:        c.ID,
			EntityID:  c.EntityID,
			UserID:    c.UserID,
			Text:      c.Text,
			ImageKey:  c.ImageKey,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return &CommentListResponse{Data: responses}, nil
}

// DeleteComment deletes a comment
func (s *ServiceImpl) DeleteComment(ctx context.Context, userID, commentID string) error {
	if err := s.repo.DeleteComment(ctx, commentID, userID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// UpdateComment updates a comment
func (s *ServiceImpl) UpdateComment(ctx context.Context, userID, commentID, text string) (*CommentResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("%w: comment text cannot be empty", shared.ErrInvalidInput)
	}

	// Get existing comment to verify ownership
	existing, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	// Check ownership
	if existing.UserID != userID {
		return nil, shared.ErrForbidden
	}

	// Update comment
	existing.Text = text
	updated, err := s.repo.UpdateComment(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &CommentResponse{
		ID:        updated.ID,
		EntityID:  updated.EntityID,
		UserID:    updated.UserID,
		Text:      updated.Text,
		ImageKey:  updated.ImageKey,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}
