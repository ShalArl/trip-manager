package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/repository"
	"github.com/google/uuid"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

type SocialService interface {
	LikeTrip(ctx context.Context, userID, tripID string) error
	UnlikeTrip(ctx context.Context, userID, tripID string) error
	GetTripLikeInfo(ctx context.Context, userID, tripID string) (*generated.TripLikeResponse, error)
	CreateTripComment(ctx context.Context, userID, tripID, text string) (*generated.TripCommentResponse, error)
	ListTripComments(ctx context.Context, tripID string) (*generated.TripCommentListResponse, error)
	DeleteTripComment(ctx context.Context, userID, commentID string) error
}

type SocialServiceImpl struct {
	socialRepo   repository.SocialRepository
	userRepo     repository.UserRepository
	mediaService infrastructure.MediaService
}

func NewSocialService(socialRepo repository.SocialRepository, userRepo repository.UserRepository, mediaService infrastructure.MediaService) SocialService {
	return &SocialServiceImpl{
		socialRepo:   socialRepo,
		userRepo:     userRepo,
		mediaService: mediaService,
	}
}

// ── Likes ──────────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) LikeTrip(ctx context.Context, userID, tripID string) error {
	if err := s.socialRepo.LikeTrip(ctx, userID, tripID); err != nil {
		return fmt.Errorf("failed to like trip: %w", err)
	}
	return nil
}

func (s *SocialServiceImpl) UnlikeTrip(ctx context.Context, userID, tripID string) error {
	if err := s.socialRepo.UnlikeTrip(ctx, userID, tripID); err != nil {
		return fmt.Errorf("failed to unlike trip: %w", err)
	}
	return nil
}

func (s *SocialServiceImpl) GetTripLikeInfo(ctx context.Context, userID, tripID string) (*generated.TripLikeResponse, error) {
	count, err := s.socialRepo.CountTripLikes(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to count likes: %w", err)
	}

	hasLiked := false
	if userID != "" {
		hasLiked, err = s.socialRepo.HasLiked(ctx, userID, tripID)
		if err != nil {
			return nil, fmt.Errorf("failed to check like status: %w", err)
		}
	}

	return &generated.TripLikeResponse{
		LikeCount: count,
		HasLiked:  hasLiked,
	}, nil
}

// ── Comments ───────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) CreateTripComment(ctx context.Context, userID, tripID, text string) (*generated.TripCommentResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("%w: comment text cannot be empty", domain.ErrInvalidInput)
	}

	comment := &domain.TripComment{
		TripID: tripID,
		UserID: userID,
		Text:   text,
	}

	created, err := s.socialRepo.CreateTripComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.mapToTripCommentResponse(ctx, created, user), nil
}

func (s *SocialServiceImpl) ListTripComments(ctx context.Context, tripID string) (*generated.TripCommentListResponse, error) {
	comments, err := s.socialRepo.ListTripComments(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	var responses []generated.TripCommentResponse
	for _, c := range comments {
		user, err := s.userRepo.GetUser(ctx, c.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user for comment: %w", err)
		}
		responses = append(responses, *s.mapToTripCommentResponse(ctx, c, user))
	}

	return &generated.TripCommentListResponse{Data: responses}, nil
}

func (s *SocialServiceImpl) DeleteTripComment(ctx context.Context, userID, commentID string) error {
	if err := s.socialRepo.DeleteTripComment(ctx, commentID, userID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// ── Mapper ─────────────────────────────────────────────────────────────────

func (s *SocialServiceImpl) mapToTripCommentResponse(ctx context.Context, c *domain.TripComment, user *domain.User) *generated.TripCommentResponse {
	id, _ := uuid.Parse(user.ID)

	userSummary := &generated.UserSummary{
		Id:    id,
		Email: openapitypes.Email(user.Email),
		Name:  user.Name,
	}

	if user.AvatarKey != "" {
		avatarUrl, err := s.mediaService.GetDownloadURL(ctx, user.AvatarKey)
		if err == nil {
			userSummary.AvatarUrl = &avatarUrl
		}
	}

	return &generated.TripCommentResponse{
		Id:        c.ID,
		TripId:    c.TripID,
		User:      userSummary,
		Text:      c.Text,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
