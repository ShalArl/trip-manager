package service

import (
	"log"
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/google/uuid"
)

func mapCreateUserRequestToUser(request *generated.CreateUserRequest) *domain.User {
	return &domain.User{
		ID:           uuid.New().String(),
		Email:        string(request.Email),
		PasswordHash: request.Password,
		Name:         request.Name,
		Bio:          "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

	if request.Bio != nil {
		updated.Bio = *request.Bio
	}

	if request.AvatarKey != nil {
		updated.AvatarKey = *request.AvatarKey
		log.Printf("[Mapper] mapUpdateUserRequestToUser: Set AvatarKey from request: %s", *request.AvatarKey)
	} else {
		log.Printf("[Mapper] mapUpdateUserRequestToUser: No AvatarKey in request")
	}

	updated.UpdatedAt = time.Now()

	return &updated
}
