package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapTransportToTransportResponse(transport *domain.Transport) *generated.TransportResponse {
	id, _ := uuid.Parse(transport.ID)
	fromLocID, _ := uuid.Parse(transport.FromLocationID)
	toLocID, _ := uuid.Parse(transport.ToLocationID)
	transportType := generated.TransportResponseType(transport.Type)

	return &generated.TransportResponse{
		Id:             ptr.ToPtr(id),
		FromLocationId: fromLocID,
		ToLocationId:   toLocID,
		DepartureTime:  transport.DepartureTime,
		ArrivalTime:    transport.ArrivalTime,
		Type:           transportType,
		Notes:          ptr.ToPtrNonEmpty(transport.Notes),
		CreatedAt:      &transport.CreatedAt,
		UpdatedAt:      &transport.UpdatedAt,
		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(transport.CreatedBy.ID),
			Name:  transport.CreatedBy.Name,
			Email: openapitypes.Email(transport.CreatedBy.Email),
		},
	}
}

func mapTransportsToTransportListResponse(transports []*domain.Transport, limit int, offset int, total int) *generated.TransportListResponse {
	items := make([]generated.TransportResponse, len(transports))
	for i, transport := range transports {
		items[i] = *mapTransportToTransportResponse(transport)
	}
	return &generated.TransportListResponse{
		Data:   items,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
