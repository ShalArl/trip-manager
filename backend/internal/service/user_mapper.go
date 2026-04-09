package service

import (
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/google/uuid"
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
