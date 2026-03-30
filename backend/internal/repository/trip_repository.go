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

type TripRepository interface {
	// GetTrip retrieves a trip by ID
	GetTrip(ctx context.Context, id string) (*domain.Trip, error)

	// CreateTrip creates a new trip and returns the created record
	CreateTrip(ctx context.Context, trip *domain.Trip) (*domain.Trip, error)

	// UpdateTrip updates an existing trip
	UpdateTrip(ctx context.Context, trip *domain.Trip) (*domain.Trip, error)

	// ListTrips retrieves trips for a user with pagination
	ListTrips(ctx context.Context, userID string, limit int, offset int) ([]*domain.Trip, int, error)

	// DeleteTrip removes a trip
	DeleteTrip(ctx context.Context, id string, userId string) error
}

type TripRepositoryImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	db *sqlx.DB
}

func (t *TripRepositoryImpl) GetTrip(ctx context.Context, id string) (*domain.Trip, error) {
	query := `
		SELECT t.*, u.id AS user_id, u.email AS user_email, u.name AS user_name 
		FROM trips t JOIN users u ON t.user_id = u.id 
		WHERE t.id = $1`

	var rec tripRecord
	if err := t.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toTrip(), nil
}

func (t *TripRepositoryImpl) CreateTrip(ctx context.Context, trip *domain.Trip) (*domain.Trip, error) {
	rec, err := tripToRecord(trip)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
		INSERT INTO trips (user_id, title, short_description, description, start_date, end_date, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err = t.db.QueryRowContext(
		ctx, query, rec.UserID, rec.Title, rec.ShortDescription, rec.Description, rec.StartDate, rec.EndDate, rec.Status,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign_key_violation
				return nil, fmt.Errorf("%w: referenced trip not found", domain.ErrInvalidInput)
			case "23505": // unique_violation
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toTrip(), nil
}

func (t *TripRepositoryImpl) UpdateTrip(ctx context.Context, trip *domain.Trip) (*domain.Trip, error) {
	rec, err := tripToRecord(trip)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
		UPDATE trips 
		SET title = $1, short_description = $2, description = $3, start_date = $4, end_date = $5, status = $6 
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at`

	err = t.db.QueryRowContext(ctx, query,
		rec.Title, rec.ShortDescription, rec.Description, rec.StartDate, rec.EndDate, rec.Status, rec.ID, rec.UserID,
	).Scan(&rec.UpdatedAt)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: referenced trip not found", domain.ErrInvalidInput)
			case "23505":
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toTrip(), nil
}

func (t *TripRepositoryImpl) ListTrips(ctx context.Context, userID string, limit int, offset int) ([]*domain.Trip, int, error) {
	var results []struct {
		tripRecord
		TotalCount int `db:"total_count"`
	}

	query := `
		SELECT 
			t.*, 
			u.id AS user_id, 
			u.email AS user_email, 
			u.name AS user_name, 
			COUNT(*) OVER() as total_count
		FROM trips t JOIN users u ON t.user_id = u.id 
		WHERE t.user_id = $1 
		ORDER BY t.start_date ASC, t.created_at ASC
		LIMIT $2 OFFSET $3`

	if err := t.db.SelectContext(ctx, &results, query, userID, limit, offset); err != nil {
		fmt.Printf("ListTrips DB error: %v\n", err)
		return nil, 0, domain.ErrInternal
	}

	if len(results) == 0 {
		return []*domain.Trip{}, 0, nil
	}

	trips := make([]*domain.Trip, len(results))
	for i, res := range results {
		trips[i] = res.tripRecord.toTrip()
	}

	return trips, results[0].TotalCount, nil
}

func (t *TripRepositoryImpl) DeleteTrip(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM trips WHERE id = $1 AND user_id = $2`
	result, err := t.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return domain.ErrInternal
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.ErrInternal
	}

	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func NewTripRepository(db *sqlx.DB) TripRepository {
	return &TripRepositoryImpl{db: db}
}
