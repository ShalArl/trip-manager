package trip

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

// ── Record ────────────────────────────────────────────────────────────────────

type tripRecord struct {
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	Title            string    `db:"title"`
	ShortDescription string    `db:"short_description"`
	Description      *string   `db:"description"`
	StartDate        time.Time `db:"start_date"`
	EndDate          time.Time `db:"end_date"`
	Status           string    `db:"status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// ── Domain Types ──────────────────────────────────────────────────────────────

type UserSummary struct {
	ID        string
	Name      string
	Email     string
	AvatarKey *string
}

type Trip struct {
	ID               string
	Title            string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	Status           string
	CreatedBy        UserSummary
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ── Errors ────────────────────────────────────────────────────────────────────

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
)

// ── Mapper ────────────────────────────────────────────────────────────────────

func (r *tripRecord) toTrip() *Trip {
	return &Trip{
		ID:               r.ID.String(),
		Title:            r.Title,
		ShortDescription: r.ShortDescription,
		Description:      fromPtr(r.Description),
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		Status:           r.Status,
		CreatedBy: UserSummary{
			ID: r.UserID.String(),
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func toRecord(t *Trip) (*tripRecord, error) {
	var id uuid.UUID
	var err error
	if t.ID != "" {
		id, err = uuid.Parse(t.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid trip ID: %w", err)
		}
	}
	userID, err := uuid.Parse(t.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return &tripRecord{
		ID:               id,
		UserID:           userID,
		Title:            t.Title,
		ShortDescription: t.ShortDescription,
		Description:      toPtr(t.Description),
		StartDate:        t.StartDate,
		EndDate:          t.EndDate,
		Status:           t.Status,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}, nil
}

func fromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	GetByID(ctx context.Context, id string) (*Trip, error)
	Create(ctx context.Context, trip *Trip) (*Trip, error)
	Update(ctx context.Context, trip *Trip) (*Trip, error)
	Delete(ctx context.Context, id, userID string) error
	List(ctx context.Context, userID string, limit, offset int) ([]*Trip, int, error)
	ListRecent(ctx context.Context, limit, offset int) ([]*Trip, int, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Trip, int, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Trip, error) {
	var rec tripRecord
	query := `SELECT * FROM trips WHERE id = $1`
	if err := r.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toTrip(), nil
}

func (r *repositoryImpl) Create(ctx context.Context, trip *Trip) (*Trip, error) {
	rec, err := toRecord(trip)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	query := `
		INSERT INTO trips (user_id, title, short_description, description, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.UserID,
		rec.Title, rec.ShortDescription, rec.Description,
		rec.StartDate, rec.EndDate, rec.Status,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toTrip(), nil
}

func (r *repositoryImpl) Update(ctx context.Context, trip *Trip) (*Trip, error) {
	rec, err := toRecord(trip)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	query := `
		UPDATE trips
		SET title = $1, short_description = $2, description = $3, start_date = $4, end_date = $5, status = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.Title, rec.ShortDescription, rec.Description,
		rec.StartDate, rec.EndDate, rec.Status,
		rec.ID, rec.UserID,
	).Scan(&rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toTrip(), nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM trips WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repositoryImpl) List(ctx context.Context, userID string, limit, offset int) ([]*Trip, int, error) {
	var results []struct {
		tripRecord
		TotalCount int `db:"total_count"`
	}
	query := `
		SELECT *, COUNT(*) OVER() as total_count
		FROM trips
		WHERE user_id = $1
		ORDER BY start_date ASC, created_at ASC
		LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &results, query, userID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Trip{}, 0, nil
	}
	trips := make([]*Trip, len(results))
	for i, res := range results {
		trips[i] = res.tripRecord.toTrip()
	}
	return trips, results[0].TotalCount, nil
}

func (r *repositoryImpl) ListRecent(ctx context.Context, limit, offset int) ([]*Trip, int, error) {
	var results []tripRecord
	query := `
		SELECT * FROM trips
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	if err := r.db.SelectContext(ctx, &results, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM trips`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Trip{}, total, nil
	}
	trips := make([]*Trip, len(results))
	for i, res := range results {
		trips[i] = res.toTrip()
	}
	return trips, total, nil
}

func (r *repositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*Trip, int, error) {
	var results []struct {
		tripRecord
		TotalCount int `db:"total_count"`
	}
	sqlQuery := `
		SELECT *, COUNT(*) OVER() as total_count
		FROM trips
		WHERE title ILIKE $1 OR short_description ILIKE $1
		ORDER BY start_date ASC, created_at ASC
		LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &results, sqlQuery, "%"+query+"%", limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Trip{}, 0, nil
	}
	trips := make([]*Trip, len(results))
	for i, res := range results {
		trips[i] = res.tripRecord.toTrip()
	}
	return trips, results[0].TotalCount, nil
}
