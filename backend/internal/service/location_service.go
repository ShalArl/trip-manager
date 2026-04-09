package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type LocationService interface {
	// GetLocation retrieves a location by its ID.
	GetLocation(ctx context.Context, id string) (*domain.Location, error)

	// CreateLocation creates a new location with the provided details.
	CreateLocation(ctx context.Context, request *generated.CreateLocationRequest, tripId string, userId string) (*domain.Location, error)

	// UpdateLocation updates an existing location's details.
	UpdateLocation(ctx context.Context, request *generated.UpdateLocationRequest, tripId string, userId string) (*domain.Location, error)

	// ListLocations retrieves all locations for a given trip ID.
	ListLocations(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error)

	// DeleteLocation removes a location from the system by its ID.
	DeleteLocation(ctx context.Context, id string, userId string) error
}

type LocationServiceImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	locationRepository repository.LocationRepository
}

// CreateLocation implements [LocationService].
func (l *LocationServiceImpl) CreateLocation(ctx context.Context, request *generated.CreateLocationRequest, tripId string, userId string) (*domain.Location, error) {
	// Validate input
	if err := validateCreateLocationRequest(*request); err != nil {
		return nil, err
	}

	// Convert from generated type to domain
	location := mapCreateLocationRequestToLocation(request, tripId, userId)

	// Call repository to persist
	createdLocation, err := l.locationRepository.CreateLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	// Convert from domain back to response type
	return createdLocation, nil
}

// DeleteLocation implements [LocationService].
func (l *LocationServiceImpl) DeleteLocation(ctx context.Context, id string, userId string) error {
	return l.locationRepository.DeleteLocation(ctx, id, userId)
}

// GetLocation implements [LocationService].
func (l *LocationServiceImpl) GetLocation(ctx context.Context, id string) (*domain.Location, error) {
	location, err := l.locationRepository.GetLocation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	return location, nil
}

func (l *LocationServiceImpl) ListLocations(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error) {
	locations, totalCount, err := l.locationRepository.ListLocations(ctx, tripId, "", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations, totalCount, nil
}

// UpdateLocation implements [LocationService].
func (l *LocationServiceImpl) UpdateLocation(ctx context.Context, request *generated.UpdateLocationRequest, tripId string, _ string) (*domain.Location, error) {
	err := validateUpdateLocationRequest(*request)
	if err != nil {
		return nil, err
	}

	// Fetch existing location
	existing, err := l.locationRepository.GetLocation(ctx, tripId)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Update with new values
	record := mapUpdateLocationRequestToLocation(request, existing)

	// Persist changes
	updated, err := l.locationRepository.UpdateLocation(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}
	return updated, nil
}

func NewLocationService(locationRepository repository.LocationRepository) LocationService {
	return &LocationServiceImpl{
		locationRepository: locationRepository,
	}
}
