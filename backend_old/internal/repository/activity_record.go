package repository

import (
	"time"

	"github.com/google/uuid"
)

// ActivityRecord represents an activity in the database
type activityRecord struct {
	ID          uuid.UUID `db:"id"`
	LocationID  uuid.UUID `db:"location_id"`
	TripID      uuid.UUID `db:"trip_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Date        time.Time `db:"date"`
	StartTime   *string   `db:"start_time"`
	EndTime     *string   `db:"end_time"`
	Category    *string   `db:"category"`
	Cost        *float64  `db:"cost"`
	Currency    string    `db:"currency"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	UserID    uuid.UUID `db:"user_id"`
	UserName  string    `db:"user_name"`
	UserEmail string    `db:"user_email"`
}
