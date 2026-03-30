package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

// GetUser implements [UserRepository].
func (u *UserRepositoryImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var rec userRecord
	query := `SELECT id, email, name, bio, password_hash, created_at, updated_at FROM users WHERE id = $1`

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
	query := `SELECT id, email, name, bio, password_hash, created_at, updated_at FROM users WHERE email = $1`
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

	query := `INSERT INTO users (id, email, name, bio, password_hash, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7)
	         RETURNING id, email, name, bio, password_hash, created_at, updated_at`

	err = u.db.QueryRowContext(ctx, query, rec.ID, rec.Email, rec.Name, rec.Bio, rec.PasswordHash, rec.CreatedAt, rec.UpdatedAt).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.PasswordHash, &rec.CreatedAt, &user.UpdatedAt)

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

	return user, nil
}

// UpdateUser implements [UserRepository].
func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	rec, err := userToRecord(user)
	if err != nil {
		return nil, err
	}

	query := `UPDATE users SET email = $1, name = $2, bio = $3, password_hash = $4, updated_at = $5 
	         WHERE id = $6
	         RETURNING id, email, name, bio, password_hash, created_at, updated_at`

	err = u.db.QueryRowContext(ctx, query, rec.Email, rec.Name, rec.Bio, rec.PasswordHash, rec.UpdatedAt, rec.ID).
		Scan(&rec.ID, &rec.Email, &rec.Name, &rec.Bio, &rec.PasswordHash, &rec.CreatedAt, &rec.UpdatedAt)

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

	return user, nil
}

// DeleteUser implements [UserRepository].
func (u *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}
