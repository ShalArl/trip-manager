package repository

import (
	"time"

	"github.com/google/uuid"
)

// UserRecord represents a user in the database
type userRecord struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Bio          *string   `db:"bio"`
	AvatarKey    *string   `db:"avatar_key"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
