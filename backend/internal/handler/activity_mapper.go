package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapActivityToActivityResponse(activity *domain.Activity) *generated.ActivityResponse {
	id, _ := uuid.Parse(activity.ID)
	locID, _ := uuid.Parse(activity.LocationID)

	category := generated.ActivityResponseCategory(activity.Category)

	return &generated.ActivityResponse{
		Id:         ptr.ToPtr(id),
		LocationId: locID,

		Name:        activity.Name,
		Description: ptr.ToPtr(activity.Description),

		Date:      openapitypes.Date{Time: activity.Date},
		StartTime: ptr.ToPtr(activity.StartTime),
		EndTime:   ptr.ToPtr(activity.EndTime),

		Category: &category,
		Cost:     ptr.ToPtr(float32(activity.Cost)),
		Currency: ptr.ToPtr(activity.Currency),

		CreatedAt: &activity.CreatedAt,
		UpdatedAt: &activity.UpdatedAt,

		CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(activity.CreatedBy.ID),
			Name:  activity.CreatedBy.Name,
			Email: openapitypes.Email(activity.CreatedBy.Email),
		},
	}
}

func mapActivitiesToActivityListResponse(activities []*domain.Activity, limit int, offset int, total int) *generated.ActivityListResponse {
	act := make([]generated.ActivityResponse, len(activities))
	for i, activity := range activities {
		act[i] = *mapActivityToActivityResponse(activity)
	}

	return &generated.ActivityListResponse{
		Data:   act,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
