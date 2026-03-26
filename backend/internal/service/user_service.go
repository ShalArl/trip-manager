package service

import (
	"context"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type UserService interface {
	// GetUser retrieves a user by their ID.
	GetUser(ctx context.Context, id string) (*generated.UserResponse, error)

	// CreateUser creates a new user with the provided details.
	CreateUser(ctx context.Context, request *generated.CreateUserRequest) (*generated.UserResponse, error)

	// UpdateUser updates an existing user's details.
	UpdateUser(ctx context.Context, id string, request *generated.UpdateUserRequest) (*generated.UserResponse, error)

	// DeleteUser removes a user from the system by their ID.
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

// CreateUser implements [UserService].
func (u *UserServiceImpl) CreateUser(ctx context.Context, request *generated.CreateUserRequest) (*generated.UserResponse, error) {
	// 1. Validate input (business logic validation)
	if err := validateCreateUserRequest(*request); err != nil {
		return nil, err
	}

	// 2. Convert from generated type to domain
	user := mapCreateUserRequestToUser(request)

	// 3. Call repository to persist
	createdUser, err := u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 4. Convert from domain back to response type
	response := mapUserToUserResponse(createdUser)
	return response, nil
}

// DeleteUser implements [UserService].
func (u *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	return u.userRepository.DeleteUser(ctx, id)
}

// GetUser implements [UserService].
func (u *UserServiceImpl) GetUser(ctx context.Context, id string) (*generated.UserResponse, error) {
	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	response := mapUserToUserResponse(user)
	return response, nil
}

// UpdateUser implements [UserService].
func (u *UserServiceImpl) UpdateUser(ctx context.Context, id string, request *generated.UpdateUserRequest) (*generated.UserResponse, error) {
	// 1. Validate input (business logic validation)
	if err := validateUpdateUserRequest(*request); err != nil {
		return nil, err
	}

	// 2. Fetch existing user
	existingUser, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 3. Convert from generated type to domain
	user := mapUpdateUserRequestToUser(request, existingUser)

	// 4. Call repository to update and get updated record
	updatedUser, err := u.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// 5. Convert back to response type
	response := mapUserToUserResponse(updatedUser)
	return response, nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}
