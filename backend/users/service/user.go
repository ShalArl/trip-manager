package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/shared/firebaseclient"
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
	TenantID    string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProvisionInput struct {
	FirebaseUID string
	Email       string
	Name        string
	TenantID    string // optional, default "default"
	Role        string // optional, default "tenant_member"
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
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*User, error)
	Provision(ctx context.Context, input ProvisionInput) (*User, bool, error)
	Update(ctx context.Context, input UpdateInput) (*User, error)
}

// ── Implementation ────────────────────────────────────────────────────────────

type serviceImpl struct {
	repo         repository.Repository
	firebaseAuth *firebaseclient.Client
}

func NewService(repo repository.Repository, firebaseAuth *firebaseclient.Client) Service {
	return &serviceImpl{repo: repo, firebaseAuth: firebaseAuth}
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoToService(u), nil
}

func (s *serviceImpl) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*User, error) {
	u, err := s.repo.GetByFirebaseUID(ctx, firebaseUID)
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

	// Defaults
	tenantID := input.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	role := input.Role
	if role == "" {
		role = "tenant_member"
	}

	// Create new user
	created, err := s.repo.Create(ctx, &repository.User{
		Email:       input.Email,
		Name:        input.Name,
		FirebaseUID: input.FirebaseUID,
		TenantID:    tenantID,
		Role:        role,
	})
	if err != nil {
		return nil, false, err
	}

	// Firebase Custom Claims setzen
	if s.firebaseAuth != nil {
		claims := map[string]interface{}{
			"tenant_id": tenantID,
			"role":      role,
		}
		if err := s.firebaseAuth.SetCustomClaims(ctx, input.FirebaseUID, claims); err != nil {
			// Nicht fatal – User ist angelegt, Claims können nachträglich gesetzt werden
			fmt.Printf("warn: failed to set firebase custom claims for %s: %v\n", input.FirebaseUID, err)
		}
	}

	return repoToService(created), true, nil
}

func (s *serviceImpl) Update(ctx context.Context, input UpdateInput) (*User, error) {
	existing, err := s.repo.GetByFirebaseUID(ctx, input.ID)
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
		TenantID:    u.TenantID,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
