package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type LocationRepository interface {
	GetLocation(ctx context.Context, id string) (*domain.Location, error)
	CreateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error)
	UpdateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error)
	ListLocations(ctx context.Context, tripId string, limit int, offset int) ([]*domain.Location, int, error)
	DeleteLocation(ctx context.Context, id string, userId string) error
	AddLocationImage(ctx context.Context, locationID string, imageKey string, sequence *int) (*domain.LocationImage, error)
	DeleteLocationImage(ctx context.Context, imageID string, locationID string) error
	ListLocationImages(ctx context.Context, locationID string) ([]domain.LocationImage, error)
}

type LocationRepositoryImpl struct {
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

	location := rec.toLocation()

	images, err := l.ListLocationImages(ctx, id)
	if err != nil {
		return nil, err
	}
	location.Images = images

	return location, nil
}

// CreateLocation implements [LocationRepository].
func (l *LocationRepositoryImpl) CreateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error) {
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
			case "23503":
				return nil, fmt.Errorf("%w: referenced trip not found", domain.ErrInvalidInput)
			case "23505":
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

	return l.GetLocation(ctx, location.ResourceMeta.ID)
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
		loc := res.locationRecord.toLocation()
		images, err := l.ListLocationImages(ctx, loc.ID)
		if err != nil {
			return nil, 0, err
		}
		loc.Images = images
		locations[i] = loc
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

// AddLocationImage implements [LocationRepository].
func (l *LocationRepositoryImpl) AddLocationImage(ctx context.Context, locationID string, imageKey string, sequence *int) (*domain.LocationImage, error) {
	locID, err := uuid.Parse(locationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid location ID", domain.ErrInvalidInput)
	}

	var rec locationImageRecord
	query := `
		INSERT INTO location_images (location_id, image_key, sequence)
		VALUES ($1, $2, $3)
		RETURNING id, location_id, image_key, sequence, created_at`

	err = l.db.QueryRowContext(ctx, query, locID, imageKey, sequence).
		Scan(&rec.ID, &rec.LocationID, &rec.ImageKey, &rec.Sequence, &rec.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	img := rec.toLocationImage()
	return &img, nil
}

// DeleteLocationImage implements [LocationRepository].
func (l *LocationRepositoryImpl) DeleteLocationImage(ctx context.Context, imageID string, locationID string) error {
	query := `DELETE FROM location_images WHERE id = $1 AND location_id = $2`
	result, err := l.db.ExecContext(ctx, query, imageID, locationID)
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

// ListLocationImages implements [LocationRepository].
func (l *LocationRepositoryImpl) ListLocationImages(ctx context.Context, locationID string) ([]domain.LocationImage, error) {
	var recs []locationImageRecord
	query := `
		SELECT id, location_id, image_key, sequence, created_at
		FROM location_images
		WHERE location_id = $1
		ORDER BY sequence ASC, created_at ASC`

	if err := l.db.SelectContext(ctx, &recs, query, locationID); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	images := make([]domain.LocationImage, len(recs))
	for i, rec := range recs {
		images[i] = rec.toLocationImage()
	}
	return images, nil
}

func NewLocationRepository(db *sqlx.DB) LocationRepository {
	return &LocationRepositoryImpl{db: db}
}
