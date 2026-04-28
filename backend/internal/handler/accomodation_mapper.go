package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapAccommodationToAccommodationResponse(a *domain.Accommodation) *generated.AccommodationResponse {
	id, _ := uuid.Parse(a.ID)
	locationID, _ := uuid.Parse(a.LocationID)

	return &generated.AccommodationResponse{
		Id:            ptr.ToPtr(id),
		LocationId:    locationID,
		Name:          a.Name,
		Address:       ptr.ToPtrNonEmpty(a.Address),
		CheckIn:       a.CheckIn,
		CheckOut:      a.CheckOut,
		PricePerNight: a.PricePerNight,
		Notes:         ptr.ToPtrNonEmpty(a.Notes),
		CreatedAt:     &a.CreatedAt,
		UpdatedAt:     &a.UpdatedAt,
		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(a.CreatedBy.ID),
			Name:  a.CreatedBy.Name,
			Email: openapitypes.Email(a.CreatedBy.Email),
		},
	}
}

func mapAccommodationsToAccommodationListResponse(accommodations []*domain.Accommodation, limit int, offset int, total int) *generated.AccommodationListResponse {
	items := make([]generated.AccommodationResponse, len(accommodations))
	for i, a := range accommodations {
		items[i] = *mapAccommodationToAccommodationResponse(a)
	}
	return &generated.AccommodationListResponse{
		Data:   items,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
