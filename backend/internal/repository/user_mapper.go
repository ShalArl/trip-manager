package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (u *userRecord) toUser() *domain.User {
	return &domain.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		Name:         u.Name,
		Bio:          ptr.FromPtr(u.Bio),
		AvatarURL:    ptr.FromPtr(u.AvatarURL),
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func userToRecord(user *domain.User) (*userRecord, error) {
	var id uuid.UUID
	var err error

	if user.ID != "" {
		id, err = uuid.Parse(user.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID for user ID: %v", err)
		}
	}

	return &userRecord{
		ID:           id,
		Email:        user.Email,
		Name:         user.Name,
		Bio:          ptr.ToPtr(user.Bio),
		AvatarURL:    ptr.ToPtr(user.AvatarURL),
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
