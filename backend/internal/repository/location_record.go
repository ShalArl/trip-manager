package repository

import (
	"time"

	"github.com/google/uuid"
)

// LocationRecord represents a location in the database
type locationRecord struct {
	ID        uuid.UUID `db:"id"`
	TripID    uuid.UUID `db:"trip_id"`
	Name      string    `db:"name"`
	City      string    `db:"city"`
	Country   string    `db:"country"`
	Latitude  *float64  `db:"latitude"`
	Longitude *float64  `db:"longitude"`
	Notes     *string   `db:"notes"`
	Sequence  *int      `db:"sequence"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	UserID    uuid.UUID `db:"user_id"`
	UserName  string    `db:"user_name"`
	UserEmail string    `db:"user_email"`
}
