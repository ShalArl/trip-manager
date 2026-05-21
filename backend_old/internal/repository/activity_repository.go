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

type ActivityRepository interface {
	GetActivity(ctx context.Context, id string) (*domain.Activity, error)
	CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error)
	ListActivitiesForLocation(ctx context.Context, locationID string, limit int, offset int) ([]*domain.Activity, int, error)
	ListActivitiesForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Activity, int, error)
	DeleteActivity(ctx context.Context, id string, userId string) error
}

type ActivityRepositoryImpl struct {
	// You can add dependencies here, such as a database connection or logger.
	db *sqlx.DB
}

// GetActivity implements [ActivityRepository].
func (a *ActivityRepositoryImpl) GetActivity(ctx context.Context, id string) (*domain.Activity, error) {
	var rec activityRecord
	query := `
		SELECT a.*, u.id AS user_id, u.name AS user_name, u.email AS user_email
		FROM activities a JOIN users u ON a.user_id = u.id
		WHERE a.id = $1`

	if err := a.db.GetContext(ctx, &rec, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return rec.toActivity(), nil
}

// CreateActivity implements [ActivityRepository].
func (a *ActivityRepositoryImpl) CreateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	rec, err := activityToRecord(activity)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
        INSERT INTO activities 
            (trip_id, location_id, user_id, name, description, date, start_time, end_time, category, cost, currency)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at`

	err = a.db.QueryRowContext(ctx, query,
		rec.TripID, rec.LocationID, rec.UserID, rec.Name, rec.Description, rec.Date,
		rec.StartTime, rec.EndTime, rec.Category, rec.Cost, rec.Currency,
	).Scan(&activity.ResourceMeta.ID, &activity.ResourceMeta.CreatedAt, &activity.ResourceMeta.UpdatedAt)

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

	return activity, nil
}

// UpdateActivity implements [ActivityRepository].
func (a *ActivityRepositoryImpl) UpdateActivity(ctx context.Context, activity *domain.Activity) (*domain.Activity, error) {
	rec, err := activityToRecord(activity)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	query := `
		UPDATE activities 
        SET trip_id = $1, location_id = $2, name = $3, description = $4, date = $5, 
            start_time = $6, end_time = $7, category = $8, cost = $9, currency = $10, 
            updated_at = NOW() 
        WHERE id = $11 AND user_id = $12
        RETURNING updated_at`

	err = a.db.QueryRowContext(ctx, query,
		rec.TripID, rec.LocationID, rec.Name, rec.Description, rec.Date,
		rec.StartTime, rec.EndTime, rec.Category, rec.Cost, rec.Currency,
		rec.ID, rec.UserID,
	).Scan(&activity.ResourceMeta.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return rec.toActivity(), nil
}

// ListActivitiesForLocation implements [ActivityRepository] calls [listActivitiesByField].
func (a *ActivityRepositoryImpl) ListActivitiesForLocation(ctx context.Context, locationID string, limit int, offset int) ([]*domain.Activity, int, error) {

	query := `
		SELECT
			a.*,
			u.id AS user_id,
			u.name AS user_name,
			u.email AS user_email,
			COUNT(*) OVER () as total_count
		FROM activities a
		JOIN users u ON a.user_id = u.id
		WHERE a.location_id = $1
		ORDER BY a.date ASC, a.start_time ASC
		LIMIT $2 OFFSET $3`

	return a.listActivitiesByField(ctx, query, locationID, limit, offset)
}

// ListActivitiesForTrip implements [ActivityRepository] calls [listActivitiesByField].
func (a *ActivityRepositoryImpl) ListActivitiesForTrip(ctx context.Context, tripID string, limit int, offset int) ([]*domain.Activity, int, error) {
	query := `
        SELECT 
            a.*, 
            u.id AS user_id, 
            u.name AS user_name, 
            u.email AS user_email,
            COUNT(*) OVER () as total_count 
        FROM activities a
        	JOIN users u ON a.user_id = u.id
        WHERE a.trip_id = $1 
        ORDER BY a.date ASC, a.start_time ASC
		LIMIT $2 OFFSET $3`

	return a.listActivitiesByField(ctx, query, tripID, limit, offset)
}

// DeleteActivity implements [ActivityRepository].
func (a *ActivityRepositoryImpl) DeleteActivity(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM activities WHERE id = $1 AND user_id = $2`

	result, err := a.db.ExecContext(ctx, query, id, userID)
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

func NewActivityRepository(db *sqlx.DB) ActivityRepository {
	return &ActivityRepositoryImpl{db: db}
}

// Helper Methods

// listActivitiesByField utility method to return a list and total count of [domain.Activity] for the specified [fieldName]
func (a *ActivityRepositoryImpl) listActivitiesByField(ctx context.Context, query string, value string, limit int, offset int) ([]*domain.Activity, int, error) {
	var results []struct {
		activityRecord
		TotalCount int `db:"total_count"`
	}

	if err := a.db.SelectContext(ctx, &results, query, value, limit, offset); err != nil {
		return nil, 0, domain.ErrInternal
	}

	if len(results) == 0 {
		return []*domain.Activity{}, 0, nil
	}

	activities := make([]*domain.Activity, len(results))
	for i, res := range results {
		activities[i] = res.activityRecord.toActivity()
	}

	return activities, results[0].TotalCount, nil
}
