package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
)

// LoginHandler handles POST /api/auth/login
func LoginHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req generated.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("Login: email=%s", req.Email)

		authResp, err := app.Services.Auth.Login(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, authResp)
	}
}

// GetMeHandler handles GET /api/users/me
func GetMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		app.Logger.Printf("GetMe: userId=%s", userId)

		user, err := app.Services.User.GetUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, user)
	}
}

// UpdateMeHandler handles PUT /api/users/me
func UpdateMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context (middleware will add it)
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		var req generated.UpdateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("UpdateMe: userId=%s", userId)

		user, err := app.Services.User.UpdateUser(r.Context(), userId, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, user)
	}
}

// ChangePasswordHandler handles PUT /api/users/me/password
func ChangePasswordHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context (middleware will add it)
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		var req generated.ChangePasswordRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("ChangePassword: userId=%s", userId)

		err := app.Services.Auth.ChangePassword(r.Context(), userId, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetUserHandler handles GET /api/users/{userId}
func GetUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		app.Logger.Printf("GetUser: id=%s", userId)

		user, err := app.Services.User.GetUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, user)
	}
}

// CreateUserHandler handles POST /api/auth/register
func CreateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req generated.CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("Register: email=%s, name=%s", req.Email, req.Name)

		authResp, err := app.Services.Auth.Register(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, authResp)
	}
}

// UpdateUserHandler handles PUT /api/users/{userId}
func UpdateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
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

		respondJSON(w, http.StatusOK, user)
	}
}

// DeleteUserHandler handles DELETE /api/users/{userId}
func DeleteUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")
		if userId == "" {
			respondError(w, http.StatusBadRequest, "User ID is required")
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
