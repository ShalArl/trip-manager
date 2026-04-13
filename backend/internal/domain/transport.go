package domain

import "time"

type TransportType string

const (
	TransportFlight TransportType = "flight"
	TransportTrain  TransportType = "train"
	TransportCar    TransportType = "car"
	TransportBus    TransportType = "bus"
)

type Transport struct {
	ResourceMeta
	TripID         string
	FromLocationID string
	ToLocationID   string
	Date           time.Time
	Type           TransportType
	Notes          string
}
