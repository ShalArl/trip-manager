package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ShalArl/trip-manager/internal/app"
	"github.com/ShalArl/trip-manager/internal/domain"
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
// Supports optional avatar file upload via multipart/form-data
func UpdateMeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from JWT token in context
		userId, _, _, err := middleware.GetUserInfoFromContext(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Parse multipart form (supports both JSON and multipart)
		contentType := r.Header.Get("Content-Type")
		var req generated.UpdateUserRequest
		var avatarFile io.Reader
		var avatarFileName string

		if strings.HasPrefix(contentType, "multipart/form-data") {
			// Handle multipart form data with optional file
			if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
				respondError(w, http.StatusBadRequest, "Invalid form data")
				return
			}

			// Parse JSON fields from form
			if jsonData := r.FormValue("data"); jsonData != "" {
				if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
					respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid data field: %v", err))
					return
				}
			}

			// Handle optional avatar file
			file, header, err := r.FormFile("avatar")
			if err == nil {
				// File was provided
				defer file.Close()

				// Validate file type
				validImageTypes := map[string]bool{
					"image/jpeg": true,
					"image/png":  true,
					"image/gif":  true,
					"image/webp": true,
				}

				contentType := header.Header.Get("Content-Type")
				if !validImageTypes[contentType] {
					respondError(w, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
					return
				}

				// Validate file size
				if header.Size > MaxUploadSize {
					respondError(w, http.StatusBadRequest, "File is too large (max 5MB)")
					return
				}

				avatarFile = file
				avatarFileName = header.Filename
			}
		} else {
			// Handle regular JSON
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
				return
			}
		}

		app.Logger.Printf("UpdateMe: userId=%s", userId)

		// Use UserService to handle update (with optional avatar upload)
		var user *domain.User
		if avatarFile != nil {
			user, err = app.Services.User.UpdateUserWithAvatar(r.Context(), userId, &req, avatarFile, avatarFileName)
		} else {
			user, err = app.Services.User.UpdateUser(r.Context(), userId, &req)
		}

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
