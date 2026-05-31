package newsletter

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type TripNode struct {
	TripID          string    `json:"tripId"`
	Title           string    `json:"title"`
	CreatorID       string    `json:"creatorId"`
	CreatorName     string    `json:"creatorName"`
	LikeCount       int64     `json:"likeCount"`
	CommentCount    int64     `json:"commentCount"`
	RelevanceReason string    `json:"relevanceReason"`
	CreatedAt       time.Time `json:"createdAt"`
}

type StoredSection struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Trips       []TripNode `json:"trips"`
}

type Repository interface {
	GetStoredNewsletter(ctx context.Context, firebaseUID string) ([]StoredSection, time.Time, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetStoredNewsletter(ctx context.Context, firebaseUID string) ([]StoredSection, time.Time, error) {
	var content []byte
	var generatedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT content, generated_at FROM newsletters WHERE firebase_uid = $1`,
		firebaseUID,
	).Scan(&content, &generatedAt)

	if err == sql.ErrNoRows {
		return nil, time.Time{}, nil
	}
	if err != nil {
		return nil, time.Time{}, err
	}

	var sections []StoredSection
	if err := json.Unmarshal(content, &sections); err != nil {
		return nil, time.Time{}, err
	}

	return sections, generatedAt, nil
}
