package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type TripService interface {
	// GetTrip retrieves a trip by its ID.
	GetTrip(ctx context.Context, id string) (*domain.Trip, error)

	// CreateTrip creates a new trip with the provided details.
	CreateTrip(ctx context.Context, request *generated.CreateTripRequest, userID, userName, userEmail string) (*domain.Trip, error)

	// UpdateTrip updates an existing trip's details.
	UpdateTrip(ctx context.Context, request *generated.UpdateTripRequest, id, userId string) (*domain.Trip, error)

	// ListTrips retrieves all trips for a given user ID.
	ListTrips(ctx context.Context, userID string, limit, offset int) ([]*domain.Trip, int, error)

	// DeleteTrip removes a trip from the system by its ID.
	DeleteTrip(ctx context.Context, id, userId string) error
}

type TripServiceImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	tripRepository     repository.TripRepository
	locationRepository repository.LocationRepository
	activityRepository repository.ActivityRepository
}

// CreateTrip implements [TripService].
// Service layer: validates input + converts types + coordinates repository
func (t *TripServiceImpl) CreateTrip(ctx context.Context, request *generated.CreateTripRequest, userID, userName, userEmail string) (*domain.Trip, error) {
	// 1. Validate input (business logic validation)
	if err := validateCreateTripRequest(*request); err != nil {
		return nil, err
	}

	// 2. Convert from generated type to domain
	trip := mapCreateTripRequestToTrip(request, userID, userName, userEmail)

	// 3. Call repository to persist
	createdTrip, err := t.tripRepository.CreateTrip(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	print("Created Trip-id:", createdTrip.ID)

	return createdTrip, nil
}

// DeleteTrip implements [TripService].
func (t *TripServiceImpl) DeleteTrip(ctx context.Context, id, userId string) error {
	return t.tripRepository.DeleteTrip(ctx, id, userId)
}

// GetTrip implements [TripService].
func (t *TripServiceImpl) GetTrip(ctx context.Context, id string) (*domain.Trip, error) {
	trip, err := t.tripRepository.GetTrip(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	if trip == nil {
		return nil, fmt.Errorf("trip not found")
	}

	return trip, nil
}

func (t *TripServiceImpl) ListTrips(ctx context.Context, userID string, limit, offset int) ([]*domain.Trip, int, error) {
	// Use repository to get trips
	trips, totalCount, err := t.tripRepository.ListTrips(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list trips: %w", err)
	}
	fmt.Printf("ListTrips trips count: %d, totalCount: %d\n", len(trips), totalCount)

	// Build response with pagination metadata
	return trips, totalCount, nil
}

// UpdateTrip implements [TripService].
func (t *TripServiceImpl) UpdateTrip(ctx context.Context, request *generated.UpdateTripRequest, id, userID string) (*domain.Trip, error) {
	// 1. Validate input (business logic validation)
	if err := validateUpdateTripRequest(*request); err != nil {
		return nil, err
	}

	// 2. Fetch existing trip (needed to merge with updates)
	existingTrip, err := t.tripRepository.GetTrip(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	// 3. Check authorization - user can only update their own trips
	if existingTrip.CreatedBy.ID != userID {
		return nil, fmt.Errorf("unauthorized to update trip")
	}

	// 4. Convert from generated type to domain (merges request with existing)
	trip := mapUpdateTripRequestToTrip(request, existingTrip)

	// 5. Call repository to update
	updatedTrip, err := t.tripRepository.UpdateTrip(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	return updatedTrip, nil
}

func NewTripService(tripRepository repository.TripRepository, locationRepository repository.LocationRepository, activityRepository repository.ActivityRepository) TripService {
	return &TripServiceImpl{
		tripRepository:     tripRepository,
		locationRepository: locationRepository,
		activityRepository: activityRepository,
	}
}
