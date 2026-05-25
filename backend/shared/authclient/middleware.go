package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type contextKey string

const (
	contextKeyUserID     contextKey = "userId"
	contextKeyUserEmail  contextKey = "userEmail"
	contextKeyUserClaims contextKey = "userClaims"
)

func RequireAuth(client *Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "authorization header required")
				return
			}
			result, err := client.ValidateBearerToken(r.Context(), authHeader)
			if err != nil {
				respondError(w, http.StatusUnauthorized, err.Error())
				return
			}
			if !result.Valid {
				respondError(w, http.StatusUnauthorized, result.Error)
				return
			}
			ctx := context.WithValue(r.Context(), contextKeyUserID, result.UserID)
			ctx = context.WithValue(ctx, contextKeyUserEmail, result.Email)
			ctx = context.WithValue(ctx, contextKeyUserClaims, result.Claims)
			next(w, r.WithContext(ctx))
		}
	}
}

func OptionalAuth(client *Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			var userID, userEmail string
			if authHeader != "" {
				result, err := client.ValidateBearerToken(r.Context(), authHeader)
				if err == nil && result.Valid {
					userID = result.UserID
					userEmail = result.Email
				}
			}
			ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
			ctx = context.WithValue(ctx, contextKeyUserEmail, userEmail)
			next(w, r.WithContext(ctx))
		}
	}
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(contextKeyUserID).(string)
	return userID, ok && userID != ""
}

func GetUserEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(contextKeyUserEmail).(string)
	return email, ok && email != ""
}

func GetUserClaims(r *http.Request) map[string]interface{} {
	claims, ok := r.Context().Value(contextKeyUserClaims).(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return claims
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		fmt.Println(err)
	}
}
