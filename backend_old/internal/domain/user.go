package domain

import (
	"time"
)

type User struct {
	ID          string
	FirebaseUID string
	Email       string
	Name        string
	Bio         string
	AvatarKey   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
