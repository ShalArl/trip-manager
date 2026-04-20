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

type LocationRepository interface {
	GetLocation(context context.Context, id string) (*domain.Location, error)
	CreateLocation(context context.Context, location *domain.Location) (*domain.Location, error)
	UpdateLocation(context context.Context, location *domain.Location) (*domain.Location, error)
	ListLocations(context context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error)
	DeleteLocation(context context.Context, id string, userId string) error
}

type LocationRepositoryImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	db *sqlx.DB
}

// GetLocation implements [LocationRepository].
func (l *LocationRepositoryImpl) GetLocation(ctx context.Context, id string) (*domain.Location, error) {
	var rec locationRecord
	query := `
		SELECT l.*, u.id AS user_id, u.name AS user_name, u.email AS user_email
		FROM locations l
		JOIN users u ON l.user_id = u.id
		WHERE l.id = $1`

	if err := l.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toLocation(), nil
}

// CreateLocation implements [LocationRepository].
func (l *LocationRepositoryImpl) CreateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	// Verify trip belongs to user
	rec, err := locationToRecord(location)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
	INSERT INTO locations (trip_id, user_id, name, city, country, short_description, date_from, date_to, latitude, longitude, notes, sequence)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, created_at, updated_at`

	err = l.db.QueryRowContext(ctx, query,
		rec.TripID, rec.UserID, rec.Name, rec.City, rec.Country,
		rec.ShortDescription, rec.DateFrom, rec.DateTo,
		rec.Latitude, rec.Longitude, rec.Notes, rec.Sequence,
	).Scan(&location.ResourceMeta.ID, &location.ResourceMeta.CreatedAt, &location.ResourceMeta.UpdatedAt)

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

	return l.GetLocation(ctx, location.ResourceMeta.ID)
}

// UpdateLocation implements [LocationRepository].
func (l *LocationRepositoryImpl) UpdateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	rec, err := locationToRecord(location)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
	UPDATE locations 
	SET name = $1, city = $2, country = $3, short_description = $4, date_from = $5, date_to = $6,
		latitude = $7, longitude = $8, notes = $9, sequence = $10, updated_at = NOW() 
	WHERE id = $11 AND user_id = $12
	RETURNING updated_at`

	err = l.db.QueryRowContext(ctx, query,
		rec.Name, rec.City, rec.Country, rec.ShortDescription, rec.DateFrom, rec.DateTo,
		rec.Latitude, rec.Longitude, rec.Notes, rec.Sequence,
		rec.ID, rec.UserID,
	).Scan(&location.ResourceMeta.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toLocation(), nil
}

// ListLocations implements [LocationRepository].
func (l *LocationRepositoryImpl) ListLocations(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Location, int, error) {
	var results []struct {
		locationRecord
		TotalCount int `db:"total_count"`
	}

	query := `
		SELECT 
		    l.*, 
		    u.id AS user_id, 
            u.name AS user_name, 
            u.email AS user_email, 
            COUNT(*) OVER () as total_count 
		FROM locations l JOIN users u ON u.id = l.user_id
		WHERE trip_id = $1
		ORDER BY sequence ASC, created_at ASC
		LIMIT $2 OFFSET $3`

	if err := l.db.SelectContext(ctx, &results, query, tripID, limit, offset); err != nil {
		fmt.Printf("ListLocations DB error: %v\n", err)
		return nil, 0, domain.ErrInternal
	}

	if len(results) == 0 {
		return []*domain.Location{}, 0, nil
	}

	locations := make([]*domain.Location, len(results))
	for i, res := range results {
		locations[i] = res.locationRecord.toLocation()
	}

	return locations, results[0].TotalCount, nil
}

// DeleteLocation implements [LocationRepository].
func (l *LocationRepositoryImpl) DeleteLocation(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM locations WHERE id = $1 AND user_id = $2`
	result, err := l.db.ExecContext(ctx, query, id, userID)
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

func NewLocationRepository(db *sqlx.DB) LocationRepository {
	return &LocationRepositoryImpl{db: db}
}
