package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapCreateUserRequestToUser(request *generated.CreateUserRequest) *domain.User {
	return &domain.User{
		ID:        uuid.New().String(),
		Email:     string(request.Email),
		Name:      request.Name,
		Bio:       "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func mapUpdateUserRequestToUser(request *generated.UpdateUserRequest, existing *domain.User) *domain.User {
	updated := *existing

	if request.Name != nil {
		updated.Name = *request.Name
	}

	if request.Email != nil {
		updated.Email = string(*request.Email)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

func mapUserToUserResponse(user *domain.User) *generated.UserResponse {
	id, _ := uuid.Parse(user.ID)

	return &generated.UserResponse{
		Id:        ptr.ToPtr(id),
		Email:     openapitypes.Email(user.Email),
		Name:      user.Name,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
}

func mapUserToUserSummary(user *domain.User) *generated.UserSummary {
	id, _ := uuid.Parse(user.ID)

	return &generated.UserSummary{
		Id:    id,
		Email: openapitypes.Email(user.Email),
		Name:  user.Name,
	}
}
