package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ── Record ────────────────────────────────────────────────────────────────────

type userRecord struct {
	ID          uuid.UUID `db:"id"`
	Email       string    `db:"email"`
	Name        string    `db:"name"`
	Bio         *string   `db:"bio"`
	AvatarKey   *string   `db:"avatar_key"`
	FirebaseUID string    `db:"firebase_uid"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ── Domain Type ───────────────────────────────────────────────────────────────

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

// ── Errors ────────────────────────────────────────────────────────────────────

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")
)

// ── Mapper ────────────────────────────────────────────────────────────────────

func (r *userRecord) toUser() *User {
	return &User{
		ID:          r.ID.String(),
		Email:       r.Email,
		Name:        r.Name,
		Bio:         fromPtr(r.Bio),
		AvatarKey:   fromPtr(r.AvatarKey),
		FirebaseUID: r.FirebaseUID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toRecord(u *User) (*userRecord, error) {
	var id uuid.UUID
	var err error
	if u.ID != "" {
		id, err = uuid.Parse(u.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
	}
	return &userRecord{
		ID:          id,
		Email:       u.Email,
		Name:        u.Name,
		Bio:         toPtr(u.Bio),
		AvatarKey:   toPtr(u.AvatarKey),
		FirebaseUID: u.FirebaseUID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}, nil
}

func fromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByFirebaseUID(ctx context.Context, uid string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*User, error) {
	var rec userRecord
	query := `SELECT id, email, name, bio, avatar_key, firebase_uid, created_at, updated_at 
	          FROM users WHERE id = $1`
	if err := r.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toUser(), nil
}

func (r *repositoryImpl) GetByFirebaseUID(ctx context.Context, uid string) (*User, error) {
	var rec userRecord
	log.Printf("[Repository] GetByFirebaseUID: %s", uid)
	query := `SELECT id, email, name, bio, avatar_key, firebase_uid, created_at, updated_at 
	          FROM users WHERE firebase_uid = $1`
	if err := r.db.GetContext(ctx, &rec, query, uid); err != nil {
		log.Printf("[Repository] GetByFirebaseUID error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	log.Printf("[Repository] GetByFirebaseUID found: %s", rec.ID)
	return rec.toUser(), nil
}

func (r *repositoryImpl) Create(ctx context.Context, user *User) (*User, error) {
	rec, err := toRecord(user)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users (email, name, bio, firebase_uid, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, NOW(), NOW())
	          RETURNING id, email, name, bio, avatar_key, firebase_uid, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query, rec.Email, rec.Name, rec.Bio, rec.FirebaseUID).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.AvatarKey, &rec.FirebaseUID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toUser(), nil
}

func (r *repositoryImpl) Update(ctx context.Context, user *User) (*User, error) {
	rec, err := toRecord(user)
	if err != nil {
		return nil, err
	}
	query := `UPDATE users SET email = $1, name = $2, bio = $3, avatar_key = $4, updated_at = NOW()
	          WHERE id = $5
	          RETURNING id, email, name, bio, avatar_key, firebase_uid, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query, rec.Email, rec.Name, rec.Bio, rec.AvatarKey, rec.ID).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.AvatarKey, &rec.FirebaseUID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toUser(), nil
}
