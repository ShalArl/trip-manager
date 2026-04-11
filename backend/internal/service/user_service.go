package service

import (
	"context"
	"fmt"
	"io"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/repository"
)

type UserService interface {
	// GetUser retrieves a user by their ID.
	GetUser(ctx context.Context, id string) (*domain.User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// CreateUser creates a new user with the provided details.
	CreateUser(ctx context.Context, request *generated.CreateUserRequest) (*domain.User, error)

	// UpdateUser updates an existing user's details.
	UpdateUser(ctx context.Context, id string, request *generated.UpdateUserRequest) (*domain.User, error)

	// UpdateUserWithAvatar updates user and optionally handles avatar upload
	UpdateUserWithAvatar(ctx context.Context, id string, request *generated.UpdateUserRequest, avatarFile io.Reader, avatarFileName string) (*domain.User, error)

	// UpdateUserPassword only called internally! Updates an existing user's details without validation (used by AuthService to update password)
	UpdateUserPassword(ctx context.Context, user *domain.User) (*domain.User, error)

	// DeleteUser removes a user from the system by their ID.
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	mediaService   *MediaService
}

// NewUserService creates a new UserService with optional MediaService
func NewUserService(userRepo repository.UserRepository, mediaService *MediaService) UserService {
	return &UserServiceImpl{
		userRepository: userRepo,
		mediaService:   mediaService,
	}
}

// CreateUser implements [UserService].
func (u *UserServiceImpl) CreateUser(ctx context.Context, request *generated.CreateUserRequest) (*domain.User, error) {
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
	updatedUser, err := u.userRepository.UpdateUserProfile(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// UpdateUserWithAvatar updates user profile and optionally handles avatar upload
func (u *UserServiceImpl) UpdateUserWithAvatar(ctx context.Context, id string, request *generated.UpdateUserRequest, avatarFile io.Reader, avatarFileName string) (*domain.User, error) {
	// If avatar file is provided and MediaService is available, upload it
	if avatarFile != nil && u.mediaService != nil && avatarFileName != "" {
		// Upload avatar via MediaService
		fileUrl, err := u.mediaService.UploadImage(ctx, avatarFile, UploadImageOptions{
			MediaType: MediaTypeAvatar,
			UserID:    id,
			FileName:  avatarFileName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload avatar: %w", err)
		}

		// Set avatar URL in request
		request.AvatarUrl = &fileUrl
	}

	// Now update the user with the potentially updated avatar URL
	return u.UpdateUser(ctx, id, request)
}

// UpdateUserPassword called only internally by AuthService therefore no validation
func (u *UserServiceImpl) UpdateUserPassword(ctx context.Context, user *domain.User) (*domain.User, error) {
	updatedUser, err := u.userRepository.UpdateUserPassword(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

