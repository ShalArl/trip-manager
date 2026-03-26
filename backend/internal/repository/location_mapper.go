package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (rec *locationRecord) toLocation() *domain.Location {
	return &domain.Location{
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
		TripID:  rec.TripID.String(),
		Name:    rec.Name,
		City:    rec.City,
		Country: rec.Country,
		Coordinates: domain.Point{
			Lat: ptr.FromPtr(rec.Latitude),
			Lon: ptr.FromPtr(rec.Longitude),
		},
		Notes:    ptr.FromPtr(rec.Notes),
		Sequence: ptr.FromPtr(rec.Sequence),
	}
}

func locationToRecord(location *domain.Location) (*locationRecord, error) {
	var locationID uuid.UUID
	var err error

	if location.ID != "" {
		locationID, err = uuid.Parse(location.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid location ID: %w", err)
		}
	}

	tripID, err := uuid.Parse(location.TripID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for trip ID: %v", err)
	}

	userID, err := uuid.Parse(location.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for creator ID: %v", err)
	}

	return &locationRecord{
		ID:        locationID,
		TripID:    tripID,
		Name:      location.Name,
		City:      location.City,
		Country:   location.Country,
		Latitude:  ptr.ToPtr(location.Coordinates.Lat),
		Longitude: ptr.ToPtr(location.Coordinates.Lon),
		Notes:     ptr.ToPtr(location.Notes),
		Sequence:  ptr.ToPtrNonEmpty(location.Sequence),
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
		UserID:    userID,
		UserName:  location.CreatedBy.Name,
		UserEmail: location.CreatedBy.Email,
	}, nil
}
