package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
)

func mapCreateLocationRequestToLocation(req *generated.CreateLocationRequest, tripID string, userID string) *domain.Location {
	return &domain.Location{
		ResourceMeta: domain.ResourceMeta{
			CreatedBy: domain.UserSummary{ID: userID},
		},
		TripID:           tripID,
		Name:             req.Name,
		City:             req.City,
		Country:          req.Country,
		ShortDescription: req.ShortDescription,
		DateFrom:         req.DateFrom.Time,
		DateTo:           req.DateTo.Time,
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
	if req.ShortDescription != nil {
		updated.ShortDescription = *req.ShortDescription
	}
	if req.DateFrom != nil {
		updated.DateFrom = req.DateFrom.Time
	}
	if req.DateTo != nil {
		updated.DateTo = req.DateTo.Time
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
