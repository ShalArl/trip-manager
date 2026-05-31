package newsletter

import (
	"context"
	"time"
)

type NewsletterTrip struct {
	TripID          string    `json:"tripId"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	Destination     string    `json:"destination,omitempty"`
	CoverImageURL   string    `json:"coverImageUrl,omitempty"`
	CreatorID       string    `json:"creatorId"`
	CreatorName     string    `json:"creatorName"`
	LikeCount       int64     `json:"likeCount"`
	CommentCount    int64     `json:"commentCount"`
	RelevanceReason string    `json:"relevanceReason"`
	CreatedAt       time.Time `json:"createdAt"`
}

type NewsletterSection struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Trips       []NewsletterTrip `json:"trips"`
}

type NewsletterResponse struct {
	Sections    []NewsletterSection `json:"sections"`
	GeneratedAt time.Time           `json:"generatedAt"`
}

type Service interface {
	GetNewsletter(ctx context.Context, userID string, limit int) (*NewsletterResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNewsletter(ctx context.Context, userID string, limit int) (*NewsletterResponse, error) {
	perSection := limit / 3
	if perSection < 3 {
		perSection = 3
	}

	creatorTrips, err := s.repo.GetCreatorTrips(ctx, userID, perSection)
	if err != nil {
		return nil, err
	}

	socialTrips, err := s.repo.GetSocialGraphTrips(ctx, userID, perSection)
	if err != nil {
		return nil, err
	}

	collaborativeTrips, err := s.repo.GetCollaborativeTrips(ctx, userID, perSection)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var sections []NewsletterSection

	if section := buildSection(
		"From Travellers You Follow",
		"New trips from creators you've interacted with",
		creatorTrips,
		seen,
	); len(section.Trips) > 0 {
		sections = append(sections, section)
	}

	if section := buildSection(
		"Popular in Your Network",
		"Trips liked by travellers with similar tastes",
		socialTrips,
		seen,
	); len(section.Trips) > 0 {
		sections = append(sections, section)
	}

	if section := buildSection(
		"Trending Among Your Peers",
		"Highly liked trips from your travel community",
		collaborativeTrips,
		seen,
	); len(section.Trips) > 0 {
		sections = append(sections, section)
	}

	return &NewsletterResponse{
		Sections:    sections,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

func buildSection(title, description string, nodes []TripNode, seen map[string]bool) NewsletterSection {
	section := NewsletterSection{
		Title:       title,
		Description: description,
		Trips:       []NewsletterTrip{},
	}
	for _, n := range nodes {
		if seen[n.TripID] {
			continue
		}
		seen[n.TripID] = true
		section.Trips = append(section.Trips, NewsletterTrip{
			TripID:          n.TripID,
			Title:           n.Title,
			CreatorID:       n.CreatorID,
			CreatorName:     n.CreatorName,
			LikeCount:       n.LikeCount,
			CommentCount:    n.CommentCount,
			RelevanceReason: n.RelevanceReason,
			CreatedAt:       n.CreatedAt,
		})
	}
	return section
}
