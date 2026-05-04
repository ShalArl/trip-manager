package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserProfile(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByFirebaseUID(ctx context.Context, uid string) (*domain.User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

// GetUser implements [UserRepository].
func (u *UserRepositoryImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var rec userRecord
	query := `SELECT id, email, name, bio, avatar_key, created_at, updated_at, firebase_uid FROM users WHERE id = $1`

	if err := u.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toUser(), nil
}

// GetUserByEmail implements [UserRepository].
func (u *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var rec userRecord
	query := `SELECT id, email, name, bio, avatar_key, created_at, updated_at, firebase_uid FROM users WHERE email = $1`
	if err := u.db.GetContext(ctx, &rec, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toUser(), nil
}

// CreateUser implements [UserRepository].
func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	rec, err := userToRecord(user)
	if err != nil {
		return nil, err
	}

	println("[Repository] CreateUser with id: ", user.ID)

	query := `INSERT INTO users (email, name, bio, firebase_uid, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6)
	         RETURNING id, email, name, bio, firebase_uid, created_at, updated_at`

	err = u.db.QueryRowContext(ctx, query, rec.Email, rec.Name, rec.Bio, rec.FirebaseUID, rec.CreatedAt, rec.UpdatedAt).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.FirebaseUID, &rec.CreatedAt, &rec.UpdatedAt)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: referenced user not found", domain.ErrInvalidInput)
			case "23505":
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toUser(), nil
}

// UpdateUserProfile implements [UserRepository].
func (u *UserRepositoryImpl) UpdateUserProfile(ctx context.Context, user *domain.User) (*domain.User, error) {
	rec, err := userToRecord(user)
	if err != nil {
		return nil, err
	}

	log.Printf("[Repository] UpdateUserProfile: User before DB update - ID=%s, AvatarKey=%s", rec.ID, *rec.AvatarKey)

	query := `UPDATE users SET email = $1, name = $2, bio = $3, avatar_key = $4, updated_at = $5
	         WHERE id = $6
	         RETURNING id, email, name, bio, avatar_key, created_at, updated_at, firebase_uid`

	err = u.db.QueryRowContext(ctx, query, rec.Email, rec.Name, rec.Bio, rec.AvatarKey, rec.UpdatedAt, rec.ID).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.AvatarKey, &rec.CreatedAt, &rec.UpdatedAt, &rec.FirebaseUID)

	log.Printf("[Repository] UpdateUserProfile: User after DB update - ID=%s, AvatarKey=%s", rec.ID, *rec.AvatarKey)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: referenced user not found", domain.ErrInvalidInput)
			case "23505":
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toUser(), nil
}

// DeleteUser implements [UserRepository].
func (u *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (u *UserRepositoryImpl) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `SELECT id, firebase_uid, email, name, bio, avatar_key, created_at, updated_at
              FROM users WHERE firebase_uid = $1`

	var rec userRecord
	err := u.db.QueryRowContext(ctx, query, firebaseUID).Scan(
		&rec.ID, &rec.FirebaseUID, &rec.Email, &rec.Name,
		&rec.Bio, &rec.AvatarKey, &rec.CreatedAt, &rec.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toUser(), nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}
