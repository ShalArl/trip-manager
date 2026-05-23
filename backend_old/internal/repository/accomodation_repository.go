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

type AccommodationRepository interface {
	GetAccommodation(ctx context.Context, id string) (*domain.Accommodation, error)
	CreateAccommodation(ctx context.Context, accommodation *domain.Accommodation) (*domain.Accommodation, error)
	UpdateAccommodation(ctx context.Context, accommodation *domain.Accommodation) (*domain.Accommodation, error)
	ListAccommodationsForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Accommodation, int, error)
	DeleteAccommodation(ctx context.Context, id string, userID string) error
}

type AccommodationRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccommodationRepository(db *sqlx.DB) AccommodationRepository {
	return &AccommodationRepositoryImpl{db: db}
}

func (r *AccommodationRepositoryImpl) GetAccommodation(ctx context.Context, id string) (*domain.Accommodation, error) {
	var rec accommodationRecord
	query := `
		SELECT a.*, u.id AS user_id, u.name AS user_name, u.email AS user_email
		FROM accommodations a JOIN users u ON a.user_id = u.id
		WHERE a.id = $1`

	if err := r.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toAccommodation(), nil
}

func (r *AccommodationRepositoryImpl) CreateAccommodation(ctx context.Context, accommodation *domain.Accommodation) (*domain.Accommodation, error) {
	rec, err := accommodationToRecord(accommodation)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
		INSERT INTO accommodations
			(trip_id, user_id, location_id, name, address, check_in, check_out, price_per_night, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRowContext(ctx, query,
		rec.TripID, rec.UserID, rec.LocationID, rec.Name, rec.Address,
		rec.CheckIn, rec.CheckOut, rec.PricePerNight, rec.Notes,
	).Scan(&accommodation.ResourceMeta.ID, &accommodation.ResourceMeta.CreatedAt, &accommodation.ResourceMeta.UpdatedAt)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: referenced trip or location not found", domain.ErrInvalidInput)
			case "23505":
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return r.GetAccommodation(ctx, accommodation.ResourceMeta.ID)
}

func (r *AccommodationRepositoryImpl) UpdateAccommodation(ctx context.Context, accommodation *domain.Accommodation) (*domain.Accommodation, error) {
	rec, err := accommodationToRecord(accommodation)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
		UPDATE accommodations
		SET location_id = $1, name = $2, address = $3, check_in = $4, check_out = $5,
		    price_per_night = $6, notes = $7, updated_at = NOW()
		WHERE id = $8 AND user_id = $9
		RETURNING updated_at`

	err = r.db.QueryRowContext(ctx, query,
		rec.LocationID, rec.Name, rec.Address, rec.CheckIn, rec.CheckOut,
		rec.PricePerNight, rec.Notes, rec.ID, rec.UserID,
	).Scan(&accommodation.ResourceMeta.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return r.GetAccommodation(ctx, accommodation.ResourceMeta.ID)
}

func (r *AccommodationRepositoryImpl) ListAccommodationsForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Accommodation, int, error) {
	query := `
		SELECT
			a.*,
			u.id AS user_id,
			u.name AS user_name,
			u.email AS user_email,
			COUNT(*) OVER () AS total_count
		FROM accommodations a
		JOIN users u ON a.user_id = u.id
		WHERE a.trip_id = $1
		ORDER BY a.check_in ASC NULLS LAST, a.created_at ASC
		LIMIT $2 OFFSET $3`

	var results []struct {
		accommodationRecord
		TotalCount int `db:"total_count"`
	}

	if err := r.db.SelectContext(ctx, &results, query, tripID, limit, offset); err != nil {
		return nil, 0, domain.ErrInternal
	}

	if len(results) == 0 {
		return []*domain.Accommodation{}, 0, nil
	}

	accommodations := make([]*domain.Accommodation, len(results))
	for i, res := range results {
		accommodations[i] = res.accommodationRecord.toAccommodation()
	}

	return accommodations, results[0].TotalCount, nil
}

func (r *AccommodationRepositoryImpl) DeleteAccommodation(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM accommodations WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
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
