package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ShalArl/trip-manager/backend/weather-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/fetcher"
)

type Handler struct {
	cache   *cache.WeatherCache
	fetcher *fetcher.Client
}

func NewHandler(cache *cache.WeatherCache, fetcher *fetcher.Client) *Handler {
	return &Handler{
		cache:   cache,
		fetcher: fetcher,
	}
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		respondError(w, http.StatusBadRequest, "lat and lng are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid lat")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid lng")
		return
	}

	// Cache prüfen
	weather, err := h.cache.Get(r.Context(), lat, lng)
	if err != nil {
		log.Printf("cache get error: %v", err)
	}

	// Cache Hit
	if weather != nil {
		respondJSON(w, http.StatusOK, weather)
		return
	}

	// Cache Miss – On-Demand von Open-Meteo holen
	log.Printf("cache miss for lat=%.2f lng=%.2f – fetching from Open-Meteo", lat, lng)
	weather, err = h.fetcher.FetchForecast(r.Context(), lat, lng)
	if err != nil {
		log.Printf("fetch forecast error: %v", err)
		respondError(w, http.StatusServiceUnavailable, "failed to fetch weather data")
		return
	}

	// In Cache schreiben – Fehler nur loggen, nicht an Client weitergeben
	if err := h.cache.Set(r.Context(), weather); err != nil {
		log.Printf("cache set error: %v", err)
	}

	respondJSON(w, http.StatusOK, weather)
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
