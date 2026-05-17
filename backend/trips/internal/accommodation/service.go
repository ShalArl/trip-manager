package accommodation

import (
	"context"
	"fmt"

	generated "github.com/ShalArl/trip-manager/backend/trips/generated"
	"github.com/google/uuid"
)

// ── Service ───────────────────────────────────────────────────────────────────

type Service interface {
	GetByID(ctx context.Context, id string) (*Accommodation, error)
	Create(ctx context.Context, req *generated.CreateAccommodationRequest, tripID, userID, userName, userEmail string) (*Accommodation, error)
	Update(ctx context.Context, req *generated.UpdateAccommodationRequest, accommodationID, userID string) (*Accommodation, error)
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Accommodation, int, error)
	Delete(ctx context.Context, id, userID string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

// ── Validation ────────────────────────────────────────────────────────────────

func validateCreate(req *generated.CreateAccommodationRequest) error {
	if req.LocationId == uuid.Nil {
		return fmt.Errorf("%w: location_id is required", ErrInvalidInput)
	}
	if req.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	return nil
}

func validateUpdate(req *generated.UpdateAccommodationRequest) error {
	if req.LocationId != nil && *req.LocationId == uuid.Nil {
		return fmt.Errorf("%w: location_id cannot be empty", ErrInvalidInput)
	}
	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
	}
	return nil
}

// ── Mappers ───────────────────────────────────────────────────────────────────

func fromCreateRequest(req *generated.CreateAccommodationRequest, tripID, userID, userName, userEmail string) *Accommodation {
	address := ""
	if req.Address != nil {
		address = *req.Address
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}
	return &Accommodation{
		TripID: tripID,
		CreatedBy: UserSummary{
			ID:    userID,
			Name:  userName,
			Email: userEmail,
		},
		LocationID:    req.LocationId.String(),
		Name:          req.Name,
		Address:       address,
		CheckIn:       req.CheckIn,
		CheckOut:      req.CheckOut,
		PricePerNight: req.PricePerNight,
		Notes:         notes,
	}
}

func applyUpdate(req *generated.UpdateAccommodationRequest, existing *Accommodation) *Accommodation {
	updated := *existing
	if req.LocationId != nil {
		updated.LocationID = req.LocationId.String()
	}
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Address != nil {
		updated.Address = *req.Address
	}
	if req.CheckIn != nil {
		updated.CheckIn = req.CheckIn
	}
	if req.CheckOut != nil {
		updated.CheckOut = req.CheckOut
	}
	if req.PricePerNight != nil {
		updated.PricePerNight = req.PricePerNight
	}
	if req.Notes != nil {
		updated.Notes = *req.Notes
	}
	return &updated
}

// ── Implementation ────────────────────────────────────────────────────────────

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Accommodation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, req *generated.CreateAccommodationRequest, tripID, userID, userName, userEmail string) (*Accommodation, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}
	a := fromCreateRequest(req, tripID, userID, userName, userEmail)
	created, err := s.repo.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("failed to create accommodation: %w", err)
	}
	return created, nil
}

func (s *serviceImpl) Update(ctx context.Context, req *generated.UpdateAccommodationRequest, accommodationID, userID string) (*Accommodation, error) {
	if err := validateUpdate(req); err != nil {
		return nil, err
	}
	existing, err := s.repo.GetByID(ctx, accommodationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accommodation: %w", err)
	}
	if existing.CreatedBy.ID != userID {
		return nil, ErrUnauthorized
	}
	updated := applyUpdate(req, existing)
	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update accommodation: %w", err)
	}
	return result, nil
}

func (s *serviceImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Accommodation, int, error) {
	accommodations, total, err := s.repo.ListByTrip(ctx, tripID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list accommodations: %w", err)
	}
	return accommodations, total, nil
}

func (s *serviceImpl) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
