package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/cache"
)

type Handler struct {
	cache *cache.WarningCache
}

func NewHandler(cache *cache.WarningCache) *Handler {
	return &Handler{cache: cache}
}

func (h *Handler) GetWarning(w http.ResponseWriter, r *http.Request) {
	countryCode := r.PathValue("countryCode")
	if countryCode == "" {
		respondError(w, http.StatusBadRequest, "invalid country code")
		return
	}

	warning, err := h.cache.Get(r.Context(), countryCode)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get warning")
		return
	}

	if warning == nil {
		respondError(w, http.StatusNotFound, "warning not found")
		return
	}

	respondJSON(w, http.StatusOK, warning)
}

// GetWarnings potentially implement in the future
func (h *Handler) GetWarnings(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "not yet implemented")
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
