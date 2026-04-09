package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapLocationToLocationResponse(l *domain.Location) *generated.LocationResponse {
	id, _ := uuid.Parse(l.ID)

	return &generated.LocationResponse{
		Id:        ptr.ToPtr(id),
		Name:      l.Name,
		City:      l.City,
		Country:   l.Country,
		Latitude:  ptr.ToPtr(float32(l.Coordinates.Lat)),
		Longitude: ptr.ToPtr(float32(l.Coordinates.Lon)),
		Notes:     ptr.ToPtr(l.Notes),
		Sequence:  ptr.ToPtr(l.Sequence),
		CreatedAt: &l.CreatedAt,
		UpdatedAt: &l.UpdatedAt,
		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(l.CreatedBy.ID),
			Name:  l.CreatedBy.Name,
			Email: openapitypes.Email(l.CreatedBy.Email),
		},
	}
}

func mapLocationsToLocationListResponse(locations []*domain.Location, limit int, offset int, total int) *generated.LocationListResponse {
	loc := make([]generated.LocationResponse, len(locations))
	for i, location := range locations {
		loc[i] = *mapLocationToLocationResponse(location)
	}

	return &generated.LocationListResponse{
		Data:   loc,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
