package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/middleware"
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

		token, exp, user, err := app.Services.Auth.Login(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userResponse := mapUserToUserResponse(user, nil)

		authResp := generated.AuthResponse{
			ExpiresIn: exp,
			Token:     token,
			User:      *userResponse,
		}

		respondJSON(w, http.StatusOK, authResp)
	}
}

// GetMeHandler handles GET /api/users/me
func GetMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		app.Logger.Printf("GetMe: userId=%s", userId)

		user, err := app.Services.User.GetUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(user, nil)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// UpdateMeHandler handles PUT /api/users/me
// Expects JSON request body with UpdateUserRequest
// For avatar uploads, use POST /api/uploads/presigned to get a presigned URL first
func UpdateMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var req generated.UpdateUserRequest

		// Handle regular JSON only (multipart uploads no longer supported)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("[Handler] UpdateMe: userId=%s, name=%s, email=%s", userId, req.Name, req.Email)

		user, err := app.Services.User.UpdateUser(r.Context(), userId, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(user, nil)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// ChangePasswordHandler handles PUT /api/users/me/password
func ChangePasswordHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var req generated.ChangePasswordRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		app.Logger.Printf("ChangePassword: userId=%s", userId)

		err = app.Services.Auth.ChangePassword(r.Context(), userId, &req)
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
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		app.Logger.Printf("GetUser: id=%s", userId)

		user, err := app.Services.User.GetUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		userResponse := mapUserToUserResponse(user, nil)

		respondJSON(w, http.StatusOK, userResponse)
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

		token, exp, user, err := app.Services.Auth.Register(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		userResponse := mapUserToUserResponse(user, nil)
		authResp := generated.AuthResponse{
			ExpiresIn: exp,
			Token:     token,
			User:      *userResponse,
		}

		respondJSON(w, http.StatusCreated, authResp)
	}
}

// UpdateUserHandler handles PUT /api/users/{userId}
func UpdateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
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

		userResponse := mapUserToUserResponse(user, nil)

		respondJSON(w, http.StatusOK, userResponse)
	}
}

// DeleteUserHandler handles DELETE /api/users/{userId}
func DeleteUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		app.Logger.Printf("DeleteUser: id=%s", userId)

		err = app.Services.User.DeleteUser(r.Context(), userId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
