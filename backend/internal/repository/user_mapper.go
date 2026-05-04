package repository

import (
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
)

func (u *userRecord) toUser() *domain.User {
	return &domain.User{
		ID:          u.ID.String(),
		Email:       u.Email,
		Name:        u.Name,
		Bio:         ptr.FromPtr(u.Bio),
		AvatarKey:   ptr.FromPtr(u.AvatarKey),
		FirebaseUID: u.FirebaseUID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
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
		ID:          id,
		Email:       user.Email,
		Name:        user.Name,
		Bio:         ptr.ToPtr(user.Bio),
		AvatarKey:   ptr.ToPtr(user.AvatarKey),
		FirebaseUID: user.FirebaseUID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
