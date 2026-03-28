package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ShalArl/trip-manager/internal/auth"
)

// UserIDContextKey is the key for storing user ID in context
type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthMiddleware validates JWT tokens and extracts user ID
func AuthMiddleware(authManager *auth.AuthManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Verify token
			claims, err := authManager.VerifyToken(tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID from request context
func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

