package handler

import (
	"context"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapUserToUserResponse(ctx context.Context, media infrastructure.MediaService, user *domain.User) *generated.UserResponse {
	id, _ := uuid.Parse(user.ID)

	var avatarURL *string
	if user.AvatarKey != "" {
		if url, err := media.GetDownloadURL(ctx, user.AvatarKey); err == nil {
			avatarURL = &url
		}
	}

	resp := &generated.UserResponse{
		Id:        ptr.ToPtr(id),
		Email:     openapitypes.Email(user.Email),
		Name:      user.Name,
		Bio:       &user.Bio,
		AvatarUrl: avatarURL,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
	return resp
}

func mapUserToUserSummary(user *domain.User) *generated.UserSummary {
	id, _ := uuid.Parse(user.ID)

	return &generated.UserSummary{
		Id:    id,
		Email: openapitypes.Email(user.Email),
		Name:  user.Name,
	}
}
