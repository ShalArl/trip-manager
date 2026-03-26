package domain

import "time"

type ActivityCategory string

const (
	CatSightseeing   ActivityCategory = "sightseeing"
	CatDining        ActivityCategory = "dining"
	CatTransport     ActivityCategory = "transport"
	CatAccommodation ActivityCategory = "accommodation"
	CatOther         ActivityCategory = "other"
)

type Activity struct {
	ResourceMeta
	LocationID  string
	TripID      string
	Name        string
	Description string
	Date        time.Time
	StartTime   string // Oder time.Time, falls du nur die Uhrzeit brauchst
	EndTime     string
	Category    ActivityCategory
	Cost        float64
	Currency    string
}
