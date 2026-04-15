package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// ListTripsHandler handles GET /api/trips with pagination
func ListTripsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		// Parse query parameters
		limit, offset := handlePaginationParams(r)

		app.Logger.Printf("ListTrips: limit=%d, offset=%d", limit, offset)

		// Handler only parses parameters - Service does validation + coordination
		trips, total, err := app.Services.Trip.ListTrips(r.Context(), userID, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			app.Logger.Printf("ListTrips error: %v", err)
			return
		}

		tripsResponse := mapTripsToTripListResponse(trips, limit, offset, total)

		app.Logger.Printf("ListTrips response: %+v", tripsResponse)
		respondJSON(w, http.StatusOK, tripsResponse)
	}
}

// CreateTripHandler handles POST /api/trips
func CreateTripHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userEmail, userName, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req generated.CreateTripRequest
		// Handler only decodes JSON - validation belongs in Service layer
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("CreateTrip: title=%s, startDate=%v, endDate=%v", req.Title, req.StartDate, req.EndDate)

		trip, err := app.Services.Trip.CreateTrip(r.Context(), &req, userID, userName, userEmail)
		print("trip created: ", trip)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		tripResponse := mapTripToTripResponse(trip)

		respondJSON(w, http.StatusCreated, tripResponse)
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

		tripResponse := mapTripToTripResponse(trip)

		respondJSON(w, http.StatusOK, tripResponse)
	}
}

// UpdateTripHandler handles PUT /api/trips/{tripId}
func UpdateTripHandler(app *app.App) http.HandlerFunc {
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

		var req generated.UpdateTripRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateTrip: id=%s", tripId)

		trip, err := app.Services.Trip.UpdateTrip(r.Context(), &req, tripId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		tripResponse := mapTripToTripResponse(trip)

		respondJSON(w, http.StatusOK, tripResponse)
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

// SearchTripsHandler handles GET /api/trips/search?q=...
func SearchTripsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query holen
		query := r.URL.Query().Get("q")
		if query == "" {
			respondError(w, http.StatusBadRequest, "query is required")
			return
		}

		// Pagination
		limit, offset := handlePaginationParams(r)

		// Service call
		trips, total, err := app.Services.Trip.SearchTrips(r.Context(), query, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		respondJSON(w, http.StatusOK, mapTripsToTripListResponse(trips, limit, offset, total))
	}
}
