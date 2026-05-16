package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func getToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// ── Request / Response types ──────────────────────────────────────────────────

type createRequest struct {
	Type           string     `json:"type"`
	DeparturePlace string     `json:"departurePlace"`
	ArrivalPlace   string     `json:"arrivalPlace"`
	DepartureTime  *time.Time `json:"departureTime"`
	ArrivalTime    *time.Time `json:"arrivalTime"`
	BookingRef     *string    `json:"bookingRef"`
	Notes          *string    `json:"notes"`
}

type updateRequest struct {
	Type           *string    `json:"type"`
	DeparturePlace *string    `json:"departurePlace"`
	ArrivalPlace   *string    `json:"arrivalPlace"`
	DepartureTime  *time.Time `json:"departureTime"`
	ArrivalTime    *time.Time `json:"arrivalTime"`
	BookingRef     *string    `json:"bookingRef"`
	Notes          *string    `json:"notes"`
}

type transportResponse struct {
	ID             string     `json:"id"`
	TripID         string     `json:"tripId"`
	Type           string     `json:"type"`
	DeparturePlace string     `json:"departurePlace"`
	ArrivalPlace   string     `json:"arrivalPlace"`
	DepartureTime  *time.Time `json:"departureTime"`
	ArrivalTime    *time.Time `json:"arrivalTime"`
	BookingRef     *string    `json:"bookingRef"`
	Notes          *string    `json:"notes"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func toResponse(t *Transport) transportResponse {
	return transportResponse{
		ID:             t.ID,
		TripID:         t.TripID,
		Type:           t.Type,
		DeparturePlace: t.DeparturePlace,
		ArrivalPlace:   t.ArrivalPlace,
		DepartureTime:  t.DepartureTime,
		ArrivalTime:    t.ArrivalTime,
		BookingRef:     t.BookingRef,
		Notes:          t.Notes,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func ListHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "tripId is required")
			return
		}
		transports, err := svc.ListByTrip(r.Context(), tripID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		data := make([]transportResponse, len(transports))
		for i, t := range transports {
			data[i] = toResponse(t)
		}
		respondJSON(w, http.StatusOK, map[string]any{"data": data})
	}
}

func CreateHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("tripId")
		if tripID == "" {
			respondError(w, http.StatusBadRequest, "tripId is required")
			return
		}
		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		t, err := svc.Create(r.Context(), CreateInput{
			TripID:         tripID,
			Type:           req.Type,
			DeparturePlace: req.DeparturePlace,
			ArrivalPlace:   req.ArrivalPlace,
			DepartureTime:  req.DepartureTime,
			ArrivalTime:    req.ArrivalTime,
			BookingRef:     req.BookingRef,
			Notes:          req.Notes,
		})
		if err != nil {
			if errors.Is(err, ErrInvalidInput) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, toResponse(t))
	}
}

func UpdateHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transportID := r.PathValue("transportId")
		if transportID == "" {
			respondError(w, http.StatusBadRequest, "transportId is required")
			return
		}
		var req updateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		t, err := svc.Update(r.Context(), UpdateInput{
			ID:             transportID,
			Type:           req.Type,
			DeparturePlace: req.DeparturePlace,
			ArrivalPlace:   req.ArrivalPlace,
			DepartureTime:  req.DepartureTime,
			ArrivalTime:    req.ArrivalTime,
			BookingRef:     req.BookingRef,
			Notes:          req.Notes,
		})
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, "transport not found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, toResponse(t))
	}
}

func DeleteHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transportID := r.PathValue("transportId")
		if transportID == "" {
			respondError(w, http.StatusBadRequest, "transportId is required")
			return
		}
		if err := svc.Delete(r.Context(), transportID); err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, "transport not found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
