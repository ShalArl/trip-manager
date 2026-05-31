package newsletter

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
)

func GetNewsletterHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authclient.GetUserID(r)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		limit := 10
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
				limit = parsed
			}
		}

		newsletter, err := svc.GetNewsletter(r.Context(), userID, limit)
		if err != nil {
			log.Printf("newsletter: GetNewsletter error: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newsletter)
	}
}
