package transport

import (
	"context"
	"fmt"
	"time"
)

// ── Service ───────────────────────────────────────────────────────────────────

type CreateInput struct {
	TripID         string
	Type           string
	DeparturePlace string
	ArrivalPlace   string
	DepartureTime  *time.Time
	ArrivalTime    *time.Time
	BookingRef     *string
	Notes          *string
}

type UpdateInput struct {
	ID             string
	Type           *string
	DeparturePlace *string
	ArrivalPlace   *string
	DepartureTime  *time.Time
	ArrivalTime    *time.Time
	BookingRef     *string
	Notes          *string
}

type Service interface {
	ListByTrip(ctx context.Context, tripID string) ([]*Transport, error)
	GetByID(ctx context.Context, id string) (*Transport, error)
	Create(ctx context.Context, input CreateInput) (*Transport, error)
	Update(ctx context.Context, input UpdateInput) (*Transport, error)
	Delete(ctx context.Context, id string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListByTrip(ctx context.Context, tripID string) ([]*Transport, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Transport, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, input CreateInput) (*Transport, error) {
	if input.Type == "" {
		return nil, fmt.Errorf("%w: type is required", ErrInvalidInput)
	}
	if input.DeparturePlace == "" {
		return nil, fmt.Errorf("%w: departure_place is required", ErrInvalidInput)
	}
	if input.ArrivalPlace == "" {
		return nil, fmt.Errorf("%w: arrival_place is required", ErrInvalidInput)
	}
	return s.repo.Create(ctx, &Transport{
		TripID:         input.TripID,
		Type:           input.Type,
		DeparturePlace: input.DeparturePlace,
		ArrivalPlace:   input.ArrivalPlace,
		DepartureTime:  input.DepartureTime,
		ArrivalTime:    input.ArrivalTime,
		BookingRef:     input.BookingRef,
		Notes:          input.Notes,
	})
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*Transport, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Type != nil {
		existing.Type = *input.Type
	}
	if input.DeparturePlace != nil {
		existing.DeparturePlace = *input.DeparturePlace
	}
	if input.ArrivalPlace != nil {
		existing.ArrivalPlace = *input.ArrivalPlace
	}
	if input.DepartureTime != nil {
		existing.DepartureTime = input.DepartureTime
	}
	if input.ArrivalTime != nil {
		existing.ArrivalTime = input.ArrivalTime
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
