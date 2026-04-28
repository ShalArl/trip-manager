package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (rec *accommodationRecord) toAccommodation() *domain.Accommodation {
	return &domain.Accommodation{
		ResourceMeta: domain.ResourceMeta{
			ID:        rec.ID.String(),
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			CreatedBy: domain.UserSummary{
				ID:    rec.UserID.String(),
				Name:  rec.UserName,
				Email: rec.UserEmail,
			},
		},
		TripID:        rec.TripID.String(),
		LocationID:    rec.LocationID.String(),
		Name:          rec.Name,
		Address:       ptr.FromPtr(rec.Address),
		CheckIn:       rec.CheckIn,
		CheckOut:      rec.CheckOut,
		PricePerNight: rec.PricePerNight,
		Notes:         ptr.FromPtr(rec.Notes),
	}
}

func accommodationToRecord(a *domain.Accommodation) (*accommodationRecord, error) {
	var accommodationID uuid.UUID
	var err error
	if a.ID != "" {
		accommodationID, err = uuid.Parse(a.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid accommodation ID: %w", err)
		}
	}
	tripID, err := uuid.Parse(a.TripID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for trip ID: %v", err)
	}
	locationID, err := uuid.Parse(a.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for location ID: %v", err)
	}
	userID, err := uuid.Parse(a.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for creator ID: %v", err)
	}
	return &accommodationRecord{
		ID:            accommodationID,
		TripID:        tripID,
		UserID:        userID,
		LocationID:    locationID,
		Name:          a.Name,
		Address:       ptr.ToPtrNonEmpty(a.Address),
		CheckIn:       a.CheckIn,
		CheckOut:      a.CheckOut,
		PricePerNight: a.PricePerNight,
		Notes:         ptr.ToPtrNonEmpty(a.Notes),
		UserName:      a.CreatedBy.Name,
		UserEmail:     a.CreatedBy.Email,
	}, nil
}
