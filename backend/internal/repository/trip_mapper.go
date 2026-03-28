package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (t *tripRecord) toTrip() *domain.Trip {
	return &domain.Trip{
		ResourceMeta: domain.ResourceMeta{
			ID:        t.ID.String(),
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			CreatedBy: domain.UserSummary{
				ID:    t.UserID.String(),
				Name:  t.UserName,
				Email: t.UserEmail,
			},
		},
		Title:       t.Title,
		Description: ptr.FromPtr(t.Description),
		StartDate:   t.StartDate,
		EndDate:     t.EndDate,
		Status:      domain.TripStatus(t.Status),
	}
}

func tripToRecord(trip *domain.Trip) (*tripRecord, error) {
	var tripID uuid.UUID
	var err error

	if trip.ID != "" {
		tripID, err = uuid.Parse(trip.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid trip ID: %w", err)
		}
	}

	userID, err := uuid.Parse(trip.CreatedBy.ID)
	if err != nil {
		return nil, err
	}

	return &tripRecord{
		ID:          tripID,
		UserID:      userID,
		Title:       trip.Title,
		Description: ptr.ToPtr(trip.Description),
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Status:      string(trip.Status),
		CreatedAt:   trip.CreatedAt,
		UpdatedAt:   trip.UpdatedAt,
		UserName:    trip.CreatedBy.Name,
		UserEmail:   trip.CreatedBy.Email,
	}, nil
}
