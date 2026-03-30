package repository

import (
	"time"

	"github.com/google/uuid"
)

// TripRecord represents a trip in the database
type tripRecord struct {
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	Title            string    `db:"title"`
	ShortDescription string    `db:"short_description"`
	Destination      string    `db:"destination"`
	Description      *string   `db:"description"`
	StartDate        time.Time `db:"start_date"`
	EndDate          time.Time `db:"end_date"`
	Status           string    `db:"status"` // planned | ongoing | completed | cancelled
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	UserName         string    `db:"user_name"`
	UserEmail        string    `db:"user_email"`
}
