package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type ActivityService interface {
	// GetActivity retrieves an activity by its ID.
	GetActivity(ctx context.Context, id string) (*generated.ActivityResponse, error)

	// CreateActivity creates a new activity with the provided details.
	CreateActivity(ctx context.Context, request *generated.CreateActivityRequest, tripId string, userId string) (*generated.ActivityResponse, error)

	// UpdateActivity updates an existing activity's details.
	UpdateActivity(ctx context.Context, request *generated.UpdateActivityRequest, activityId string, userId string) (*generated.ActivityResponse, error)

	// ListActivitiesForTrip retrieves all activities for a given trip ID.
	ListActivitiesForTrip(ctx context.Context, limit int, offset int, tripId string) (*generated.ActivityListResponse, error)

	// ListActivitiesForLocation retrieves all activities for a given location ID.
	ListActivitiesForLocation(ctx context.Context, limit int, offset int, locationId string) (*generated.ActivityListResponse, error)

	// DeleteActivity removes an activity from the system by its ID.
	DeleteActivity(ctx context.Context, id string, userId string) error
}

type ActivityServiceImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	activityRepository repository.ActivityRepository
	tripRepository     repository.TripRepository
}

// CreateActivity implements [ActivityService].
func (a *ActivityServiceImpl) CreateActivity(ctx context.Context, request *generated.CreateActivityRequest, tripId string, userId string) (*generated.ActivityResponse, error) {
	// 1. Validate input (business logic validation)
	if err := validateCreateActivityRequest(*request); err != nil {
		return nil, err
	}

	// 2. fetch trip which this activity belongs to, and verify it belongs to
	trip, err := a.tripRepository.GetTrip(ctx, tripId)
	if err != nil {
		return nil, err
	}

	activity := mapActivityCreateRequestToActivity(request, tripId, userId, trip.CreatedBy.Name, trip.CreatedBy.Email)

	// 3. validate activity in context of trip
	if trip.CreatedBy.ID != userId {
		return nil, fmt.Errorf("%w: user is not authorized to create activity", domain.ErrUnauthorized)
	}

	err = validateActivity(activity, trip)
	if err != nil {
		return nil, err
	}

	// 5. Call repository to persist
	activity, err = a.activityRepository.CreateActivity(ctx, activity)
	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	// 6. Convert from domain record back to response type
	response := mapActivityToActivityResponse(activity)
	return response, nil
}

// DeleteActivity implements [ActivityService].
func (a *ActivityServiceImpl) DeleteActivity(ctx context.Context, id string, userId string) error {
	return a.activityRepository.DeleteActivity(ctx, id, userId)
}

// GetActivity implements [ActivityService].
func (a *ActivityServiceImpl) GetActivity(ctx context.Context, id string) (*generated.ActivityResponse, error) {
	activityRecord, err := a.activityRepository.GetActivity(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}
	return mapActivityToActivityResponse(activityRecord), nil
}

// ListActivitiesForTrip implements [ActivityService].
func (a *ActivityServiceImpl) ListActivitiesForTrip(ctx context.Context, limit int, offset int, tripId string) (*generated.ActivityListResponse, error) {
	// Use repository to get activities for the specified trip
	activities, totalCount, err := a.activityRepository.ListActivitiesForTrip(ctx, tripId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities for trip: %w", err)
	}

	activityResponses := make([]generated.ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = *mapActivityToActivityResponse(activity)
	}

	return &generated.ActivityListResponse{
		Data:   activityResponses,
		Limit:  limit,
		Offset: offset,
		Total:  totalCount,
	}, nil
}

// ListActivitiesForLocation implements [ActivityService].
func (a *ActivityServiceImpl) ListActivitiesForLocation(ctx context.Context, limit int, offset int, locationId string) (*generated.ActivityListResponse, error) {
	activities, totalCount, err := a.activityRepository.ListActivitiesForLocation(ctx, locationId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities for location: %w", err)
	}

	activityResponses := make([]generated.ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = *mapActivityToActivityResponse(activity)
	}

	return &generated.ActivityListResponse{
		Data:   activityResponses,
		Limit:  limit,
		Offset: offset,
		Total:  totalCount,
	}, nil
}

// UpdateActivity implements [ActivityService].
func (a *ActivityServiceImpl) UpdateActivity(ctx context.Context, request *generated.UpdateActivityRequest, activityId string, userId string) (*generated.ActivityResponse, error) {
	// 1. validate business logic
	if err := validateUpdateActivityRequest(*request); err != nil {
		return nil, err
	}

	// 2. fetch activity to be updated, and verify it belongs to user
	activity, err := a.activityRepository.GetActivity(ctx, activityId)
	if err != nil {
		return nil, err
	}

	// 3. get corresponding trip
	trip, err := a.tripRepository.GetTrip(ctx, activity.TripID)
	if err != nil {
		return nil, err
	}

	// 4. validate in context of trip
	if trip.CreatedBy.ID != userId {
		return nil, fmt.Errorf("%w: user is not authorized to update activity", domain.ErrUnauthorized)
	}

	// 5. convert to activity
	activity = mapActivityUpdateRequestToActivity(request, activity)
	err = validateActivity(activity, trip)
	if err != nil {
		return nil, err
	}

	// 6. Call repository to persist
	activity, err = a.activityRepository.UpdateActivity(ctx, activity)
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}
	response := mapActivityToActivityResponse(activity)
	return response, nil
}

func NewActivityService(activityRepository repository.ActivityRepository) ActivityService {
	return &ActivityServiceImpl{
		activityRepository: activityRepository,
	}
}
