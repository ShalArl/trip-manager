package repository

import (
	"time"

	"github.com/google/uuid"
)

type accommodationRecord struct {
	ID            uuid.UUID  `db:"id"`
	TripID        uuid.UUID  `db:"trip_id"`
	UserID        uuid.UUID  `db:"user_id"`
	LocationID    uuid.UUID  `db:"location_id"`
	Name          string     `db:"name"`
	Address       *string    `db:"address"`
	CheckIn       *time.Time `db:"check_in"`
	CheckOut      *time.Time `db:"check_out"`
	PricePerNight *float32   `db:"price_per_night"`
	Notes         *string    `db:"notes"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	UserName      string     `db:"user_name"`
	UserEmail     string     `db:"user_email"`
}
