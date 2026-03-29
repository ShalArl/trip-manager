package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type TripService interface {
	// GetTrip retrieves a trip by its ID.
	GetTrip(ctx context.Context, id string) (*generated.TripResponse, error)

	// CreateTrip creates a new trip with the provided details.
	CreateTrip(ctx context.Context, request *generated.CreateTripRequest, userID, userName, userEmail string) (*generated.TripResponse, error)

	// UpdateTrip updates an existing trip's details.
	UpdateTrip(ctx context.Context, request *generated.UpdateTripRequest, id string) (*generated.TripResponse, error)

	// ListTrips retrieves all trips for a given user ID.
	ListTrips(ctx context.Context, userID string, limit int, offset int) (*generated.TripListResponse, error)

	// DeleteTrip removes a trip from the system by its ID.
	DeleteTrip(ctx context.Context, id string, userId string) error
}

type TripServiceImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	tripRepository     repository.TripRepository
	locationRepository repository.LocationRepository
	activityRepository repository.ActivityRepository
}

// CreateTrip implements [TripService].
// Service layer: validates input + converts types + coordinates repository
func (t *TripServiceImpl) CreateTrip(ctx context.Context, request *generated.CreateTripRequest, userID, userName, userEmail string) (*generated.TripResponse, error) {
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

	// 4. Convert from domain back to response type
	response := mapTripToTripResponse(createdTrip)
	return response, nil
}

// DeleteTrip implements [TripService].
func (t *TripServiceImpl) DeleteTrip(ctx context.Context, id string, userId string) error {
	return t.tripRepository.DeleteTrip(ctx, id, userId)
}

// GetTrip implements [TripService].
func (t *TripServiceImpl) GetTrip(ctx context.Context, id string) (*generated.TripResponse, error) {
	trip, err := t.tripRepository.GetTrip(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	if trip == nil {
		return nil, fmt.Errorf("trip not found")
	}

	response := mapTripToTripResponse(trip)
	return response, nil
}

func (t *TripServiceImpl) ListTrips(ctx context.Context, userID string, limit int, offset int) (*generated.TripListResponse, error) {
	// Use repository to get trips
	trips, totalCount, err := t.tripRepository.ListTrips(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list trips: %w", err)
	}

	// Convert to responses
	tripResponses := make([]generated.TripResponse, len(trips))
	for i, trip := range trips {
		resp := mapTripToTripResponse(trip)
		if resp != nil {
			tripResponses[i] = *resp
		}
	}

	// Build response with pagination metadata
	return &generated.TripListResponse{
		Data:   tripResponses,
		Limit:  limit,
		Offset: offset,
		Total:  totalCount,
	}, nil
}

// UpdateTrip implements [TripService].
func (t *TripServiceImpl) UpdateTrip(ctx context.Context, request *generated.UpdateTripRequest, id string) (*generated.TripResponse, error) {
	// 1. Validate input (business logic validation)
	if err := validateUpdateTripRequest(*request); err != nil {
		return nil, err
	}

	// 2. Fetch existing trip
	existingTrip, err := t.tripRepository.GetTrip(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	// 3. Convert from generated type to domain
	trip := mapUpdateTripRequestToTrip(request, existingTrip)

	// 4. Call repository to update and get updated record
	updatedTrip, err := t.tripRepository.UpdateTrip(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	// 5. Convert back to response type
	response := mapTripToTripResponse(updatedTrip)
	return response, nil
}

func NewTripService(tripRepository repository.TripRepository, locationRepository repository.LocationRepository, activityRepository repository.ActivityRepository) TripService {
	return &TripServiceImpl{
		tripRepository:     tripRepository,
		locationRepository: locationRepository,
		activityRepository: activityRepository,
	}
}
