package repository

import (
	"time"

	"github.com/google/uuid"
)

type locationRecord struct {
	ID               uuid.UUID `db:"id"`
	TripID           uuid.UUID `db:"trip_id"`
	Name             string    `db:"name"`
	City             string    `db:"city"`
	Country          string    `db:"country"`
	ShortDescription string    `db:"short_description"`
	DateFrom         time.Time `db:"date_from"`
	DateTo           time.Time `db:"date_to"`
	Latitude         *float64  `db:"latitude"`
	Longitude        *float64  `db:"longitude"`
	Notes            *string   `db:"notes"`
	Sequence         *int      `db:"sequence"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	UserID           uuid.UUID `db:"user_id"`
	UserName         string    `db:"user_name"`
	UserEmail        string    `db:"user_email"`
}

type locationImageRecord struct {
	ID         uuid.UUID `db:"id"`
	LocationID uuid.UUID `db:"location_id"`
	ImageKey   string    `db:"image_key"`
	Sequence   *int      `db:"sequence"`
	CreatedAt  time.Time `db:"created_at"`
}
