package service

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
)

func validateCreateLocationRequest(req generated.CreateLocationRequest) error {
	if req.Name == "" || req.City == "" || req.Country == "" {
		return fmt.Errorf("%w: name, city, and country are required", domain.ErrInvalidInput)
	}

	return validateCoordinates(req.Latitude, req.Longitude)
}

func validateUpdateLocationRequest(req generated.UpdateLocationRequest) error {
	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidInput)
	}

	return validateCoordinates(req.Latitude, req.Longitude)
}

// validateCoordinates checks if latitude and longitude are within valid ranges
func validateCoordinates(latitude, longitude *float32) error {
	if latitude != nil && (*latitude < -90 || *latitude > 90) {
		return fmt.Errorf("%w: latitude must be between -90 and 90", domain.ErrInvalidInput)
	}
	if longitude != nil && (*longitude < -180 || *longitude > 180) {
		return fmt.Errorf("%w: longitude must be between -180 and 180", domain.ErrInvalidInput)
	}

	return nil
}

