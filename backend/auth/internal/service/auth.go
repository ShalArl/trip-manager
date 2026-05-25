package service

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// TokenValidationResponse is returned after validating a token
type TokenValidationResponse struct {
	Valid   bool              `json:"valid"`
	UserID  string            `json:"userId,omitempty"`
	Email   string            `json:"email,omitempty"`
	Claims  map[string]interface{} `json:"claims,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// Service defines authentication operations
type Service interface {
	// ValidateToken validates a Bearer token and returns user info
	ValidateToken(ctx context.Context, token string) (*TokenValidationResponse, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	authClient *auth.Client
}

// NewService creates a new auth service
func NewService(authClient *auth.Client) Service {
	return &ServiceImpl{
		authClient: authClient,
	}
}

// ValidateToken validates a Firebase ID token
func (s *ServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenValidationResponse, error) {
	if token == "" {
		return &TokenValidationResponse{
			Valid: false,
			Error: "token is required",
		}, nil
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	// Verify the token
	decodedToken, err := s.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		return &TokenValidationResponse{
			Valid: false,
			Error: fmt.Sprintf("invalid token: %v", err),
		}, nil
	}

	// Return validation response with user info
	return &TokenValidationResponse{
		Valid:   true,
		UserID:  decodedToken.UID,
		Email:   decodedToken.Claims["email"].(string),
		Claims:  decodedToken.Claims,
	}, nil
}

