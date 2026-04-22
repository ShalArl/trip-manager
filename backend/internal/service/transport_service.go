package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type TransportService interface {
	GetTransport(ctx context.Context, id string) (*domain.Transport, error)
	CreateTransport(ctx context.Context, request *generated.CreateTransportRequest, tripId string, userId string) (*domain.Transport, error)
	UpdateTransport(ctx context.Context, request *generated.UpdateTransportRequest, transportId string, userId string) (*domain.Transport, error)
	ListTransports(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Transport, int, error)
	DeleteTransport(ctx context.Context, id string, userId string) error
}

type TransportServiceImpl struct {
	transportRepository repository.TransportRepository
}

func (t *TransportServiceImpl) GetTransport(ctx context.Context, id string) (*domain.Transport, error) {
	transport, err := t.transportRepository.GetTransport(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transport: %w", err)
	}
	return transport, nil
}

func (t *TransportServiceImpl) CreateTransport(ctx context.Context, request *generated.CreateTransportRequest, tripId string, userId string) (*domain.Transport, error) {
	if err := validateCreateTransportRequest(*request); err != nil {
		return nil, err
	}

	transport := mapCreateTransportRequestToTransport(request, tripId, userId)

	created, err := t.transportRepository.CreateTransport(ctx, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	return created, nil
}

func (t *TransportServiceImpl) UpdateTransport(ctx context.Context, request *generated.UpdateTransportRequest, transportId string, userId string) (*domain.Transport, error) {
	if err := validateUpdateTransportRequest(*request); err != nil {
		return nil, err
	}

	existing, err := t.transportRepository.GetTransport(ctx, transportId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transport: %w", err)
	}

	updated := mapUpdateTransportRequestToTransport(request, existing)

	result, err := t.transportRepository.UpdateTransport(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update transport: %w", err)
	}

	return result, nil
}

func (t *TransportServiceImpl) ListTransports(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Transport, int, error) {
	transports, totalCount, err := t.transportRepository.ListTransportsForTrip(ctx, tripId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list transports: %w", err)
	}
	return transports, totalCount, nil
}

func (t *TransportServiceImpl) DeleteTransport(ctx context.Context, id string, userId string) error {
	return t.transportRepository.DeleteTransport(ctx, id, userId)
}

func NewTransportService(transportRepository repository.TransportRepository) TransportService {
	return &TransportServiceImpl{
		transportRepository: transportRepository,
	}
}
