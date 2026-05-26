package feed

import (
	"context"

	generated "github.com/ShalArl/trip-manager/backend/feed/generated"
)

type Service interface {
	GetFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Personalisierter Feed wenn eingeloggt, sonst globaler Feed
	if userID != "" {
		return s.repo.GetPersonalizedFeed(ctx, userID, limit, offset)
	}
	return s.repo.GetFeed(ctx, limit, offset)
}
