package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ShalArl/trip-manager/internal/auth"
)

type contextKey string

const (
	firebaseUIDKey contextKey = "firebaseUID"
	userIDKey      contextKey = "userID"
	emailKey       contextKey = "email"
	nameKey        contextKey = "name"
)

// UserResolver wandelt eine Firebase-UID in die interne Postgres-UUID um.
// Wird vom UserService implementiert.
type UserResolver interface {
	ResolveByFirebaseUID(ctx context.Context, firebaseUID string) (string, error)
}

func FirebaseAuthMiddleware(fb *auth.FirebaseAuth, resolver UserResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			idToken := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := fb.VerifyToken(r.Context(), idToken)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			email, _ := token.Claims["email"].(string)

			userID, err := resolver.ResolveByFirebaseUID(r.Context(), token.UID)
			if err != nil {
				http.Error(w, "user not provisioned", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), firebaseUIDKey, token.UID)
			ctx = context.WithValue(ctx, userIDKey, userID)
			ctx = context.WithValue(ctx, emailKey, email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalFirebaseAuthMiddleware(fb *auth.FirebaseAuth, resolver UserResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// Kein Token — einfach weiterreichen, kein 401
				next.ServeHTTP(w, r)
				return
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := fb.VerifyToken(r.Context(), idToken)
			if err != nil {
				// Token ungültig — als anonym behandeln, nicht 401 werfen.
				// Frontend könnte ein abgelaufenes Token mitschicken, für public
				// Routes ist das kein Problem.
				next.ServeHTTP(w, r)
				return
			}

			email, _ := token.Claims["email"].(string)

			userID, err := resolver.ResolveByFirebaseUID(r.Context(), token.UID)
			if err != nil {
				// User in Firebase vorhanden aber noch nicht provisioniert.
				// Request geht durch, aber ohne userID im Context.
				ctx := context.WithValue(r.Context(), firebaseUIDKey, token.UID)
				ctx = context.WithValue(ctx, emailKey, email)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := context.WithValue(r.Context(), firebaseUIDKey, token.UID)
			ctx = context.WithValue(ctx, userIDKey, userID)
			ctx = context.WithValue(ctx, emailKey, email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ProvisionMiddleware nur Token-Check ohne DB-Lookup. Für /provision-Endpoint.
func ProvisionMiddleware(fb *auth.FirebaseAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			idToken := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := fb.VerifyToken(r.Context(), idToken)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			email, _ := token.Claims["email"].(string)
			name, _ := token.Claims["name"].(string)

			ctx := context.WithValue(r.Context(), firebaseUIDKey, token.UID)
			ctx = context.WithValue(ctx, emailKey, email)
			ctx = context.WithValue(ctx, "name", name)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(userIDKey).(string)
	return id, ok
}

func GetFirebaseUID(r *http.Request) (string, bool) {
	uid, ok := r.Context().Value(firebaseUIDKey).(string)
	return uid, ok
}

func GetEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(emailKey).(string)
	return email, ok
}

func GetName(r *http.Request) (string, bool) {
	name, ok := r.Context().Value(nameKey).(string)
	return name, ok
}
