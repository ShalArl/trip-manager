package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

func ListAccommodationsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}
		limit, offset := handlePaginationParams(r)
		app.Logger.Printf("ListAccommodations: tripId=%s, limit=%d, offset=%d", tripId, limit, offset)
		accommodations, total, err := app.Services.Accommodation.ListAccommodations(r.Context(), tripId, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, mapAccommodationsToAccommodationListResponse(accommodations, limit, offset, total))
	}
}

func CreateAccommodationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}
		var req generated.CreateAccommodationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}
		app.Logger.Printf("CreateAccommodation: tripId=%s", tripId)
		accommodation, err := app.Services.Accommodation.CreateAccommodation(r.Context(), &req, tripId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, mapAccommodationToAccommodationResponse(accommodation))
	}
}

func UpdateAccommodationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		accommodationId := r.PathValue("accommodationId")
		if accommodationId == "" {
			respondError(w, http.StatusBadRequest, "Accommodation ID is required")
			return
		}
		var req generated.UpdateAccommodationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}
		app.Logger.Printf("UpdateAccommodation: id=%s", accommodationId)
		accommodation, err := app.Services.Accommodation.UpdateAccommodation(r.Context(), &req, accommodationId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, mapAccommodationToAccommodationResponse(accommodation))
	}
}

func DeleteAccommodationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		accommodationId := r.PathValue("accommodationId")
		if accommodationId == "" {
			respondError(w, http.StatusBadRequest, "Accommodation ID is required")
			return
		}
		app.Logger.Printf("DeleteAccommodation: id=%s", accommodationId)
		err := app.Services.Accommodation.DeleteAccommodation(r.Context(), accommodationId, userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
