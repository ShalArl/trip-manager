package accommodation

import (
	"context"
	"fmt"
	"time"
)

// ── Input Types ───────────────────────────────────────────────────────────────

type CreateInput struct {
	TripID     string
	Name       string
	Address    *string
	CheckIn    *time.Time
	CheckOut   *time.Time
	BookingRef *string
	Notes      *string
}

type UpdateInput struct {
	ID         string
	Name       *string
	Address    *string
	CheckIn    *time.Time
	CheckOut   *time.Time
	BookingRef *string
	Notes      *string
}

// ── Service ───────────────────────────────────────────────────────────────────

type Service interface {
	ListByTrip(ctx context.Context, tripID string) ([]*Accommodation, error)
	GetByID(ctx context.Context, id string) (*Accommodation, error)
	Create(ctx context.Context, input CreateInput) (*Accommodation, error)
	Update(ctx context.Context, input UpdateInput) (*Accommodation, error)
	Delete(ctx context.Context, id string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListByTrip(ctx context.Context, tripID string) ([]*Accommodation, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Accommodation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, input CreateInput) (*Accommodation, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	return s.repo.Create(ctx, &Accommodation{
		TripID:     input.TripID,
		Name:       input.Name,
		Address:    input.Address,
		CheckIn:    input.CheckIn,
		CheckOut:   input.CheckOut,
		BookingRef: input.BookingRef,
		Notes:      input.Notes,
	})
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*Accommodation, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Address != nil {
		existing.Address = input.Address
	}
	if input.CheckIn != nil {
		existing.CheckIn = input.CheckIn
	}
	if input.CheckOut != nil {
		existing.CheckOut = input.CheckOut
	}
	if input.BookingRef != nil {
		existing.BookingRef = input.BookingRef
	}
	if input.Notes != nil {
		existing.Notes = input.Notes
	}
	return s.repo.Update(ctx, existing)
}

func (s *serviceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
