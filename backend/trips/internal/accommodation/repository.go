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

type UserSummary struct {
	ID    string
	Name  string
	Email string
}

type Accommodation struct {
	ID            string
	TripID        string
	CreatedBy     UserSummary
	LocationID    string
	Name          string
	Address       string
	CheckIn       *time.Time
	CheckOut      *time.Time
	PricePerNight *float32
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ── Record ────────────────────────────────────────────────────────────────────

type accommodationRecord struct {
	ID            uuid.UUID  `db:"id"`
	TripID        uuid.UUID  `db:"trip_id"`
	UserID        uuid.UUID  `db:"user_id"`
	UserName      string     `db:"user_name"`
	UserEmail     string     `db:"user_email"`
	LocationID    uuid.UUID  `db:"location_id"`
	Name          string     `db:"name"`
	Address       *string    `db:"address"`
	CheckIn       *time.Time `db:"check_in"`
	CheckOut      *time.Time `db:"check_out"`
	PricePerNight *float32   `db:"price_per_night"`
	Notes         *string    `db:"notes"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

func (r *accommodationRecord) toDomain() *Accommodation {
	address := ""
	if r.Address != nil {
		address = *r.Address
	}
	notes := ""
	if r.Notes != nil {
		notes = *r.Notes
	}
	return &Accommodation{
		ID:     r.ID.String(),
		TripID: r.TripID.String(),
		CreatedBy: UserSummary{
			ID:    r.UserID.String(),
			Name:  r.UserName,
			Email: r.UserEmail,
		},
		LocationID:    r.LocationID.String(),
		Name:          r.Name,
		Address:       address,
		CheckIn:       r.CheckIn,
		CheckOut:      r.CheckOut,
		PricePerNight: r.PricePerNight,
		Notes:         notes,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func toRecord(a *Accommodation) (*accommodationRecord, error) {
	var id uuid.UUID
	var err error
	if a.ID != "" {
		id, err = uuid.Parse(a.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid accommodation ID: %w", err)
		}
	}
	tripID, err := uuid.Parse(a.TripID)
	if err != nil {
		return nil, fmt.Errorf("invalid trip ID: %w", err)
	}
	locationID, err := uuid.Parse(a.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}
	userID, err := uuid.Parse(a.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var addressPtr *string
	if a.Address != "" {
		addressPtr = &a.Address
	}
	var notesPtr *string
	if a.Notes != "" {
		notesPtr = &a.Notes
	}

	return &accommodationRecord{
		ID:            id,
		TripID:        tripID,
		UserID:        userID,
		UserName:      a.CreatedBy.Name,
		UserEmail:     a.CreatedBy.Email,
		LocationID:    locationID,
		Name:          a.Name,
		Address:       addressPtr,
		CheckIn:       a.CheckIn,
		CheckOut:      a.CheckOut,
		PricePerNight: a.PricePerNight,
		Notes:         notesPtr,
	}, nil
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	GetByID(ctx context.Context, id string) (*Accommodation, error)
	Create(ctx context.Context, a *Accommodation) (*Accommodation, error)
	Update(ctx context.Context, a *Accommodation) (*Accommodation, error)
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Accommodation, int, error)
	Delete(ctx context.Context, id, userID string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
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
	rec, err := toRecord(a)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	var id uuid.UUID
	query := `
		INSERT INTO accommodations (trip_id, user_id, user_name, user_email, location_id, name, address, check_in, check_out, price_per_night, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`
	err = r.db.QueryRowContext(ctx, query,
		rec.TripID, rec.UserID, rec.UserName, rec.UserEmail,
		rec.LocationID, rec.Name, rec.Address,
		rec.CheckIn, rec.CheckOut, rec.PricePerNight, rec.Notes,
	).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: referenced trip or location not found", ErrInvalidInput)
			case "23505":
				return nil, fmt.Errorf("%w: conflict", ErrInternal)
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return r.GetByID(ctx, id.String())
}

func (r *repositoryImpl) Update(ctx context.Context, a *Accommodation) (*Accommodation, error) {
	rec, err := toRecord(a)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
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
	).Scan(&a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return r.GetByID(ctx, a.ID)
}

func (r *repositoryImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Accommodation, int, error) {
	var results []struct {
		accommodationRecord
		TotalCount int `db:"total_count"`
	}
	query := `
		SELECT *, COUNT(*) OVER() as total_count
		FROM accommodations
		WHERE trip_id = $1
		ORDER BY check_in ASC NULLS LAST, created_at ASC
		LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &results, query, tripID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Accommodation{}, 0, nil
	}
	accommodations := make([]*Accommodation, len(results))
	for i, res := range results {
		accommodations[i] = res.accommodationRecord.toDomain()
	}
	return accommodations, results[0].TotalCount, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM accommodations WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
