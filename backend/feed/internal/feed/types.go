package feed

// FeedTrip ist ein Trip-Eintrag im Feed mit Score-Informationen.
type FeedTrip struct {
	TripID    string  `json:"tripId"`
	Title     string  `json:"title"`
	CreatedAt string  `json:"createdAt"`
	CreatorID string  `json:"creatorId"`
	Likes     int64   `json:"likes"`
	Comments  int64   `json:"comments"`
	Score     float64 `json:"score"`
}

// FeedResponse ist die Antwort des Feed-Endpoints.
type FeedResponse struct {
	Data   []FeedTrip `json:"data"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
