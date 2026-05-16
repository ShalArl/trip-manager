package trip

import (
	"context"
	"fmt"
	"time"
)

// ── Input Types ───────────────────────────────────────────────────────────────

type CreateInput struct {
	Title            string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	UserID           string
	UserName         string
	UserEmail        string
}

type UpdateInput struct {
	ID               string
	UserID           string
	Title            *string
	ShortDescription *string
	Description      *string
	StartDate        *time.Time
	EndDate          *time.Time
	Status           *string
}

// ── Interface ─────────────────────────────────────────────────────────────────

type Service interface {
	GetByID(ctx context.Context, id string) (*Trip, error)
	Create(ctx context.Context, input CreateInput) (*Trip, error)
	Update(ctx context.Context, input UpdateInput) (*Trip, error)
	Delete(ctx context.Context, id, userID string) error
	List(ctx context.Context, userID string, limit, offset int) ([]*Trip, int, error)
	ListRecent(ctx context.Context, limit, offset int) ([]*Trip, int, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Trip, int, error)
}

// ── Implementation ────────────────────────────────────────────────────────────

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Trip, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, input CreateInput) (*Trip, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if input.ShortDescription == "" {
		return nil, fmt.Errorf("%w: short description is required", ErrInvalidInput)
	}
	if len(input.ShortDescription) > 80 {
		return nil, fmt.Errorf("%w: short description is too long", ErrInvalidInput)
	}
	if input.EndDate.Before(input.StartDate) {
		return nil, fmt.Errorf("%w: end date must be after start date", ErrInvalidInput)
	}

	return s.repo.Create(ctx, &Trip{
		Title:            input.Title,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		Status:           "planned",
		CreatedBy: UserSummary{
			ID:    input.UserID,
			Name:  input.UserName,
			Email: input.UserEmail,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*Trip, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if existing.CreatedBy.ID != input.UserID {
		return nil, ErrUnauthorized
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.ShortDescription != nil {
		existing.ShortDescription = *input.ShortDescription
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.StartDate != nil {
		existing.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		existing.EndDate = *input.EndDate
	}
	if input.Status != nil {
		existing.Status = *input.Status
	}
	existing.UpdatedAt = time.Now()
	return s.repo.Update(ctx, existing)
}

func (s *serviceImpl) Delete(ctx context.Context, id, userID string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.CreatedBy.ID != userID {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id, userID)
}

func (s *serviceImpl) List(ctx context.Context, userID string, limit, offset int) ([]*Trip, int, error) {
	return s.repo.List(ctx, userID, limit, offset)
}

func (s *serviceImpl) ListRecent(ctx context.Context, limit, offset int) ([]*Trip, int, error) {
	return s.repo.ListRecent(ctx, limit, offset)
}

func (s *serviceImpl) Search(ctx context.Context, query string, limit, offset int) ([]*Trip, int, error) {
	return s.repo.Search(ctx, query, limit, offset)
}
