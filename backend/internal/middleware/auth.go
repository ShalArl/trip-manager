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

const (
	UserIDContextKey    contextKey = "userID"
	UserEmailContextKey contextKey = "userEmail"
	UserNameContextKey  contextKey = "userName"
)

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

			// Add user ID, email and name to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailContextKey, claims.Email)
			ctx = context.WithValue(ctx, UserNameContextKey, claims.Name)
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

// GetUserEmailFromContext extracts the user email from request context
func GetUserEmailFromContext(r *http.Request) (string, error) {
	email, ok := r.Context().Value(UserEmailContextKey).(string)
	if !ok || email == "" {
		return "", fmt.Errorf("user email not found in context")
	}
	return email, nil
}

// GetUserNameFromContext extracts the user name from request context
func GetUserNameFromContext(r *http.Request) (string, error) {
	name, ok := r.Context().Value(UserNameContextKey).(string)
	if !ok || name == "" {
		return "", fmt.Errorf("user name not found in context")
	}
	return name, nil
}
