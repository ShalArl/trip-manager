package handler

import (
	"context"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapTripToTripResponse(ctx context.Context, trip *domain.Trip, media infrastructure.MediaService) *generated.TripResponse {
	id, _ := uuid.Parse(trip.ID)

	status := generated.TripResponseStatus(trip.Status)

	var avatarURL *string
	if trip.CreatedBy.AvatarKey != nil {
		if url, err := media.GetDownloadURL(ctx, *trip.CreatedBy.AvatarKey); err == nil {
			avatarURL = &url
		}
	}

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
			Id:        uuid.MustParse(trip.CreatedBy.ID),
			Name:      trip.CreatedBy.Name,
			Email:     openapitypes.Email(trip.CreatedBy.Email),
			AvatarUrl: avatarURL,
		},
	}
}

func mapTripsToTripListResponse(ctx context.Context, trips []*domain.Trip, limit int, offset int, total int, media infrastructure.MediaService) generated.TripListResponse {
	tr := make([]generated.TripResponse, len(trips))
	for i, trip := range trips {
		tr[i] = *mapTripToTripResponse(ctx, trip, media)
	}

	return generated.TripListResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Data:   tr,
	}

}
