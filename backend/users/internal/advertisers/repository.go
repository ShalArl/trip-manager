package advertiser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type Advertiser struct {
	ID          string
	FirebaseUID string
	Email       string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tenants     []string // tenant IDs
}

type advertiserRecord struct {
	ID          string    `db:"id"`
	FirebaseUID string    `db:"firebase_uid"`
	Email       string    `db:"email"`
	Name        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, firebaseUID, email, name string) (*Advertiser, error)
	GetByID(ctx context.Context, id string) (*Advertiser, error)
	GetByFirebaseUID(ctx context.Context, uid string) (*Advertiser, error)
	List(ctx context.Context) ([]*Advertiser, error)
	AssignTenant(ctx context.Context, advertiserID, tenantID string) error
	RemoveTenant(ctx context.Context, advertiserID, tenantID string) error
	GetTenants(ctx context.Context, advertiserID string) ([]string, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func generateID() string {
	return fmt.Sprintf("adv-%d", time.Now().UnixNano())
}

func (r *repositoryImpl) Create(ctx context.Context, firebaseUID, email, name string) (*Advertiser, error) {
	id := generateID()
	var rec advertiserRecord
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO advertisers (id, firebase_uid, email, name)
        VALUES ($1, $2, $3, $4)
        RETURNING id, firebase_uid, email, name, created_at, updated_at`,
		id, firebaseUID, email, name,
	).Scan(&rec.ID, &rec.FirebaseUID, &rec.Email, &rec.Name, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConflict, err)
	}
	return &Advertiser{
		ID: rec.ID, FirebaseUID: rec.FirebaseUID,
		Email: rec.Email, Name: rec.Name,
		CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt,
	}, nil
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Advertiser, error) {
	var rec advertiserRecord
	err := r.db.GetContext(ctx, &rec,
		`SELECT id, firebase_uid, email, name, created_at, updated_at FROM advertisers WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tenants, _ := r.GetTenants(ctx, id)
	return &Advertiser{
		ID: rec.ID, FirebaseUID: rec.FirebaseUID,
		Email: rec.Email, Name: rec.Name,
		CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt,
		Tenants: tenants,
	}, nil
}

func (r *repositoryImpl) GetByFirebaseUID(ctx context.Context, uid string) (*Advertiser, error) {
	var rec advertiserRecord
	err := r.db.GetContext(ctx, &rec,
		`SELECT id, firebase_uid, email, name, created_at, updated_at FROM advertisers WHERE firebase_uid = $1`, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tenants, _ := r.GetTenants(ctx, rec.ID)
	return &Advertiser{
		ID: rec.ID, FirebaseUID: rec.FirebaseUID,
		Email: rec.Email, Name: rec.Name,
		CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt,
		Tenants: tenants,
	}, nil
}

func (r *repositoryImpl) List(ctx context.Context) ([]*Advertiser, error) {
	var recs []advertiserRecord
	err := r.db.SelectContext(ctx, &recs,
		`SELECT id, firebase_uid, email, name, created_at, updated_at FROM advertisers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	var result []*Advertiser
	for _, rec := range recs {
		tenants, _ := r.GetTenants(ctx, rec.ID)
		result = append(result, &Advertiser{
			ID: rec.ID, FirebaseUID: rec.FirebaseUID,
			Email: rec.Email, Name: rec.Name,
			CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt,
			Tenants: tenants,
		})
	}
	return result, nil
}

func (r *repositoryImpl) AssignTenant(ctx context.Context, advertiserID, tenantID string) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO advertiser_tenants (advertiser_id, tenant_id)
        VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		advertiserID, tenantID,
	)
	return err
}

func (r *repositoryImpl) RemoveTenant(ctx context.Context, advertiserID, tenantID string) error {
	_, err := r.db.ExecContext(ctx, `
        DELETE FROM advertiser_tenants WHERE advertiser_id = $1 AND tenant_id = $2`,
		advertiserID, tenantID,
	)
	return err
}

func (r *repositoryImpl) GetTenants(ctx context.Context, advertiserID string) ([]string, error) {
	var tenantIDs []string
	err := r.db.SelectContext(ctx, &tenantIDs, `
        SELECT tenant_id FROM advertiser_tenants WHERE advertiser_id = $1 ORDER BY created_at`,
		advertiserID,
	)
	return tenantIDs, err
}
