package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
)

func mapActivityCreateRequestToActivity(request *generated.CreateActivityRequest, tripID string, userID string, userName string, userEmail string) *domain.Activity {
	category := ""
	if request.Category != nil {
		category = string(*request.Category)
	}

	// Cost ist im Request ein *float32, im Domain oft float64 oder float32
	var cost float32
	if request.Cost != nil {
		cost = *request.Cost
	}

	return &domain.Activity{
		ResourceMeta: domain.ResourceMeta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: domain.UserSummary{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
			},
		},
		TripID:     tripID,
		LocationID: request.LocationId.String(),

		Name:        request.Name,
		Description: ptr.FromPtr(request.Description),
		Date:        request.Date.Time,

		StartTime: ptr.FromPtr(request.StartTime),
		EndTime:   ptr.FromPtr(request.EndTime),

		Category: domain.ActivityCategory(category),
		Cost:     float64(cost),
		Currency: ptr.FromPtr(request.Currency),
	}
}

func mapActivityUpdateRequestToActivity(request *generated.UpdateActivityRequest, existing *domain.Activity) *domain.Activity {
	updated := *existing

	if request.Name != nil {
		updated.Name = *request.Name
	}

	if request.Description != nil {
		updated.Description = *request.Description
	}

	if request.Date != nil {
		updated.Date = request.Date.Time
	}

	if request.StartTime != nil {
		updated.StartTime = *request.StartTime
	}

	if request.EndTime != nil {
		updated.EndTime = *request.EndTime
	}

	if request.Category != nil {
		updated.Category = domain.ActivityCategory(*request.Category)
	}

	if request.Cost != nil {
		updated.Cost = float64(*request.Cost)
	}

	if request.Currency != nil {
		updated.Currency = *request.Currency
	}

	if request.LocationId != nil {
		updated.LocationID = request.LocationId.String()
	}

	updated.UpdatedAt = time.Now()

	return &updated
}
