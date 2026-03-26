package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/auth"
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

// AuthService handles authentication operations
type AuthService interface {
	// Register creates a new user and returns an auth response
	Register(ctx context.Context, request *generated.CreateUserRequest) (*generated.AuthResponse, error)

	// Login authenticates a user and returns an auth response
	Login(ctx context.Context, request *generated.LoginRequest) (*generated.AuthResponse, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID string, request *generated.ChangePasswordRequest) error
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepository repository.UserRepository
	authManager    *auth.AuthManager
	userService    UserService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepository repository.UserRepository, authManager *auth.AuthManager, userService UserService) AuthService {
	return &AuthServiceImpl{
		userRepository: userRepository,
		authManager:    authManager,
		userService:    userService,
	}
}

// Register implements AuthService.Register
func (as *AuthServiceImpl) Register(ctx context.Context, request *generated.CreateUserRequest) (*generated.AuthResponse, error) {
	// 1. Validate the request
	if err := validateCreateUserRequest(*request); err != nil {
		return nil, err
	}

	// 2. Check if user already exists
	existingUser, err := as.userRepository.GetUserByEmail(ctx, string(request.Email))
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("%w: email already in use", domain.ErrConflict)
	}

	// 3. Hash the password
	hashedPassword, err := as.authManager.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to hash password", domain.ErrInternal)
	}

	// 4. Create the user domain object
	user := mapCreateUserRequestToUser(request)
	user.PasswordHash = hashedPassword

	// 5. Persist the user
	createdUser, err := as.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. Generate JWT token
	token, err := as.authManager.GenerateToken(createdUser.ID, createdUser.Email, createdUser.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 7. Return auth response
	userResp := mapUserToUserResponse(createdUser)
	return &generated.AuthResponse{
		Token:     token,
		ExpiresIn: as.authManager.GetTokenExpiresIn(),
		User:      *userResp,
	}, nil
}

// Login implements AuthService.Login
func (as *AuthServiceImpl) Login(ctx context.Context, request *generated.LoginRequest) (*generated.AuthResponse, error) {
	// 1. Fetch user by email
	user, err := as.userRepository.GetUserByEmail(ctx, string(request.Email))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", domain.ErrUnauthorized)
	}

	// 2. Verify password
	if !as.authManager.VerifyPassword(user.PasswordHash, request.Password) {
		return nil, fmt.Errorf("%w: invalid credentials", domain.ErrUnauthorized)
	}

	// 3. Generate JWT token
	token, err := as.authManager.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 4. Return auth response
	userResp := mapUserToUserResponse(user)
	return &generated.AuthResponse{
		Token:     token,
		ExpiresIn: as.authManager.GetTokenExpiresIn(),
		User:      *userResp,
	}, nil
}

// ChangePassword implements AuthService.ChangePassword
func (as *AuthServiceImpl) ChangePassword(ctx context.Context, userID string, request *generated.ChangePasswordRequest) error {
	// 1. Fetch the user
	user, err := as.userRepository.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Verify current password
	if !as.authManager.VerifyPassword(user.PasswordHash, request.CurrentPassword) {
		return fmt.Errorf("%w: current password is incorrect", domain.ErrUnauthorized)
	}

	// 3. Validate new password
	if err := validatePassword(request.NewPassword); err != nil {
		return err
	}

	// 4. Hash new password
	hashedPassword, err := as.authManager.HashPassword(request.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: failed to hash password", domain.ErrInternal)
	}

	// 5. Update the user
	user.PasswordHash = hashedPassword
	_, err = as.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}



