package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type LocationService interface {
	GetLocation(ctx context.Context, id string) (*domain.Location, error)
	CreateLocation(ctx context.Context, request *generated.CreateLocationRequest, tripId string, userId string) (*domain.Location, error)
	UpdateLocation(ctx context.Context, request *generated.UpdateLocationRequest, tripId string, userId string) (*domain.Location, error)
	ListLocations(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error)
	DeleteLocation(ctx context.Context, id string, userId string) error
	AddLocationImage(ctx context.Context, locationId string, userId string, imageKey string, sequence *int) (*domain.LocationImage, error)
	DeleteLocationImage(ctx context.Context, locationId string, imageId string, userId string) error
}

type LocationServiceImpl struct {
	locationRepository repository.LocationRepository
	mediaService       infrastructure.MediaService
}

// CreateLocation implements [LocationService].
func (l *LocationServiceImpl) CreateLocation(ctx context.Context, request *generated.CreateLocationRequest, tripId string, userId string) (*domain.Location, error) {
	if err := validateCreateLocationRequest(*request); err != nil {
		return nil, err
	}

	location := mapCreateLocationRequestToLocation(request, tripId, userId)

	createdLocation, err := l.locationRepository.CreateLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

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

// ListLocations implements [LocationService].
func (l *LocationServiceImpl) ListLocations(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error) {
	locations, totalCount, err := l.locationRepository.ListLocations(ctx, tripId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}
	return locations, totalCount, nil
}

// UpdateLocation implements [LocationService].
func (l *LocationServiceImpl) UpdateLocation(ctx context.Context, request *generated.UpdateLocationRequest, locationId string, _ string) (*domain.Location, error) {
	if err := validateUpdateLocationRequest(*request); err != nil {
		return nil, err
	}

	existing, err := l.locationRepository.GetLocation(ctx, locationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	record := mapUpdateLocationRequestToLocation(request, existing)

	updated, err := l.locationRepository.UpdateLocation(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return updated, nil
}

// AddLocationImage implements [LocationService].
func (l *LocationServiceImpl) AddLocationImage(ctx context.Context, locationId string, userId string, imageKey string, sequence *int) (*domain.LocationImage, error) {
	// Validate prefix: locations/{locationId}/...
	if !strings.HasPrefix(imageKey, "locations/") {
		return nil, fmt.Errorf("%w: image key must be a location image", domain.ErrInvalidInput)
	}

	fmt.Printf("[AddLocationImage] locationId=%s, imageKey=%s\n", locationId, imageKey)

	// Check file exists in storage
	exists, err := l.mediaService.ConfirmUpload(ctx, imageKey)
	if err != nil {
		return nil, fmt.Errorf("verify image upload: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("%w: image not uploaded", domain.ErrInvalidInput)
	}

	image, err := l.locationRepository.AddLocationImage(ctx, locationId, imageKey, sequence)
	if err != nil {
		return nil, fmt.Errorf("failed to add location image: %w", err)
	}

	return image, nil
}

// DeleteLocationImage implements [LocationService].
func (l *LocationServiceImpl) DeleteLocationImage(ctx context.Context, locationId string, imageId string, userId string) error {
	// Verify location exists and belongs to user
	location, err := l.locationRepository.GetLocation(ctx, locationId)
	if err != nil {
		return fmt.Errorf("failed to get location: %w", err)
	}

	if location.CreatedBy.ID != userId {
		return domain.ErrForbidden
	}

	if err := l.locationRepository.DeleteLocationImage(ctx, imageId, locationId); err != nil {
		return fmt.Errorf("failed to delete location image: %w", err)
	}

	return nil
}

func NewLocationService(locationRepository repository.LocationRepository, mediaService infrastructure.MediaService) LocationService {
	return &LocationServiceImpl{
		locationRepository: locationRepository,
		mediaService:       mediaService,
	}
}
