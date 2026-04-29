package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
)

func mapCreateAccommodationRequestToAccommodation(req *generated.CreateAccommodationRequest, tripID string, userID string) *domain.Accommodation {
	return &domain.Accommodation{
		ResourceMeta: domain.ResourceMeta{
			CreatedBy: domain.UserSummary{ID: userID},
		},
		TripID:        tripID,
		LocationID:    req.LocationId.String(),
		Name:          req.Name,
		Address:       ptr.FromPtr(req.Address),
		CheckIn:       req.CheckIn,
		CheckOut:      req.CheckOut,
		PricePerNight: req.PricePerNight,
		Notes:         ptr.FromPtr(req.Notes),
	}
}

func mapUpdateAccommodationRequestToAccommodation(req *generated.UpdateAccommodationRequest, existing *domain.Accommodation) *domain.Accommodation {
	updated := *existing
	if req.LocationId != nil {
		updated.LocationID = req.LocationId.String()
	}
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Address != nil {
		updated.Address = *req.Address
	}
	if req.CheckIn != nil {
		updated.CheckIn = req.CheckIn
	}
	if req.CheckOut != nil {
		updated.CheckOut = req.CheckOut
	}
	if req.PricePerNight != nil {
		updated.PricePerNight = req.PricePerNight
	}
	if req.Notes != nil {
		updated.Notes = *req.Notes
	}
	updated.UpdatedAt = time.Now()
	return &updated
}
