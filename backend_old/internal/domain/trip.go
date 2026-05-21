package domain

import "time"

type TripStatus string

const (
	TripStatusPlanned   TripStatus = "planned"
	TripStatusOngoing   TripStatus = "ongoing"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

type Trip struct {
	ResourceMeta
	Title            string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	Status           TripStatus
}
