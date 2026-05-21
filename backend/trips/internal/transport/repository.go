package transport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

type Place struct {
	Name    string
	City    string
	Country string
	Lat     *float64
	Lng     *float64
}

type Transport struct {
	ID            string
	TripID        string
	CreatedBy     UserSummary
	From          Place
	To            Place
	DepartureTime *time.Time
	ArrivalTime   *time.Time
	Type          string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ── Record ────────────────────────────────────────────────────────────────────

type transportRecord struct {
	ID            uuid.UUID  `db:"id"`
	TripID        uuid.UUID  `db:"trip_id"`
	UserID        uuid.UUID  `db:"user_id"`
	UserName      string     `db:"user_name"`
	UserEmail     string     `db:"user_email"`
	FromName      *string    `db:"from_name"`
	FromCity      *string    `db:"from_city"`
	FromCountry   *string    `db:"from_country"`
	FromLat       *float64   `db:"from_lat"`
	FromLng       *float64   `db:"from_lng"`
	ToName        *string    `db:"to_name"`
	ToCity        *string    `db:"to_city"`
	ToCountry     *string    `db:"to_country"`
	ToLat         *float64   `db:"to_lat"`
	ToLng         *float64   `db:"to_lng"`
	DepartureTime *time.Time `db:"departure_time"`
	ArrivalTime   *time.Time `db:"arrival_time"`
	Type          string     `db:"type"`
	Notes         *string    `db:"notes"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (rec *transportRecord) toDomain() *Transport {
	return &Transport{
		ID:     rec.ID.String(),
		TripID: rec.TripID.String(),
		CreatedBy: UserSummary{
			ID:    rec.UserID.String(),
			Name:  rec.UserName,
			Email: rec.UserEmail,
		},
		From: Place{
			Name:    derefStr(rec.FromName),
			City:    derefStr(rec.FromCity),
			Country: derefStr(rec.FromCountry),
			Lat:     rec.FromLat,
			Lng:     rec.FromLng,
		},
		To: Place{
			Name:    derefStr(rec.ToName),
			City:    derefStr(rec.ToCity),
			Country: derefStr(rec.ToCountry),
			Lat:     rec.ToLat,
			Lng:     rec.ToLng,
		},
		DepartureTime: rec.DepartureTime,
		ArrivalTime:   rec.ArrivalTime,
		Type:          rec.Type,
		Notes:         derefStr(rec.Notes),
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     rec.UpdatedAt,
	}
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	GetByID(ctx context.Context, id string) (*Transport, error)
	Create(ctx context.Context, t *Transport) (*Transport, error)
	Update(ctx context.Context, t *Transport) (*Transport, error)
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Transport, int, error)
	Delete(ctx context.Context, id, userID string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
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
	userID, err := uuid.Parse(t.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user_id", ErrInvalidInput)
	}

	var notesPtr *string
	if t.Notes != "" {
		notesPtr = &t.Notes
	}

	var id uuid.UUID
	query := `
		INSERT INTO transports (
			trip_id, user_id, user_name, user_email,
			from_name, from_city, from_country, from_lat, from_lng,
			to_name, to_city, to_country, to_lat, to_lng,
			departure_time, arrival_time, type, notes
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18
		) RETURNING id`
	err = r.db.QueryRowContext(ctx, query,
		tripID, userID, t.CreatedBy.Name, t.CreatedBy.Email,
		t.From.Name, t.From.City, t.From.Country, t.From.Lat, t.From.Lng,
		t.To.Name, t.To.City, t.To.Country, t.To.Lat, t.To.Lng,
		t.DepartureTime, t.ArrivalTime, t.Type, notesPtr,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return r.GetByID(ctx, id.String())
}

func (r *repositoryImpl) Update(ctx context.Context, t *Transport) (*Transport, error) {
	id, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	userID, err := uuid.Parse(t.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user_id", ErrInvalidInput)
	}

	var notesPtr *string
	if t.Notes != "" {
		notesPtr = &t.Notes
	}

	query := `
		UPDATE transports
		SET from_name = $1, from_city = $2, from_country = $3, from_lat = $4, from_lng = $5,
		    to_name = $6, to_city = $7, to_country = $8, to_lat = $9, to_lng = $10,
		    departure_time = $11, arrival_time = $12, type = $13, notes = $14,
		    updated_at = NOW()
		WHERE id = $15 AND user_id = $16
		RETURNING updated_at`
	err = r.db.QueryRowContext(ctx, query,
		t.From.Name, t.From.City, t.From.Country, t.From.Lat, t.From.Lng,
		t.To.Name, t.To.City, t.To.Country, t.To.Lat, t.To.Lng,
		t.DepartureTime, t.ArrivalTime, t.Type, notesPtr,
		id, userID,
	).Scan(&t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return r.GetByID(ctx, t.ID)
}

func (r *repositoryImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Transport, int, error) {
	var results []struct {
		transportRecord
		TotalCount int `db:"total_count"`
	}
	query := `
		SELECT *, COUNT(*) OVER() as total_count
		FROM transports
		WHERE trip_id = $1
		ORDER BY departure_time ASC NULLS LAST
		LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &results, query, tripID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Transport{}, 0, nil
	}
	transports := make([]*Transport, len(results))
	for i, res := range results {
		transports[i] = res.transportRecord.toDomain()
	}
	return transports, results[0].TotalCount, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM transports WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
