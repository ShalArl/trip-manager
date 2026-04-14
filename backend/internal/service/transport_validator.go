package service

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/google/uuid"
)

func validateCreateTransportRequest(req generated.CreateTransportRequest) error {
	if req.FromLocationId == uuid.Nil || req.ToLocationId == uuid.Nil {
		return fmt.Errorf("%w: from_location_id and to_location_id are required", domain.ErrInvalidInput)
	}
	if req.Date.IsZero() {
		return fmt.Errorf("%w: date is required", domain.ErrInvalidInput)
	}
	return nil
}

func validateUpdateTransportRequest(req generated.UpdateTransportRequest) error {
	if req.FromLocationId != nil && *req.FromLocationId == uuid.Nil {
		return fmt.Errorf("%w: from_location_id cannot be empty", domain.ErrInvalidInput)
	}
	if req.ToLocationId != nil && *req.ToLocationId == uuid.Nil {
		return fmt.Errorf("%w: to_location_id cannot be empty", domain.ErrInvalidInput)
	}
	return nil
}
