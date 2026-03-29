package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

type APIErrorResponse struct {
	Message string `json:"message"`
}

// ListTripsHandler handles GET /api/trips with pagination
func ListTripsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		limit, offset := handlePaginationParams(r)

		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		app.Logger.Printf("ListTrips: limit=%d, offset=%d", limit, offset)

		// Handler only parses parameters - Service does validation + coordination
		tripsResp, err := app.Services.Trip.ListTrips(r.Context(), userID, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			app.Logger.Printf("ListTrips error: %v", err)
			return
		}

		app.Logger.Printf("ListTrips response: %+v", tripsResp)
		respondJSON(w, http.StatusOK, tripsResp)
	}
}

// CreateTripHandler handles POST /api/trips
func CreateTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req generated.CreateTripRequest

		// Handler only decodes JSON - validation belongs in Service layer
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("CreateTrip: title=%s, startDate=%v, endDate=%v", req.Title, req.StartDate, req.EndDate)

		// Delegate to Service layer (includes validation + business logic)
		/*trip, err := app.Services.Trip.CreateTrip(r.Context(), &req, "user-id-placeholder")
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}*/
		userID, err := middleware.GetUserIDFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userName, err := middleware.GetUserNameFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userEmail, err := middleware.GetUserEmailFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		trip, err := app.Services.Trip.CreateTrip(r.Context(), &req, userID, userName, userEmail)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, trip)
	}
}

// GetTripHandler handles GET /api/trips/{tripId}
func GetTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		app.Logger.Printf("GetTrip: id=%s", tripId)

		trip, err := app.Services.Trip.GetTrip(r.Context(), tripId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, trip)
	}
}

// UpdateTripHandler handles PUT /api/trips/{tripId}
func UpdateTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		var req generated.UpdateTripRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateTrip: id=%s", tripId)

		trip, err := app.Services.Trip.UpdateTrip(r.Context(), &req, tripId)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, trip)
	}
}

// DeleteTripHandler handles DELETE /api/trips/{tripId}
func DeleteTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		app.Logger.Printf("DeleteTrip: id=%s", tripId)

		err := app.Services.Trip.DeleteTrip(r.Context(), tripId, "user-id-placeholder")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
