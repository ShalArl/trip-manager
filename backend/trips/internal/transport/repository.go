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

type UserSummary struct {
	ID    string
	Name  string
	Email string
}

type Transport struct {
	ID             string
	TripID         string
	CreatedBy      UserSummary
	FromLocationID string
	ToLocationID   string
	DepartureTime  *time.Time
	ArrivalTime    *time.Time
	Type           string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ── Record ────────────────────────────────────────────────────────────────────

type transportRecord struct {
	ID             uuid.UUID  `db:"id"`
	TripID         uuid.UUID  `db:"trip_id"`
	UserID         uuid.UUID  `db:"user_id"`
	UserName       string     `db:"user_name"`
	UserEmail      string     `db:"user_email"`
	FromLocationID uuid.UUID  `db:"from_location_id"`
	ToLocationID   uuid.UUID  `db:"to_location_id"`
	DepartureTime  *time.Time `db:"departure_time"`
	ArrivalTime    *time.Time `db:"arrival_time"`
	Type           string     `db:"type"`
	Notes          *string    `db:"notes"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (r *transportRecord) toDomain() *Transport {
	notes := ""
	if r.Notes != nil {
		notes = *r.Notes
	}
	return &Transport{
		ID:     r.ID.String(),
		TripID: r.TripID.String(),
		CreatedBy: UserSummary{
			ID:    r.UserID.String(),
			Name:  r.UserName,
			Email: r.UserEmail,
		},
		FromLocationID: r.FromLocationID.String(),
		ToLocationID:   r.ToLocationID.String(),
		DepartureTime:  r.DepartureTime,
		ArrivalTime:    r.ArrivalTime,
		Type:           r.Type,
		Notes:          notes,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
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
	fromID, err := uuid.Parse(t.FromLocationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid from_location_id", ErrInvalidInput)
	}
	toID, err := uuid.Parse(t.ToLocationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid to_location_id", ErrInvalidInput)
	}

	var notesPtr *string
	if t.Notes != "" {
		notesPtr = &t.Notes
	}

	var id uuid.UUID
	var createdAt, updatedAt time.Time

	query := `
		INSERT INTO transports (trip_id, user_id, user_name, user_email, from_location_id, to_location_id, departure_time, arrival_time, type, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query,
		tripID, userID, t.CreatedBy.Name, t.CreatedBy.Email,
		fromID, toID, t.DepartureTime, t.ArrivalTime, t.Type, notesPtr,
	).Scan(&id, &createdAt, &updatedAt)
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

func (r *repositoryImpl) Update(ctx context.Context, t *Transport) (*Transport, error) {
	id, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	userID, err := uuid.Parse(t.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user_id", ErrInvalidInput)
	}
	fromID, err := uuid.Parse(t.FromLocationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid from_location_id", ErrInvalidInput)
	}
	toID, err := uuid.Parse(t.ToLocationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid to_location_id", ErrInvalidInput)
	}

	var notesPtr *string
	if t.Notes != "" {
		notesPtr = &t.Notes
	}

	query := `
		UPDATE transports
		SET from_location_id = $1, to_location_id = $2, departure_time = $3,
		    arrival_time = $4, type = $5, notes = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at`
	err = r.db.QueryRowContext(ctx, query,
		fromID, toID, t.DepartureTime, t.ArrivalTime, t.Type, notesPtr,
		id, userID,
	).Scan(&t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return t, nil
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
