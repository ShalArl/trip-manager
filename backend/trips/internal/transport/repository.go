package transport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ── Errors ────────────────────────────────────────────────────────────────────

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")
)

// ── Domain Types ──────────────────────────────────────────────────────────────

type Transport struct {
	ID             string
	TripID         string
	Type           string
	DeparturePlace string
	ArrivalPlace   string
	DepartureTime  *time.Time
	ArrivalTime    *time.Time
	BookingRef     *string
	Notes          *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ── Record ────────────────────────────────────────────────────────────────────

type transportRecord struct {
	ID             uuid.UUID  `db:"id"`
	TripID         uuid.UUID  `db:"trip_id"`
	Type           string     `db:"type"`
	DeparturePlace string     `db:"departure_place"`
	ArrivalPlace   string     `db:"arrival_place"`
	DepartureTime  *time.Time `db:"departure_time"`
	ArrivalTime    *time.Time `db:"arrival_time"`
	BookingRef     *string    `db:"booking_ref"`
	Notes          *string    `db:"notes"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (r *transportRecord) toDomain() *Transport {
	return &Transport{
		ID:             r.ID.String(),
		TripID:         r.TripID.String(),
		Type:           r.Type,
		DeparturePlace: r.DeparturePlace,
		ArrivalPlace:   r.ArrivalPlace,
		DepartureTime:  r.DepartureTime,
		ArrivalTime:    r.ArrivalTime,
		BookingRef:     r.BookingRef,
		Notes:          r.Notes,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	ListByTrip(ctx context.Context, tripID string) ([]*Transport, error)
	GetByID(ctx context.Context, id string) (*Transport, error)
	Create(ctx context.Context, t *Transport) (*Transport, error)
	Update(ctx context.Context, t *Transport) (*Transport, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) ListByTrip(ctx context.Context, tripID string) ([]*Transport, error) {
	var records []transportRecord
	query := `SELECT * FROM transports WHERE trip_id = $1 ORDER BY departure_time ASC NULLS LAST, created_at ASC`
	if err := r.db.SelectContext(ctx, &records, query, tripID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	result := make([]*Transport, len(records))
	for i, rec := range records {
		result[i] = rec.toDomain()
	}
	return result, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Transport, error) {
	var rec transportRecord
	if err := r.db.GetContext(ctx, &rec, `SELECT * FROM transports WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toDomain(), nil
}

func (r *repositoryImpl) Create(ctx context.Context, t *Transport) (*Transport, error) {
	tripID, err := uuid.Parse(t.TripID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid trip_id", ErrInvalidInput)
	}
	rec := &transportRecord{
		TripID:         tripID,
		Type:           t.Type,
		DeparturePlace: t.DeparturePlace,
		ArrivalPlace:   t.ArrivalPlace,
		DepartureTime:  t.DepartureTime,
		ArrivalTime:    t.ArrivalTime,
		BookingRef:     t.BookingRef,
		Notes:          t.Notes,
	}
	query := `
		INSERT INTO transports (trip_id, type, departure_place, arrival_place, departure_time, arrival_time, booking_ref, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.TripID, rec.Type, rec.DeparturePlace, rec.ArrivalPlace,
		rec.DepartureTime, rec.ArrivalTime, rec.BookingRef, rec.Notes,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: %v", ErrInternal, pgErr)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toDomain(), nil
}

func (r *repositoryImpl) Update(ctx context.Context, t *Transport) (*Transport, error) {
	id, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	rec := &transportRecord{
		ID:             id,
		Type:           t.Type,
		DeparturePlace: t.DeparturePlace,
		ArrivalPlace:   t.ArrivalPlace,
		DepartureTime:  t.DepartureTime,
		ArrivalTime:    t.ArrivalTime,
		BookingRef:     t.BookingRef,
		Notes:          t.Notes,
	}
	query := `
		UPDATE transports
		SET type = $1, departure_place = $2, arrival_place = $3,
		    departure_time = $4, arrival_time = $5, booking_ref = $6,
		    notes = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.Type, rec.DeparturePlace, rec.ArrivalPlace,
		rec.DepartureTime, rec.ArrivalTime, rec.BookingRef,
		rec.Notes, rec.ID,
	).Scan(&rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	t.UpdatedAt = rec.UpdatedAt
	return t, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM transports WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
