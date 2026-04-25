package repository

import (
	"time"

	"github.com/google/uuid"
)

// ActivityRecord represents an activity in the database
type transportRecord struct {
	ID             uuid.UUID  `db:"id"`
	TripID         uuid.UUID  `db:"trip_id"`
	FromLocationID uuid.UUID  `db:"from_location_id"`
	ToLocationID   uuid.UUID  `db:"to_location_id"`
	ArrivalTime    *time.Time `db:"arrival_time"`
	DepartureTime  *time.Time `db:"departure_time"`
	Type           string     `db:"type"`
	Notes          *string    `db:"notes"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	UserID         uuid.UUID  `db:"user_id"`
	UserName       string     `db:"user_name"`
	UserEmail      string     `db:"user_email"`
}
