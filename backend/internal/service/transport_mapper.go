package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
)

func mapCreateTransportRequestToTransport(req *generated.CreateTransportRequest, tripID string, userID string) *domain.Transport {
	return &domain.Transport{
		ResourceMeta: domain.ResourceMeta{
			CreatedBy: domain.UserSummary{ID: userID},
		},
		TripID:         tripID,
		FromLocationID: req.FromLocationId.String(),
		ToLocationID:   req.ToLocationId.String(),
		DepartureTime:  req.DepartureTime,
		ArrivalTime:    req.ArrivalTime,
		Type:           domain.TransportType(req.Type),
		Notes:          ptr.FromPtr(req.Notes),
	}
}

func mapUpdateTransportRequestToTransport(req *generated.UpdateTransportRequest, existing *domain.Transport) *domain.Transport {
	updated := *existing
	if req.FromLocationId != nil {
		updated.FromLocationID = req.FromLocationId.String()
	}
	if req.ToLocationId != nil {
		updated.ToLocationID = req.ToLocationId.String()
	}
	if req.DepartureTime != nil {
		updated.DepartureTime = req.DepartureTime
	}
	if req.ArrivalTime != nil {
		updated.ArrivalTime = req.ArrivalTime
	}
	if req.Type != nil {
		updated.Type = domain.TransportType(*req.Type)
	}
	if req.Notes != nil {
		updated.Notes = *req.Notes
	}
	updated.UpdatedAt = time.Now()
	return &updated
}
