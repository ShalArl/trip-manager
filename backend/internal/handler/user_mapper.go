package handler

import (
	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/storage"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapUserToUserResponse(user *domain.User, storage *storage.Storage) *generated.UserResponse {
	id, _ := uuid.Parse(user.ID)

	// Set avatarUrl from domain, convert to pointer if not empty
	var avatarUrl *string
	if user.AvatarURL != "" {
		avatarUrl = ptr.ToPtr(user.AvatarURL)
	}

	return &generated.UserResponse{
		AvatarUrl: avatarUrl,
		Bio:       &user.Bio,
		CreatedAt: &user.CreatedAt,
		Email:     openapitypes.Email(user.Email),
		Id:        ptr.ToPtr(id),
		Name:      user.Name,
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
