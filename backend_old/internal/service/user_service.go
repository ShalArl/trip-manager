package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type UserService interface {
	// GetUser retrieves a user by their ID.
	GetUser(ctx context.Context, id string) (*domain.User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateUser updates an existing user's details.
	UpdateUser(ctx context.Context, id string, request *generated.UpdateUserRequest) (*domain.User, error)

	// DeleteUser removes a user from the system by their ID.
	DeleteUser(ctx context.Context, id string) error

	// ResolveByFirebaseUID retrieves user details for a given firebase uid
	ResolveByFirebaseUID(ctx context.Context, firebaseUID string) (string, error)

	// ProvisionUser for a given firebase uid + requires user details to be stored in db
	ProvisionUser(ctx context.Context, firebaseUID, email, name string) (*domain.User, bool, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	mediaService   infrastructure.MediaService
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository, mediaService infrastructure.MediaService) UserService {
	return &UserServiceImpl{
		userRepository: userRepo,
		mediaService:   mediaService,
	}
}

// CreateUser implements [UserService].
func (u *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	log.Default().Printf("Creating user with email: %s", user.Email)
	createdUser, err := u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// DeleteUser implements [UserService].
func (u *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	return u.userRepository.DeleteUser(ctx, id)
}

// GetUser implements [UserService].
func (u *UserServiceImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (u *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil

}

// UpdateUser implements [UserService].
func (u *UserServiceImpl) UpdateUser(ctx context.Context, id string, request *generated.UpdateUserRequest) (*domain.User, error) {
	if err := validateUpdateUserRequest(*request); err != nil {
		return nil, err
	}

	// Verify AvatarKey if set
	if request.AvatarKey != nil && *request.AvatarKey != "" {
		// Sanity-Check: Does it belong to the current User?
		expectedPrefix := fmt.Sprintf("avatars/%s", id)
		if !strings.HasPrefix(*request.AvatarKey, expectedPrefix) {
			return nil, fmt.Errorf("%w: avatar key does not belong to user", domain.ErrInvalidInput)
		}

		// Does the file exist?
		exists, err := u.mediaService.ConfirmUpload(ctx, *request.AvatarKey)
		if err != nil {
			return nil, fmt.Errorf("verify avatar upload: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("%w: avatar not uploaded", domain.ErrInvalidInput)
		}
	}

	existingUser, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := mapUpdateUserRequestToUser(request, existingUser)

	updatedUser, err := u.userRepository.UpdateUserProfile(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// ResolveByFirebaseUID implements [UserService] + [UserResolver]
func (u *UserServiceImpl) ResolveByFirebaseUID(ctx context.Context, firebaseUID string) (string, error) {
	user, err := u.userRepository.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// ProvisionUser implements [UserService] + [UserResolver]
func (u *UserServiceImpl) ProvisionUser(ctx context.Context, firebaseUID, email, name string) (*domain.User, bool, error) {
	existing, err := u.userRepository.GetUserByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, false, fmt.Errorf("lookup existing user: %w", err)
	}

	user := &domain.User{
		FirebaseUID: firebaseUID,
		Email:       email,
		Name:        name,
	}
	created, err := u.userRepository.CreateUser(ctx, user)
	if err != nil {
		println("failed to create user")
		return nil, false, err
	}
	return created, true, nil
}
