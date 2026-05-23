package domain

import "time"

// TripLike represents a user liking a trip
type TripLike struct {
	ID        string    `json:"id"`
	TripID    string    `json:"tripId"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

// TripComment represents a comment on a trip
type TripComment struct {
	ID        string    `json:"id"`
	TripID    string    `json:"tripId"`
	UserID    string    `json:"userId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
