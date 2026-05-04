package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

func ListTransportsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripId := r.PathValue("tripId")
		if tripId == "" {
			respondError(w, http.StatusBadRequest, "Trip ID is required")
			return
		}

		limit, offset := handlePaginationParams(r)

		app.Logger.Printf("ListTransports: tripId=%s, limit=%d, offset=%d", tripId, limit, offset)

		transports, total, err := app.Services.Transport.ListTransports(r.Context(), tripId, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, mapTransportsToTransportListResponse(transports, limit, offset, total))
	}
}

func GetTransportHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transportId := r.PathValue("transportId")
		if transportId == "" {
			respondError(w, http.StatusBadRequest, "Transport ID is required")
			return
		}

		app.Logger.Printf("GetTransport: id=%s", transportId)

		transport, err := app.Services.Transport.GetTransport(r.Context(), transportId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, mapTransportToTransportResponse(transport))
	}
}

func CreateTransportHandler(app *app.App) http.HandlerFunc {
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

		var req generated.CreateTransportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("CreateTransport: tripId=%s", tripId)

		transport, err := app.Services.Transport.CreateTransport(r.Context(), &req, tripId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, mapTransportToTransportResponse(transport))
	}
}

func UpdateTransportHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		transportId := r.PathValue("transportId")
		if transportId == "" {
			respondError(w, http.StatusBadRequest, "Transport ID is required")
			return
		}

		var req generated.UpdateTransportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateTransport: id=%s", transportId)

		transport, err := app.Services.Transport.UpdateTransport(r.Context(), &req, transportId, userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, mapTransportToTransportResponse(transport))
	}
}

func DeleteTransportHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		transportId := r.PathValue("transportId")
		if transportId == "" {
			respondError(w, http.StatusBadRequest, "Transport ID is required")
			return
		}

		app.Logger.Printf("DeleteTransport: id=%s", transportId)

		err := app.Services.Transport.DeleteTransport(r.Context(), transportId, userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
