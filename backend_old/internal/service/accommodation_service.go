package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type AccommodationService interface {
	GetAccommodation(ctx context.Context, id string) (*domain.Accommodation, error)
	CreateAccommodation(ctx context.Context, request *generated.CreateAccommodationRequest, tripID string, userID string) (*domain.Accommodation, error)
	UpdateAccommodation(ctx context.Context, request *generated.UpdateAccommodationRequest, accommodationID string, userID string) (*domain.Accommodation, error)
	ListAccommodations(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Accommodation, int, error)
	DeleteAccommodation(ctx context.Context, id string, userID string) error
}

type AccommodationServiceImpl struct {
	accommodationRepository repository.AccommodationRepository
}

func NewAccommodationService(accommodationRepository repository.AccommodationRepository) AccommodationService {
	return &AccommodationServiceImpl{
		accommodationRepository: accommodationRepository,
	}
}

func (s *AccommodationServiceImpl) GetAccommodation(ctx context.Context, id string) (*domain.Accommodation, error) {
	accommodation, err := s.accommodationRepository.GetAccommodation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get accommodation: %w", err)
	}
	return accommodation, nil
}

func (s *AccommodationServiceImpl) CreateAccommodation(ctx context.Context, request *generated.CreateAccommodationRequest, tripID string, userID string) (*domain.Accommodation, error) {
	if err := validateCreateAccommodationRequest(*request); err != nil {
		return nil, err
	}
	accommodation := mapCreateAccommodationRequestToAccommodation(request, tripID, userID)
	created, err := s.accommodationRepository.CreateAccommodation(ctx, accommodation)
	if err != nil {
		return nil, fmt.Errorf("failed to create accommodation: %w", err)
	}
	return created, nil
}

func (s *AccommodationServiceImpl) UpdateAccommodation(ctx context.Context, request *generated.UpdateAccommodationRequest, accommodationID string, userID string) (*domain.Accommodation, error) {
	if err := validateUpdateAccommodationRequest(*request); err != nil {
		return nil, err
	}
	existing, err := s.accommodationRepository.GetAccommodation(ctx, accommodationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accommodation: %w", err)
	}
	if existing.CreatedBy.ID != userID {
		return nil, domain.ErrForbidden
	}
	updated := mapUpdateAccommodationRequestToAccommodation(request, existing)
	result, err := s.accommodationRepository.UpdateAccommodation(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update accommodation: %w", err)
	}
	return result, nil
}

func (s *AccommodationServiceImpl) ListAccommodations(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Accommodation, int, error) {
	accommodations, totalCount, err := s.accommodationRepository.ListAccommodationsForTrip(ctx, tripID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list accommodations: %w", err)
	}
	return accommodations, totalCount, nil
}

func (s *AccommodationServiceImpl) DeleteAccommodation(ctx context.Context, id string, userID string) error {
	return s.accommodationRepository.DeleteAccommodation(ctx, id, userID)
}
