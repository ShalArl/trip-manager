package service

import (
	"log"
	"time"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
)

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
