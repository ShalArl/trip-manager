package feed

import "context"

type Service interface {
	GetFeed(ctx context.Context, limit, offset int) ([]FeedTrip, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetFeed(ctx context.Context, limit, offset int) ([]FeedTrip, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetFeed(ctx, limit, offset)
}
