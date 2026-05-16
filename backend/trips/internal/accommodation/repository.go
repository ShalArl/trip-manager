package accommodation

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

type Accommodation struct {
	ID         string
	TripID     string
	Name       string
	Address    *string
	CheckIn    *time.Time
	CheckOut   *time.Time
	BookingRef *string
	Notes      *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ── Record ────────────────────────────────────────────────────────────────────

type accommodationRecord struct {
	ID         uuid.UUID  `db:"id"`
	TripID     uuid.UUID  `db:"trip_id"`
	Name       string     `db:"name"`
	Address    *string    `db:"address"`
	CheckIn    *time.Time `db:"check_in"`
	CheckOut   *time.Time `db:"check_out"`
	BookingRef *string    `db:"booking_ref"`
	Notes      *string    `db:"notes"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func (r *accommodationRecord) toDomain() *Accommodation {
	return &Accommodation{
		ID:         r.ID.String(),
		TripID:     r.TripID.String(),
		Name:       r.Name,
		Address:    r.Address,
		CheckIn:    r.CheckIn,
		CheckOut:   r.CheckOut,
		BookingRef: r.BookingRef,
		Notes:      r.Notes,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	ListByTrip(ctx context.Context, tripID string) ([]*Accommodation, error)
	GetByID(ctx context.Context, id string) (*Accommodation, error)
	Create(ctx context.Context, a *Accommodation) (*Accommodation, error)
	Update(ctx context.Context, a *Accommodation) (*Accommodation, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) ListByTrip(ctx context.Context, tripID string) ([]*Accommodation, error) {
	var records []accommodationRecord
	query := `SELECT * FROM accommodations WHERE trip_id = $1 ORDER BY check_in ASC NULLS LAST, created_at ASC`
	if err := r.db.SelectContext(ctx, &records, query, tripID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	result := make([]*Accommodation, len(records))
	for i, rec := range records {
		result[i] = rec.toDomain()
	}
	return result, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Accommodation, error) {
	var rec accommodationRecord
	if err := r.db.GetContext(ctx, &rec, `SELECT * FROM accommodations WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toDomain(), nil
}

func (r *repositoryImpl) Create(ctx context.Context, a *Accommodation) (*Accommodation, error) {
	tripID, err := uuid.Parse(a.TripID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid trip_id", ErrInvalidInput)
	}
	rec := &accommodationRecord{
		TripID:     tripID,
		Name:       a.Name,
		Address:    a.Address,
		CheckIn:    a.CheckIn,
		CheckOut:   a.CheckOut,
		BookingRef: a.BookingRef,
		Notes:      a.Notes,
	}
	query := `
		INSERT INTO accommodations (trip_id, name, address, check_in, check_out, booking_ref, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.TripID, rec.Name, rec.Address,
		rec.CheckIn, rec.CheckOut, rec.BookingRef, rec.Notes,
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

func (r *repositoryImpl) Update(ctx context.Context, a *Accommodation) (*Accommodation, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	rec := &accommodationRecord{
		ID:         id,
		Name:       a.Name,
		Address:    a.Address,
		CheckIn:    a.CheckIn,
		CheckOut:   a.CheckOut,
		BookingRef: a.BookingRef,
		Notes:      a.Notes,
	}
	query := `
		UPDATE accommodations
		SET name = $1, address = $2, check_in = $3, check_out = $4,
		    booking_ref = $5, notes = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.Name, rec.Address, rec.CheckIn, rec.CheckOut,
		rec.BookingRef, rec.Notes, rec.ID,
	).Scan(&rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	a.UpdatedAt = rec.UpdatedAt
	return a, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM accommodations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
