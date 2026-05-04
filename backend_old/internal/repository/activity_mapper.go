package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"

	"github.com/google/uuid"
)

func (rec *activityRecord) toActivity() *domain.Activity {
	return &domain.Activity{
		ResourceMeta: domain.ResourceMeta{
			ID:        rec.ID.String(),
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			CreatedBy: domain.UserSummary{
				ID:    rec.UserID.String(),
				Name:  rec.UserName,
				Email: rec.UserEmail,
			},
		},
		TripID:      rec.TripID.String(),
		LocationID:  rec.LocationID.String(),
		Name:        rec.Name,
		Description: ptr.FromPtr(rec.Description),
		StartTime:   ptr.FromPtr(rec.StartTime),
		EndTime:     ptr.FromPtr(rec.EndTime),
		Category:    domain.ActivityCategory(ptr.FromPtr(rec.Category)),
		Cost:        ptr.FromPtr(rec.Cost),
		Currency:    rec.Currency,
		Date:        rec.Date,
	}
}

func activityToRecord(activity *domain.Activity) (*activityRecord, error) {
	var activityID uuid.UUID
	var err error

	if activity.ID != "" {
		activityID, err = uuid.Parse(activity.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid activity ID: %w", err)
		}
	}

	tripID, err := uuid.Parse(activity.TripID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for trip ID: %v", err)
	}

	locationID, err := uuid.Parse(activity.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for location ID: %v", err)
	}

	userID, err := uuid.Parse(activity.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID for creator ID: %v", err)
	}

	categoryStr := string(activity.Category)

	return &activityRecord{
		ID:          activityID,
		TripID:      tripID,
		LocationID:  locationID,
		Name:        activity.Name,
		Description: ptr.ToPtrNonEmpty(activity.Description),
		StartTime:   ptr.ToPtr(activity.StartTime),
		EndTime:     ptr.ToPtr(activity.EndTime),
		UserID:      userID,
		UserName:    activity.CreatedBy.Name,
		UserEmail:   activity.CreatedBy.Email,
		Category:    ptr.ToPtrNonEmpty(categoryStr),
		Cost:        ptr.ToPtr(activity.Cost),
		Currency:    activity.Currency,
		Date:        activity.Date,
	}, nil
}
