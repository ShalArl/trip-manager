package service

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/google/uuid"
)

func validateCreateAccommodationRequest(req generated.CreateAccommodationRequest) error {
	if req.LocationId == uuid.Nil {
		return fmt.Errorf("%w: location_id is required", domain.ErrInvalidInput)
	}
	if req.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	return nil
}

func validateUpdateAccommodationRequest(req generated.UpdateAccommodationRequest) error {
	if req.LocationId != nil && *req.LocationId == uuid.Nil {
		return fmt.Errorf("%w: location_id cannot be empty", domain.ErrInvalidInput)
	}
	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidInput)
	}
	return nil
}
