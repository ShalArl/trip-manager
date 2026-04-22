package domain

import "time"

type User struct {
	ID           string
	Email        string
	Name         string
	Bio          string
	PasswordHash string
	AvatarKey    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
