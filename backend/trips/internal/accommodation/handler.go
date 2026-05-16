package accommodation

import (
	"encoding/json"
	"errors"
	"net/http"
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

// ── Request / Response types ──────────────────────────────────────────────────

type createRequest struct {
	Name       string     `json:"name"`
	Address    *string    `json:"address"`
	CheckIn    *time.Time `json:"checkIn"`
	CheckOut   *time.Time `json:"checkOut"`
	BookingRef *string    `json:"bookingRef"`
	Notes      *string    `json:"notes"`
}

type updateRequest struct {
	Name       *string    `json:"name"`
	Address    *string    `json:"address"`
	CheckIn    *time.Time `json:"checkIn"`
	CheckOut   *time.Time `json:"checkOut"`
	BookingRef *string    `json:"bookingRef"`
	Notes      *string    `json:"notes"`
}

type accommodationResponse struct {
	ID         string     `json:"id"`
	TripID     string     `json:"tripId"`
	Name       string     `json:"name"`
	Address    *string    `json:"address"`
	CheckIn    *time.Time `json:"checkIn"`
	CheckOut   *time.Time `json:"checkOut"`
	BookingRef *string    `json:"bookingRef"`
	Notes      *string    `json:"notes"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func toResponse(a *Accommodation) accommodationResponse {
	return accommodationResponse{
		ID:         a.ID,
		TripID:     a.TripID,
		Name:       a.Name,
		Address:    a.Address,
		CheckIn:    a.CheckIn,
		CheckOut:   a.CheckOut,
		BookingRef: a.BookingRef,
		Notes:      a.Notes,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
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
		accommodations, err := svc.ListByTrip(r.Context(), tripID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		data := make([]accommodationResponse, len(accommodations))
		for i, a := range accommodations {
			data[i] = toResponse(a)
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
		a, err := svc.Create(r.Context(), CreateInput{
			TripID:     tripID,
			Name:       req.Name,
			Address:    req.Address,
			CheckIn:    req.CheckIn,
			CheckOut:   req.CheckOut,
			BookingRef: req.BookingRef,
			Notes:      req.Notes,
		})
		if err != nil {
			if errors.Is(err, ErrInvalidInput) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, toResponse(a))
	}
}

func UpdateHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accommodationID := r.PathValue("accommodationId")
		if accommodationID == "" {
			respondError(w, http.StatusBadRequest, "accommodationId is required")
			return
		}
		var req updateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		a, err := svc.Update(r.Context(), UpdateInput{
			ID:         accommodationID,
			Name:       req.Name,
			Address:    req.Address,
			CheckIn:    req.CheckIn,
			CheckOut:   req.CheckOut,
			BookingRef: req.BookingRef,
			Notes:      req.Notes,
		})
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, "accommodation not found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, toResponse(a))
	}
}

func DeleteHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accommodationID := r.PathValue("accommodationId")
		if accommodationID == "" {
			respondError(w, http.StatusBadRequest, "accommodationId is required")
			return
		}
		if err := svc.Delete(r.Context(), accommodationID); err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, "accommodation not found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
