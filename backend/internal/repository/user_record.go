package repository

import (
	"time"

	"github.com/google/uuid"
)

// UserRecord represents a user in the database
type userRecord struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Bio          *string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
