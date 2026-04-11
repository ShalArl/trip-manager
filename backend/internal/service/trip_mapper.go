package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
)

func mapCreateTripRequestToTrip(request *generated.CreateTripRequest, userID string, userName string, userEmail string) *domain.Trip {
	return &domain.Trip{
		ResourceMeta: domain.ResourceMeta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: domain.UserSummary{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
			},
		},
		Title:            request.Title,
		ShortDescription: request.ShortDescription,
		Description:      ptr.FromPtr(request.Description),
		StartDate:        request.StartDate.Time,
		EndDate:          request.EndDate.Time,
		Status:           domain.TripStatusPlanned,
	}
}

func mapUpdateTripRequestToTrip(request *generated.UpdateTripRequest, existing *domain.Trip) *domain.Trip {
	updated := *existing

	if request.Title != nil {
		updated.Title = *request.Title
	}

	if request.Description != nil {
		updated.Description = *request.Description
	}

	if request.ShortDescription != nil {
		updated.ShortDescription = *request.ShortDescription
	}

	if request.StartDate != nil {
		updated.StartDate = request.StartDate.Time
	}

	if request.EndDate != nil {
		updated.EndDate = request.EndDate.Time
	}

	if request.Status != nil {
		updated.Status = domain.TripStatus(*request.Status)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}
