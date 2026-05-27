package feed

import (
	"context"

	generated "github.com/ShalArl/trip-manager/backend/feed/generated"
)

type Service interface {
	GetGlobalFeed(ctx context.Context, limit, offset int) ([]generated.FeedTrip, int, error)
	GetPersonalizedFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetGlobalFeed(ctx context.Context, limit, offset int) ([]generated.FeedTrip, int, error) {
	limit, offset = clamp(limit, offset)
	return s.repo.GetFeed(ctx, limit, offset)
}

func (s *service) GetPersonalizedFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error) {
	limit, offset = clamp(limit, offset)
	return s.repo.GetPersonalizedFeed(ctx, userID, limit, offset)
}

func clamp(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
