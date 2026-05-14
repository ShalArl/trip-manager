package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/trips/repository"
)

// ── Types ─────────────────────────────────────────────────────────────────────

type Trip struct {
	ID               string
	Title            string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	Status           string
	CreatedBy        UserSummary
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserSummary struct {
	ID        string
	Name      string
	Email     string
	AvatarKey *string
}

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
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Trip, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoToService(t), nil
}

func (s *serviceImpl) Create(ctx context.Context, input CreateInput) (*Trip, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("%w: title is required", repository.ErrInvalidInput)
	}
	if input.ShortDescription == "" {
		return nil, fmt.Errorf("%w: short description is required", repository.ErrInvalidInput)
	}
	if len(input.ShortDescription) > 80 {
		return nil, fmt.Errorf("%w: short description is too long", repository.ErrInvalidInput)
	}
	if input.EndDate.Before(input.StartDate) {
		return nil, fmt.Errorf("%w: end date must be after start date", repository.ErrInvalidInput)
	}

	t, err := s.repo.Create(ctx, &repository.Trip{
		Title:            input.Title,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		Status:           "planned",
		CreatedBy: repository.UserSummary{
			ID:    input.UserID,
			Name:  input.UserName,
			Email: input.UserEmail,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return repoToService(t), nil
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*Trip, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if existing.CreatedBy.ID != input.UserID {
		return nil, repository.ErrUnauthorized
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

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return repoToService(updated), nil
}

func (s *serviceImpl) Delete(ctx context.Context, id, userID string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.CreatedBy.ID != userID {
		return repository.ErrUnauthorized
	}
	return s.repo.Delete(ctx, id, userID)
}

func (s *serviceImpl) List(ctx context.Context, userID string, limit, offset int) ([]*Trip, int, error) {
	trips, total, err := s.repo.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Trip, len(trips))
	for i, t := range trips {
		result[i] = repoToService(t)
	}
	return result, total, nil
}

func (s *serviceImpl) ListRecent(ctx context.Context, limit, offset int) ([]*Trip, int, error) {
	trips, total, err := s.repo.ListRecent(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Trip, len(trips))
	for i, t := range trips {
		result[i] = repoToService(t)
	}
	return result, total, nil
}

func (s *serviceImpl) Search(ctx context.Context, query string, limit, offset int) ([]*Trip, int, error) {
	trips, total, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Trip, len(trips))
	for i, t := range trips {
		result[i] = repoToService(t)
	}
	return result, total, nil
}

// ── Mapper ────────────────────────────────────────────────────────────────────

func repoToService(t *repository.Trip) *Trip {
	return &Trip{
		ID:               t.ID,
		Title:            t.Title,
		ShortDescription: t.ShortDescription,
		Description:      t.Description,
		StartDate:        t.StartDate,
		EndDate:          t.EndDate,
		Status:           t.Status,
		CreatedBy: UserSummary{
			ID:        t.CreatedBy.ID,
			Name:      t.CreatedBy.Name,
			Email:     t.CreatedBy.Email,
			AvatarKey: t.CreatedBy.AvatarKey,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func fromErrors(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return err
	}
	return err
}
