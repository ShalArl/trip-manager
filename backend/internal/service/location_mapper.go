package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapCreateLocationRequestToLocation(req *generated.CreateLocationRequest, tripID string, userID string) *domain.Location {
	return &domain.Location{
		ResourceMeta: domain.ResourceMeta{
			CreatedBy: domain.UserSummary{ID: userID},
		},
		TripID:  tripID,
		Name:    req.Name,
		City:    req.City,
		Country: req.Country,
		Coordinates: domain.Point{
			Lat: float64(ptr.FromPtr(req.Latitude)),
			Lon: float64(ptr.FromPtr(req.Longitude)),
		},
		Notes:    ptr.FromPtr(req.Notes),
		Sequence: ptr.FromPtr(req.Sequence),
	}
}

func mapUpdateLocationRequestToLocation(req *generated.UpdateLocationRequest, existing *domain.Location) *domain.Location {
	updated := *existing

	if req.Name != nil {
		updated.Name = *req.Name
	}

	if req.City != nil {
		updated.City = *req.City
	}

	if req.Country != nil {
		updated.Country = *req.Country
	}

	if req.Latitude != nil {
		updated.Coordinates.Lat = float64(*req.Latitude)
	}

	if req.Longitude != nil {
		updated.Coordinates.Lon = float64(*req.Longitude)
	}

	if req.Notes != nil {
		updated.Notes = *req.Notes
	}

	if req.Sequence != nil {
		updated.Sequence = *req.Sequence
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

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
