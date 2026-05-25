package transport

import (
	"context"
	"fmt"

	generated "github.com/ShalArl/trip-manager/backend/trips/generated"
	"github.com/google/uuid"
)

// ── Service ───────────────────────────────────────────────────────────────────

type Service interface {
	GetByID(ctx context.Context, id string) (*Transport, error)
	Create(ctx context.Context, req *generated.CreateTransportRequest, tripID, userID, userName, userEmail string) (*Transport, error)
	Update(ctx context.Context, req *generated.UpdateTransportRequest, transportID, userID string) (*Transport, error)
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Transport, int, error)
	Delete(ctx context.Context, id, userID string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Transport, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, req *generated.CreateTransportRequest, tripID, userID, userName, userEmail string) (*Transport, error) {
	if req.FromLocationId == uuid.Nil || req.ToLocationId == uuid.Nil {
		return nil, fmt.Errorf("%w: from_location_id and to_location_id are required", ErrInvalidInput)
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	t := &Transport{
		TripID: tripID,
		CreatedBy: UserSummary{
			ID:    userID,
			Name:  userName,
			Email: userEmail,
		},
		FromLocationID: req.FromLocationId.String(),
		ToLocationID:   req.ToLocationId.String(),
		DepartureTime:  req.DepartureTime,
		ArrivalTime:    req.ArrivalTime,
		Type:           string(req.Type),
		Notes:          notes,
	}
	return s.repo.Create(ctx, t)
}

func (s *serviceImpl) Update(ctx context.Context, req *generated.UpdateTransportRequest, transportID, userID string) (*Transport, error) {
	existing, err := s.repo.GetByID(ctx, transportID)
	if err != nil {
		return nil, err
	}
	if req.FromLocationId != nil {
		existing.FromLocationID = req.FromLocationId.String()
	}
	if req.ToLocationId != nil {
		existing.ToLocationID = req.ToLocationId.String()
	}
	if req.DepartureTime != nil {
		existing.DepartureTime = req.DepartureTime
	}
	if req.ArrivalTime != nil {
		existing.ArrivalTime = req.ArrivalTime
	}
	if req.Type != nil {
		existing.Type = string(*req.Type)
	}
	if req.Notes != nil {
		existing.Notes = *req.Notes
	}
	existing.CreatedBy.ID = userID
	return s.repo.Update(ctx, existing)
}

func (s *serviceImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Transport, int, error) {
	return s.repo.ListByTrip(ctx, tripID, limit, offset)
}

func (s *serviceImpl) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
