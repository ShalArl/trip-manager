package domain

import "time"

type Point struct {
	Lat float64
	Lon float64
}

type Location struct {
	ResourceMeta
	TripID           string
	Name             string
	City             string
	Country          string
	Coordinates      Point
	Notes            string
	Sequence         int
	ShortDescription string
	DateFrom         time.Time
	DateTo           time.Time
}
