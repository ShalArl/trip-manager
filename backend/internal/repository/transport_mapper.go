package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (rec *transportRecord) toTransport() *domain.Transport {
	return &domain.Transport{
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
		TripID:         rec.TripID.String(),
		FromLocationID: rec.FromLocationID.String(),
		ToLocationID:   rec.ToLocationID.String(),
		Date:           rec.Date,
		Type:           domain.TransportType(rec.Type),
		Notes:          ptr.FromPtr(rec.Notes),
	}
}

func transportToRecord(transport *domain.Transport) (*transportRecord, error) {
	var transportID uuid.UUID
	var err error
	if transport.ID != "" {
		transportID, err = uuid.Parse(transport.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid transport ID: %w", err)
		}
	}

	tripID, err := uuid.Parse(transport.TripID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for trip ID: %v", err)
	}

	fromLocationID, err := uuid.Parse(transport.FromLocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for from_location ID: %v", err)
	}

	toLocationID, err := uuid.Parse(transport.ToLocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for to_location ID: %v", err)
	}

	userID, err := uuid.Parse(transport.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for creator ID: %v", err)
	}

	return &transportRecord{
		ID:             transportID,
		TripID:         tripID,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		Date:           transport.Date,
		Type:           string(transport.Type),
		Notes:          ptr.ToPtrNonEmpty(transport.Notes),
		UserID:         userID,
		UserName:       transport.CreatedBy.Name,
		UserEmail:      transport.CreatedBy.Email,
	}, nil
}
