package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	DEFAULT_LIMIT  = 50
	DEFAULT_OFFSET = 0
)

// respondError writes an error JSON response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// respondJSON writes a JSON response with the given status code
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// handlePaginationParams parses the "limit" and "offset" query parameters, applying defaults if they are not provided or invalid
func handlePaginationParams(r *http.Request) (int, int) {
	limit := DEFAULT_LIMIT
	offset := DEFAULT_OFFSET

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
