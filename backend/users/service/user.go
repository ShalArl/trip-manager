package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/users/repository"
)

// ── Types ─────────────────────────────────────────────────────────────────────

type User struct {
	ID          string
	Email       string
	Name        string
	Bio         string
	AvatarKey   string
	FirebaseUID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProvisionInput struct {
	FirebaseUID string
	Email       string
	Name        string
}

type UpdateInput struct {
	ID        string
	Name      *string
	Email     *string
	Bio       *string
	AvatarKey *string
}

// ── Interface ─────────────────────────────────────────────────────────────────

type Service interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Provision(ctx context.Context, input ProvisionInput) (*User, bool, error)
	Update(ctx context.Context, input UpdateInput) (*User, error)
}

// ── Implementation ────────────────────────────────────────────────────────────

type serviceImpl struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoToService(u), nil
}

func (s *serviceImpl) Provision(ctx context.Context, input ProvisionInput) (*User, bool, error) {
	// Check if user already exists
	existing, err := s.repo.GetByFirebaseUID(ctx, input.FirebaseUID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, false, fmt.Errorf("checking existing user: %w", err)
	}
	if existing != nil {
		return repoToService(existing), false, nil
	}

	// Create new user
	created, err := s.repo.Create(ctx, &repository.User{
		Email:       input.Email,
		Name:        input.Name,
		FirebaseUID: input.FirebaseUID,
	})
	if err != nil {
		return nil, false, err
	}
	return repoToService(created), true, nil
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*User, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.Bio != nil {
		existing.Bio = *input.Bio
	}
	if input.AvatarKey != nil {
		existing.AvatarKey = *input.AvatarKey
	}
	existing.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return repoToService(updated), nil
}

// ── Mapper ────────────────────────────────────────────────────────────────────

func repoToService(u *repository.User) *User {
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Bio:         u.Bio,
		AvatarKey:   u.AvatarKey,
		FirebaseUID: u.FirebaseUID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
