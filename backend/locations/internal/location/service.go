package location

import (
	"context"
	"fmt"
	"time"
)

// ── Input Types ───────────────────────────────────────────────────────────────

type CreateInput struct {
	TripID           string
	UserID           string
	UserName         string
	UserEmail        string
	UserAvatarKey    *string
	Name             string
	City             string
	Country          string
	CountryCode      string
	ShortDescription string
	DateFrom         time.Time
	DateTo           time.Time
	Latitude         *float64
	Longitude        *float64
	Notes            *string
	Sequence         *int
}

type UpdateInput struct {
	ID               string
	Name             *string
	City             *string
	Country          *string
	CountryCode      *string
	ShortDescription *string
	DateFrom         *time.Time
	DateTo           *time.Time
	Latitude         *float64
	Longitude        *float64
	Notes            *string
	Sequence         *int
}

type AddImageInput struct {
	LocationID string
	ImageKey   string
	Sequence   int
}

// ── Service ───────────────────────────────────────────────────────────────────

type Service interface {
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Location, int, error)
	GetByID(ctx context.Context, id string) (*Location, error)
	Create(ctx context.Context, input CreateInput) (*Location, error)
	Update(ctx context.Context, input UpdateInput) (*Location, error)
	Delete(ctx context.Context, id string) error
	AddImage(ctx context.Context, input AddImageInput) (*LocationImage, error)
	DeleteImage(ctx context.Context, imageID string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Location, int, error) {
	return s.repo.ListByTrip(ctx, tripID, limit, offset)
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*Location, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, input CreateInput) (*Location, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if input.City == "" {
		return nil, fmt.Errorf("%w: city is required", ErrInvalidInput)
	}
	if input.Country == "" {
		return nil, fmt.Errorf("%w: country is required", ErrInvalidInput)
	}
	if input.ShortDescription == "" {
		return nil, fmt.Errorf("%w: short_description is required", ErrInvalidInput)
	}
	if input.DateTo.Before(input.DateFrom) {
		return nil, fmt.Errorf("%w: date_to must be after date_from", ErrInvalidInput)
	}
	return s.repo.Create(ctx, &Location{
		TripID: input.TripID,
		CreatedBy: UserSummary{
			ID:        input.UserID,
			Name:      input.UserName,
			Email:     input.UserEmail,
			AvatarKey: input.UserAvatarKey,
		},
		Name:             input.Name,
		City:             input.City,
		Country:          input.Country,
		ShortDescription: input.ShortDescription,
		DateFrom:         input.DateFrom,
		DateTo:           input.DateTo,
		Latitude:         input.Latitude,
		Longitude:        input.Longitude,
		Notes:            input.Notes,
		Sequence:         input.Sequence,
	})
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*Location, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.City != nil {
		existing.City = *input.City
	}
	if input.Country != nil {
		existing.Country = *input.Country
	}
	if input.ShortDescription != nil {
		existing.ShortDescription = *input.ShortDescription
	}
	if input.DateFrom != nil {
		existing.DateFrom = *input.DateFrom
	}
	if input.DateTo != nil {
		existing.DateTo = *input.DateTo
	}
	if input.Latitude != nil {
		existing.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		existing.Longitude = input.Longitude
	}
	if input.Notes != nil {
		existing.Notes = input.Notes
	}
	if input.Sequence != nil {
		existing.Sequence = input.Sequence
	}
	return s.repo.Update(ctx, existing)
}

func (s *serviceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *serviceImpl) AddImage(ctx context.Context, input AddImageInput) (*LocationImage, error) {
	if input.ImageKey == "" {
		return nil, fmt.Errorf("%w: image_key is required", ErrInvalidInput)
	}
	if _, err := s.repo.GetByID(ctx, input.LocationID); err != nil {
		return nil, err
	}
	return s.repo.AddImage(ctx, &LocationImage{
		LocationID: input.LocationID,
		ImageKey:   input.ImageKey,
		Sequence:   input.Sequence,
	})
}

func (s *serviceImpl) DeleteImage(ctx context.Context, imageID string) error {
	return s.repo.DeleteImage(ctx, imageID)
}
