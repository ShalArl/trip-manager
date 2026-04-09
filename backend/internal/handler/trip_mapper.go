package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapTripToTripResponse(trip *domain.Trip) *generated.TripResponse {
	id, _ := uuid.Parse(trip.ID)

	status := generated.TripResponseStatus(trip.Status)

	return &generated.TripResponse{
		Id:               ptr.ToPtr(id),
		Title:            trip.Title,
		ShortDescription: trip.ShortDescription,
		Description:      ptr.ToPtr(trip.Description),
		StartDate:        openapitypes.Date{Time: trip.StartDate},
		EndDate:          openapitypes.Date{Time: trip.EndDate},
		Status:           status,
		CreatedAt:        &trip.CreatedAt,
		UpdatedAt:        &trip.UpdatedAt,
		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(trip.CreatedBy.ID),
			Name:  trip.CreatedBy.Name,
			Email: openapitypes.Email(trip.CreatedBy.Email),
		},
	}
}

func mapTripsToTripListResponse(trips []*domain.Trip, limit int, offset int, total int) generated.TripListResponse {
	tr := make([]generated.TripResponse, len(trips))
	for i, trip := range trips {
		tr[i] = *mapTripToTripResponse(trip)
	}

	return generated.TripListResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Data:   tr,
	}

}
