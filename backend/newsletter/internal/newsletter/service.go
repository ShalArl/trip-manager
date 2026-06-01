package newsletter

import (
	"context"
	"time"
)

type NewsletterTrip struct {
	TripID          string    `json:"tripId"`
	Title           string    `json:"title"`
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
	GetNewsletter(ctx context.Context, firebaseUID string) (*NewsletterResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNewsletter(ctx context.Context, firebaseUID string) (*NewsletterResponse, error) {
	sections, generatedAt, err := s.repo.GetStoredNewsletter(ctx, firebaseUID)
	if err != nil {
		return nil, err
	}

	var newsletterSections []NewsletterSection
	for _, sec := range sections {
		var trips []NewsletterTrip
		for _, t := range sec.Trips {
			trips = append(trips, NewsletterTrip{
				TripID:          t.TripID,
				Title:           t.Title,
				CreatorID:       t.CreatorID,
				CreatorName:     t.CreatorName,
				LikeCount:       t.LikeCount,
				CommentCount:    t.CommentCount,
				RelevanceReason: t.RelevanceReason,
				CreatedAt:       t.CreatedAt,
			})
		}
		newsletterSections = append(newsletterSections, NewsletterSection{
			Title:       sec.Title,
			Description: sec.Description,
			Trips:       trips,
		})
	}

	if newsletterSections == nil {
		newsletterSections = []NewsletterSection{}
	}

	return &NewsletterResponse{
		Sections:    newsletterSections,
		GeneratedAt: generatedAt,
	}, nil
}
