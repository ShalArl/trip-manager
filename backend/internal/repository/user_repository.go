package repository

import (
	"context"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/jmoiron/sqlx"
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
	var user domain.User
	query := `SELECT id, email, name, bio, password_hash, created_at, updated_at FROM users WHERE id = $1`
	
	if err := u.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail implements [UserRepository].
func (u *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, name, bio, password_hash, created_at, updated_at FROM users WHERE email = $1`
	
	if err := u.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser implements [UserRepository].
func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (id, email, name, bio, password_hash, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7)
	         RETURNING id, email, name, bio, password_hash, created_at, updated_at`
	
	err := u.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Name, user.Bio, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
		Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser implements [UserRepository].
func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `UPDATE users SET email = $1, name = $2, bio = $3, password_hash = $4, updated_at = $5 
	         WHERE id = $6
	         RETURNING id, email, name, bio, password_hash, created_at, updated_at`
	
	err := u.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Bio, user.PasswordHash, user.UpdatedAt, user.ID).
		Scan(&user.ID, &user.Email, &user.Name, &user.Bio, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
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
