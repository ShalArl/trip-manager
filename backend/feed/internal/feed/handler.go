package feed

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	generated "github.com/ShalArl/trip-manager/backend/feed/generated"
	"github.com/ShalArl/trip-manager/backend/shared/authclient"
)

func GetFeedHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := queryInt(r, "limit", 20)
		offset := queryInt(r, "offset", 0)

		// Firebase UID – leer wenn nicht eingeloggt
		userID, _ := authclient.GetUserID(r)

		trips, total, err := svc.GetFeed(r.Context(), userID, limit, offset)
		if err != nil {
			log.Printf("feed: error getting feed: %v", err)
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if trips == nil {
			trips = []generated.FeedTrip{}
		}

		respondJSON(w, http.StatusOK, generated.FeedResponse{
			Data:   trips,
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

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

func queryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}
