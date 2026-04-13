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

type TransportRepository interface {
	GetTransport(ctx context.Context, id string) (*domain.Transport, error)
	CreateTransport(ctx context.Context, transport *domain.Transport) (*domain.Transport, error)
	UpdateTransport(ctx context.Context, transport *domain.Transport) (*domain.Transport, error)
	ListTransportsForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Transport, int, error)
	DeleteTransport(ctx context.Context, id string, userID string) error
}

type TransportRepositoryImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	db *sqlx.DB
}

func (t *TransportRepositoryImpl) GetTransport(ctx context.Context, id string) (*domain.Transport, error) {
	var rec transportRecord
	query := `
        SELECT t.*, u.id AS user_id, u.name AS user_name, u.email AS user_email
        FROM transports t JOIN users u ON t.user_id = u.id
        WHERE t.id = $1`

	if err := t.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toTransport(), nil
}

func (t *TransportRepositoryImpl) CreateTransport(ctx context.Context, transport *domain.Transport) (*domain.Transport, error) {
	rec, err := transportToRecord(transport)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
        INSERT INTO transports
            (trip_id, user_id, from_location_id, to_location_id, date, type, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	err = t.db.QueryRowContext(ctx, query,
		rec.TripID, rec.UserID, rec.FromLocationID, rec.ToLocationID, rec.Date, rec.Type, rec.Notes,
	).Scan(&transport.ResourceMeta.ID, &transport.ResourceMeta.CreatedAt, &transport.ResourceMeta.UpdatedAt)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign_key_violation
				return nil, fmt.Errorf("%w: referenced trip or location not found", domain.ErrInvalidInput)
			case "23505": // unique_violation
				return nil, domain.ErrConflict
			}
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return transport, nil
}

func (t *TransportRepositoryImpl) UpdateTransport(ctx context.Context, transport *domain.Transport) (*domain.Transport, error) {
	rec, err := transportToRecord(transport)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
        UPDATE transports
        SET from_location_id = $1, to_location_id = $2, date = $3, type = $4, notes = $5,
            updated_at = NOW()
        WHERE id = $6 AND user_id = $7
        RETURNING updated_at`

	err = t.db.QueryRowContext(ctx, query,
		rec.FromLocationID, rec.ToLocationID, rec.Date, rec.Type, rec.Notes,
		rec.ID, rec.UserID,
	).Scan(&transport.ResourceMeta.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toTransport(), nil
}

func (t *TransportRepositoryImpl) ListTransportsForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Transport, int, error) {
	query := `
        SELECT
            t.*,
            u.id AS user_id,
            u.name AS user_name,
            u.email AS user_email,
            COUNT(*) OVER () as total_count
        FROM transports t
        JOIN users u ON t.user_id = u.id
        WHERE t.trip_id = $1
        ORDER BY t.date ASC
        LIMIT $2 OFFSET $3`

	return t.listTransportsByField(ctx, query, tripID, limit, offset)
}

func (t *TransportRepositoryImpl) listTransportsByField(ctx context.Context, query string, value string, limit int, offset int) ([]*domain.Transport, int, error) {
	var results []struct {
		transportRecord
		TotalCount int `db:"total_count"`
	}

	if err := t.db.SelectContext(ctx, &results, query, value, limit, offset); err != nil {
		return nil, 0, domain.ErrInternal
	}

	if len(results) == 0 {
		return []*domain.Transport{}, 0, nil
	}

	transports := make([]*domain.Transport, len(results))
	for i, res := range results {
		transports[i] = res.transportRecord.toTransport()
	}

	return transports, results[0].TotalCount, nil
}

func (t *TransportRepositoryImpl) DeleteTransport(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM transports WHERE id = $1 AND user_id = $2`

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

func NewTransportRepository(db *sqlx.DB) TransportRepository {
	return &TransportRepositoryImpl{db: db}
}
