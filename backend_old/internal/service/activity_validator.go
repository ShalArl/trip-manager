package service

import (
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
)

func validateCreateActivityRequest(request generated.CreateActivityRequest) error {
	if request.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}

	if request.Date.Time.IsZero() {
		return fmt.Errorf("%w: date is required", domain.ErrInvalidInput)
	}

	var category *generated.UpdateActivityRequestCategory
	if request.Category != nil {
		cat := generated.UpdateActivityRequestCategory(*request.Category)
		category = &cat
	}

	if err := validateCategoryAndTime(category, request.StartTime, request.EndTime); err != nil {
		return err
	}

	return nil
}

func validateUpdateActivityRequest(request generated.UpdateActivityRequest) error {
	if request.Name != nil && *request.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidInput)
	}

	if request.Date != nil && request.Date.Time.IsZero() {
		return fmt.Errorf("%w: date is invalid", domain.ErrInvalidInput)
	}

	if err := validateCategoryAndTime(request.Category, request.StartTime, request.EndTime); err != nil {
		return err
	}

	return nil
}

// validateCategoryAndTime checks if category is valid and if times are in correct order
func validateCategoryAndTime(category *generated.UpdateActivityRequestCategory, startTime, endTime *string) error {
	if category != nil && !category.Valid() {
		return fmt.Errorf("%w: invalid category", domain.ErrInvalidInput)
	}

	if startTime != nil && endTime != nil {
		if err := validateTimeRange(*startTime, *endTime); err != nil {
			return err
		}
	}
	return nil
}

// validateTimeRange checks if end time is after start time
func validateTimeRange(startTime, endTime string) error {
	const layout = "15:04"
	start, errS := time.Parse(layout, startTime)
	end, errE := time.Parse(layout, endTime)

	if errS != nil || errE != nil {
		return fmt.Errorf("%w: invalid time format (HH:mm)", domain.ErrInvalidInput)
	}

	if !end.After(start) {
		return fmt.Errorf("%w: end time must be after start time", domain.ErrInvalidInput)
	}
	return nil
}

func validateActivity(activity *domain.Activity, trip *domain.Trip) error {
	if activity.Date.Before(trip.StartDate) || activity.Date.After(trip.EndDate) {
		return fmt.Errorf("%w: activity date (%s) is outside of trip period (%s - %s)",
			domain.ErrInvalidInput,
			activity.Date.Format("2006-01-02"),
			trip.StartDate.Format("2006-01-02"),
			trip.EndDate.Format("2006-01-02"))
	}

	if activity.StartTime != "" && activity.EndTime != "" {
		const layout = "15:04"
		start, _ := time.Parse(layout, activity.StartTime)
		end, _ := time.Parse(layout, activity.EndTime)

		if !end.After(start) {
			return fmt.Errorf("%w: end time (%s) must be after start time (%s)",
				domain.ErrInvalidInput, activity.EndTime, activity.StartTime)
		}
	}

	return nil
}
