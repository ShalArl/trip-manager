package domain

import "time"

type Accommodation struct {
	ResourceMeta
	TripID        string
	LocationID    string
	Name          string
	Address       string
	CheckIn       *time.Time
	CheckOut      *time.Time
	PricePerNight *float32
	Notes         string
}
