package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

// ListLocationsHandler handles GET /api/trips/{tripId}/locations with pagination
func ListLocationsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		// Parse query parameters
		limit, offset := handlePaginationParams(r)

		app.Logger.Printf("ListLocations: tripId=%s, limit=%d, offset=%d", tripId, limit, offset)

		locations, total, err := app.Services.Location.ListLocations(r.Context(), tripId, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		locationResp := mapLocationsToLocationListResponse(locations, limit, offset, total)

		respondJSON(w, http.StatusOK, locationResp)
	}
}

// GetLocationHandler handles GET /api/locations/{locationId}
func GetLocationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locationId := r.PathValue("locationId")
		if locationId == "" {
			respondError(w, http.StatusBadRequest, "Location ID is required")
			return
		}

		app.Logger.Printf("GetLocation: id=%s", locationId)

		location, err := app.Services.Location.GetLocation(r.Context(), locationId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		locationResp := mapLocationToLocationResponse(location)

		respondJSON(w, http.StatusOK, locationResp)
	}
}

// CreateLocationHandler handles POST /api/trips/{tripId}/locations
func CreateLocationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		var req generated.CreateLocationRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("CreateLocation: tripId=%s, name=%s", tripId, req.Name)

		location, err := app.Services.Location.CreateLocation(r.Context(), &req, tripId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Printf("CreatedBy Email: %q\n", location.CreatedBy.Email) // ← hier

		locationResp := mapLocationToLocationResponse(location)

		respondJSON(w, http.StatusCreated, locationResp)
	}
}

// UpdateLocationHandler handles PUT /api/locations/{locationId}
func UpdateLocationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		locationId := r.PathValue("locationId")
		if locationId == "" {
			respondError(w, http.StatusBadRequest, "Location ID is required")
			return
		}

		var req generated.UpdateLocationRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateLocation: id=%s", locationId)

		location, err := app.Services.Location.UpdateLocation(r.Context(), &req, locationId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		locationResp := mapLocationToLocationResponse(location)

		respondJSON(w, http.StatusOK, locationResp)
	}
}

// DeleteLocationHandler handles DELETE /api/locations/{locationId}
func DeleteLocationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		locationId := r.PathValue("locationId")
		if locationId == "" {
			respondError(w, http.StatusBadRequest, "Location ID is required")
			return
		}

		app.Logger.Printf("DeleteLocation: id=%s", locationId)

		err = app.Services.Location.DeleteLocation(r.Context(), locationId, userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
