package service

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
)

func validateCreateUserRequest(request generated.CreateUserRequest) error {
	if request.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}

	if request.Email == "" {
		return fmt.Errorf("%w: email is required", domain.ErrInvalidInput)
	}

	if request.Password == "" {
		return fmt.Errorf("%w: password is required", domain.ErrInvalidInput)
	}

	return validatePassword(request.Password)
}

func validateUpdateUserRequest(request generated.UpdateUserRequest) error {
	if request.Name != nil && *request.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidInput)
	}

	if request.Email != nil && *request.Email == "" {
		return fmt.Errorf("%w: email cannot be empty", domain.ErrInvalidInput)
	}

	return nil
}

// validatePassword checks if password meets minimum requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters long", domain.ErrInvalidInput)
	}

	// Check for at least one uppercase, one lowercase, and one digit
	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, r := range password {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		}
		if r >= 'a' && r <= 'z' {
			hasLower = true
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("%w: password must contain uppercase, lowercase, and digit", domain.ErrInvalidInput)
	}

	return nil
}
