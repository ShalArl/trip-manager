package service

import (
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
)

func validateCreateTripRequest(request generated.CreateTripRequest) error {
	if request.Title == "" {
		return fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}

	if request.StartDate.Time.IsZero() {
		return fmt.Errorf("%w: start date is required", domain.ErrInvalidInput)
	}

	if request.ShortDescription == "" {
		return fmt.Errorf("%w: short description is required", domain.ErrInvalidInput)
	}

	if len(request.ShortDescription) > 80 {
		return fmt.Errorf("%w: short description is too long", domain.ErrInvalidInput)
	}

	if request.Destination == "" {
		return fmt.Errorf("%w: destination is required", domain.ErrInvalidInput)
	}

	// TODO: Reactivate
	// if request.EndDate.Time.IsZero() {
	//	return fmt.Errorf("%w: end date is required", domain.ErrInvalidInput)
	// }

	// return validateTripDateRange(request.StartDate.Time, request.EndDate.Time)

	return nil
}

func validateUpdateTripRequest(request generated.UpdateTripRequest) error {
	if request.Title != nil && *request.Title == "" {
		return fmt.Errorf("%w: title cannot be empty", domain.ErrInvalidInput)
	}

	if request.Status != nil && !request.Status.Valid() {
		return fmt.Errorf("%w: invalid status", domain.ErrInvalidInput)
	}

	if request.ShortDescription != nil && *request.ShortDescription == "" {
		return fmt.Errorf("%w: short description cannot be empty", domain.ErrInvalidInput)
	}

	if request.ShortDescription != nil && len(*request.ShortDescription) > 80 {
		return fmt.Errorf("%w: short description is too long", domain.ErrInvalidInput)
	}

	// If both dates are provided, validate the range
	// if request.StartDate != nil && request.EndDate != nil {
	// 	 return validateTripDateRange(request.StartDate.Time, request.EndDate.Time)
	// }

	return nil
}

// validateTripDateRange checks if end date is after or equal to start date
func validateTripDateRange(startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		return fmt.Errorf("%w: end date must be after or equal to start date", domain.ErrInvalidInput)
	}
	return nil
}
