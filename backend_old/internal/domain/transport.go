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
	DepartureTime  *time.Time
	ArrivalTime    *time.Time
	Type           TransportType
	Notes          string
}
