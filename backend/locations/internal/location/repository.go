package location

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
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)

// ── Domain Types ──────────────────────────────────────────────────────────────

type UserSummary struct {
	ID        string
	Name      string
	Email     string
	AvatarKey *string
}

type LocationImage struct {
	ID         string
	LocationID string
	ImageKey   string
	Sequence   int
	CreatedAt  time.Time
}

type Location struct {
	ID               string
	TripID           string
	CreatedBy        UserSummary
	Name             string
	City             string
	Country          string
	ShortDescription string
	DateFrom         time.Time
	DateTo           time.Time
	Latitude         *float64
	Longitude        *float64
	Notes            *string
	Sequence         *int
	Images           []LocationImage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ── Records ───────────────────────────────────────────────────────────────────

type locationRecord struct {
	ID               uuid.UUID `db:"id"`
	TripID           uuid.UUID `db:"trip_id"`
	UserID           uuid.UUID `db:"user_id"`
	UserName         string    `db:"user_name"`
	UserEmail        string    `db:"user_email"`
	UserAvatarKey    *string   `db:"user_avatar_key"`
	Name             string    `db:"name"`
	City             string    `db:"city"`
	Country          string    `db:"country"`
	ShortDescription string    `db:"short_description"`
	DateFrom         time.Time `db:"date_from"`
	DateTo           time.Time `db:"date_to"`
	Latitude         *float64  `db:"latitude"`
	Longitude        *float64  `db:"longitude"`
	Notes            *string   `db:"notes"`
	Sequence         *int      `db:"sequence"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type locationImageRecord struct {
	ID         uuid.UUID `db:"id"`
	LocationID uuid.UUID `db:"location_id"`
	ImageKey   string    `db:"image_key"`
	Sequence   int       `db:"sequence"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *locationRecord) toDomain(images []LocationImage) *Location {
	return &Location{
		ID:     r.ID.String(),
		TripID: r.TripID.String(),
		CreatedBy: UserSummary{
			ID:        r.UserID.String(),
			Name:      r.UserName,
			Email:     r.UserEmail,
			AvatarKey: r.UserAvatarKey,
		},
		Name:             r.Name,
		City:             r.City,
		Country:          r.Country,
		ShortDescription: r.ShortDescription,
		DateFrom:         r.DateFrom,
		DateTo:           r.DateTo,
		Latitude:         r.Latitude,
		Longitude:        r.Longitude,
		Notes:            r.Notes,
		Sequence:         r.Sequence,
		Images:           images,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func (r *locationImageRecord) toDomain() LocationImage {
	return LocationImage{
		ID:         r.ID.String(),
		LocationID: r.LocationID.String(),
		ImageKey:   r.ImageKey,
		Sequence:   r.Sequence,
		CreatedAt:  r.CreatedAt,
	}
}

// ── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Location, int, error)
	GetByID(ctx context.Context, id string) (*Location, error)
	Create(ctx context.Context, l *Location) (*Location, error)
	Update(ctx context.Context, l *Location) (*Location, error)
	Delete(ctx context.Context, id string) error
	AddImage(ctx context.Context, img *LocationImage) (*LocationImage, error)
	DeleteImage(ctx context.Context, imageID string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) listImages(ctx context.Context, locationID string) ([]LocationImage, error) {
	var records []locationImageRecord
	if err := r.db.SelectContext(ctx, &records,
		`SELECT * FROM location_images WHERE location_id = $1 ORDER BY sequence ASC, created_at ASC`,
		locationID,
	); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	images := make([]LocationImage, len(records))
	for i, rec := range records {
		images[i] = rec.toDomain()
	}
	return images, nil
}

func (r *repositoryImpl) ListByTrip(ctx context.Context, tripID string, limit, offset int) ([]*Location, int, error) {
	var results []struct {
		locationRecord
		TotalCount int `db:"total_count"`
	}
	query := `
		SELECT *, COUNT(*) OVER() as total_count
		FROM locations
		WHERE trip_id = $1
		ORDER BY sequence ASC NULLS LAST, date_from ASC, created_at ASC
		LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &results, query, tripID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if len(results) == 0 {
		return []*Location{}, 0, nil
	}
	locations := make([]*Location, len(results))
	for i, res := range results {
		images, err := r.listImages(ctx, res.locationRecord.ID.String())
		if err != nil {
			return nil, 0, err
		}
		locations[i] = res.locationRecord.toDomain(images)
	}
	return locations, results[0].TotalCount, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Location, error) {
	var rec locationRecord
	if err := r.db.GetContext(ctx, &rec, `SELECT * FROM locations WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	images, err := r.listImages(ctx, id)
	if err != nil {
		return nil, err
	}
	return rec.toDomain(images), nil
}

func (r *repositoryImpl) Create(ctx context.Context, l *Location) (*Location, error) {
	tripID, err := uuid.Parse(l.TripID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid trip_id", ErrInvalidInput)
	}
	userID, err := uuid.Parse(l.CreatedBy.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user_id", ErrInvalidInput)
	}
	rec := &locationRecord{
		TripID:           tripID,
		UserID:           userID,
		UserName:         l.CreatedBy.Name,
		UserEmail:        l.CreatedBy.Email,
		UserAvatarKey:    l.CreatedBy.AvatarKey,
		Name:             l.Name,
		City:             l.City,
		Country:          l.Country,
		ShortDescription: l.ShortDescription,
		DateFrom:         l.DateFrom,
		DateTo:           l.DateTo,
		Latitude:         l.Latitude,
		Longitude:        l.Longitude,
		Notes:            l.Notes,
		Sequence:         l.Sequence,
	}
	query := `
		INSERT INTO locations (trip_id, user_id, user_name, user_email, user_avatar_key, name, city, country, short_description, date_from, date_to, latitude, longitude, notes, sequence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, query,
		rec.TripID, rec.UserID, rec.UserName, rec.UserEmail, rec.UserAvatarKey,
		rec.Name, rec.City, rec.Country, rec.ShortDescription,
		rec.DateFrom, rec.DateTo, rec.Latitude, rec.Longitude, rec.Notes, rec.Sequence,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: %v", ErrInternal, pgErr)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return rec.toDomain([]LocationImage{}), nil
}

func (r *repositoryImpl) Update(ctx context.Context, l *Location) (*Location, error) {
	id, err := uuid.Parse(l.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	var updatedAt time.Time
	err = r.db.QueryRowContext(ctx, `
		UPDATE locations
		SET name = $1, city = $2, country = $3, short_description = $4,
		    date_from = $5, date_to = $6, latitude = $7, longitude = $8,
		    notes = $9, sequence = $10, updated_at = NOW()
		WHERE id = $11
		RETURNING updated_at`,
		l.Name, l.City, l.Country, l.ShortDescription,
		l.DateFrom, l.DateTo, l.Latitude, l.Longitude,
		l.Notes, l.Sequence, id,
	).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	l.UpdatedAt = updatedAt
	images, err := r.listImages(ctx, l.ID)
	if err != nil {
		return nil, err
	}
	l.Images = images
	return l, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM locations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repositoryImpl) AddImage(ctx context.Context, img *LocationImage) (*LocationImage, error) {
	locationID, err := uuid.Parse(img.LocationID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid location_id", ErrInvalidInput)
	}
	rec := &locationImageRecord{
		LocationID: locationID,
		ImageKey:   img.ImageKey,
		Sequence:   img.Sequence,
	}
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO location_images (location_id, image_key, sequence)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		rec.LocationID, rec.ImageKey, rec.Sequence,
	).Scan(&rec.ID, &rec.CreatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: %v", ErrInternal, pgErr)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	result := rec.toDomain()
	return &result, nil
}

func (r *repositoryImpl) DeleteImage(ctx context.Context, imageID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM location_images WHERE id = $1`, imageID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
