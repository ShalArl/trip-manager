package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
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
		Title:       request.Title,
		Description: ptr.FromPtr(request.Description),
		StartDate:   request.StartDate.Time,
		EndDate:     request.EndDate.Time,
		Status:      domain.TripStatusPlanned,
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

func mapTripToTripResponse(trip *domain.Trip) *generated.TripResponse {
	id, _ := uuid.Parse(trip.ID)

	status := generated.TripResponseStatus(trip.Status)

	return &generated.TripResponse{
		Id:          ptr.ToPtr(id),
		Title:       trip.Title,
		Description: ptr.ToPtr(trip.Description),
		StartDate:   openapitypes.Date{Time: trip.StartDate},
		EndDate:     openapitypes.Date{Time: trip.EndDate},
		Status:      status,
		CreatedAt:   &trip.CreatedAt,
		UpdatedAt:   &trip.UpdatedAt,
		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(trip.CreatedBy.ID),
			Name:  trip.CreatedBy.Name,
			Email: openapitypes.Email(trip.CreatedBy.Email),
		},
	}
}
