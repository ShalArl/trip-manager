package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
)

// ProvisionMeHandler handles POST /api/users/provision
func ProvisionMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firebaseUID, ok := middleware.GetFirebaseUID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		email, _ := middleware.GetEmail(r)
		name, _ := middleware.GetName(r)

		var req generated.ProvisionUserRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Name != nil && *req.Name != "" {
			name = *req.Name
		}

		app.Logger.Printf("ProvisionUser %s, %s, %s", name, email, firebaseUID)

		user, created, err := app.Services.User.ProvisionUser(r.Context(), firebaseUID, email, name)
		if err != nil {
			app.Logger.Printf("[ProvisionMe] err=%v", err)
			respondError(w, http.StatusInternalServerError, "provisioning failed")
			return
		}

		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}
		respondJSON(w, status, mapUserToUserResponse(r.Context(), app.Services.Media, user))
	}
}

// GetMeHandler handles GET /api/users/me
func GetMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		app.Logger.Printf("GetMe: userId=%s", userID)

		user, err := app.Services.User.GetUser(r.Context(), userID)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(r.Context(), app.Services.Media, user)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// UpdateMeHandler handles PUT /api/users/me
func UpdateMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userID, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req generated.UpdateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("[Handler] UpdateMe: userId=%s, name=%v, email=%v", userID, req.Name, req.Email)

		user, err := app.Services.User.UpdateUser(r.Context(), userID, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(r.Context(), app.Services.Media, user)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// GetUserHandler handles GET /api/users/{userId}
func GetUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		app.Logger.Printf("GetUser: id=%s", userId)

		user, err := app.Services.User.GetUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(r.Context(), app.Services.Media, user)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// UpdateUserHandler handles PUT /api/users/{userId}
func UpdateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req generated.UpdateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateUser: id=%s", userId)

		user, err := app.Services.User.UpdateUser(r.Context(), userId, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		userResponse := mapUserToUserResponse(r.Context(), app.Services.Media, user)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// DeleteUserHandler handles DELETE /api/users/{userId}
func DeleteUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := middleware.GetUserID(r)
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		app.Logger.Printf("DeleteUser: id=%s", userId)

		err := app.Services.User.DeleteUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
