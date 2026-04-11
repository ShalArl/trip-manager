package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
)

// ListActivitiesForTripHandler handles GET /api/trips/{tripId}/activities with pagination
func ListActivitiesForTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		// Parse query parameters
		limit, offset := handlePaginationParams(r)

		app.Logger.Printf("ListActivitiesForTrip: tripId=%s, limit=%d, offset=%d", tripId, limit, offset)

		activities, totalCount, err := app.Services.Activity.ListActivitiesForTrip(r.Context(), limit, offset, tripId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		activitiesResp := mapActivitiesToActivityListResponse(activities, totalCount, limit, offset)

		respondJSON(w, http.StatusOK, activitiesResp)
	}
}

// ListActivitiesForLocationHandler handles GET /api/locations/{locationId}/activities with pagination
func ListActivitiesForLocationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locationId := r.PathValue("locationId")
		if locationId == "" {
			respondError(w, http.StatusBadRequest, "Location ID is required")
			return
		}

		// Parse query parameters
		limit, offset := handlePaginationParams(r)

		app.Logger.Printf("ListActivitiesForLocation: locationId=%s, limit=%d, offset=%d", locationId, limit, offset)

		activities, totalCount, err := app.Services.Activity.ListActivitiesForLocation(r.Context(), limit, offset, locationId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		activitiesResp := mapActivitiesToActivityListResponse(activities, totalCount, limit, offset)

		respondJSON(w, http.StatusOK, activitiesResp)
	}
}

// GetActivityHandler handles GET /api/activities/{activityId}
func GetActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityId := r.PathValue("activityId")
		if activityId == "" {
			respondError(w, http.StatusBadRequest, "Activity ID is required")
			return
		}

		app.Logger.Printf("GetActivity: id=%s", activityId)

		activity, err := app.Services.Activity.GetActivity(r.Context(), activityId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		activityResponse := mapActivityToActivityResponse(activity)

		respondJSON(w, http.StatusOK, activityResponse)
	}
}

// CreateActivityHandler handles POST /api/trips/{tripId}/activities
func CreateActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		var req generated.CreateActivityRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("CreateActivity: tripId=%s, name=%s", tripId, req.Name)

		activity, err := app.Services.Activity.CreateActivity(r.Context(), &req, tripId, "user-id-placeholder")
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		activityResponse := mapActivityToActivityResponse(activity)

		respondJSON(w, http.StatusCreated, activityResponse)
	}
}

// UpdateActivityHandler handles PUT /api/activities/{activityId}
func UpdateActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityId := r.PathValue("activityId")
		if activityId == "" {
			respondError(w, http.StatusBadRequest, "Activity ID is required")
			return
		}

		var req generated.UpdateActivityRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateActivity: id=%s", activityId)

		activity, err := app.Services.Activity.UpdateActivity(r.Context(), &req, activityId, "user-id-placeholder")
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		activityResponse := mapActivityToActivityResponse(activity)

		respondJSON(w, http.StatusOK, activityResponse)
	}
}

// DeleteActivityHandler handles DELETE /api/activities/{activityId}
func DeleteActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityId := r.PathValue("activityId")
		if activityId == "" {
			respondError(w, http.StatusBadRequest, "Activity ID is required")
			return
		}

		app.Logger.Printf("DeleteActivity: id=%s", activityId)

		err := app.Services.Activity.DeleteActivity(r.Context(), activityId, "user-id-placeholder")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
