package tenant

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ShalArl/trip-manager/backend/shared/tenantdb"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("tenant not found")
	ErrConflict = errors.New("tenant already exists")
	ErrInternal = errors.New("internal error")
)

type Tenant struct {
	ID        string
	Name      string
	Slug      string
	Tier      string
	Status    string
	Branding  map[string]interface{}
	Settings  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type tenantRecord struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	Tier      string    `db:"tier"`
	Status    string    `db:"status"`
	Branding  []byte    `db:"branding"`
	Settings  []byte    `db:"settings"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, t *Tenant) (*Tenant, error)
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, t *Tenant) (*Tenant, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func GenerateSlug(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	slug := re.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(slug, "-")
}

func (r *repositoryImpl) Create(ctx context.Context, t *Tenant) (*Tenant, error) {
	if t.Slug == "" {
		t.Slug = GenerateSlug(t.Name)
	}

	defaultSettings := map[string]interface{}{"maxActiveTrips": 0}
	if t.Tier == "free" {
		defaultSettings["maxActiveTrips"] = 3
	}
	settingsJSON, err := json.Marshal(defaultSettings)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal default settings: %v", ErrInternal, err)
	}

	var result *Tenant
	err = tenantdb.WithTenant(ctx, r.db, func(tx *sqlx.Tx) error {
		var rec tenantRecord
		err := tx.QueryRowContext(ctx, `
            INSERT INTO tenants (id, name, slug, tier, status, branding, settings)
            VALUES ($1, $2, $3, $4, 'active', '{}', $5)
            RETURNING id, name, slug, tier, status, branding, settings, created_at, updated_at`,
			t.ID, t.Name, t.Slug, t.Tier, settingsJSON,
		).Scan(&rec.ID, &rec.Name, &rec.Slug, &rec.Tier, &rec.Status, &rec.Branding, &rec.Settings, &rec.CreatedAt, &rec.UpdatedAt)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		result = rec.toDomain()
		return nil
	})
	return result, err
}

func (r *repositoryImpl) GetByID(ctx context.Context, id string) (*Tenant, error) {
	var result *Tenant
	err := tenantdb.WithTenant(ctx, r.db, func(tx *sqlx.Tx) error {
		var rec tenantRecord
		if err := tx.GetContext(ctx, &rec, `SELECT * FROM tenants WHERE id = $1`, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		result = rec.toDomain()
		return nil
	})
	return result, err
}

func (r *repositoryImpl) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var result *Tenant
	err := tenantdb.WithTenant(ctx, r.db, func(tx *sqlx.Tx) error {
		var rec tenantRecord
		if err := tx.GetContext(ctx, &rec, `SELECT * FROM tenants WHERE slug = $1`, slug); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		result = rec.toDomain()
		return nil
	})
	return result, err
}

func (r *repositoryImpl) Update(ctx context.Context, t *Tenant) (*Tenant, error) {
	var result *Tenant

	brandingJSON, err := json.Marshal(t.Branding)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal branding: %v", ErrInternal, err)
	}

	settingsJSON, err := json.Marshal(t.Settings)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal settings: %v", ErrInternal, err)
	}

	err = tenantdb.WithTenant(ctx, r.db, func(tx *sqlx.Tx) error {
		var rec tenantRecord
		err := tx.QueryRowContext(ctx, `
            UPDATE tenants SET name = $1, branding = $2, settings = $3, tier = $4, updated_at = NOW()
            WHERE id = $5
            RETURNING id, name, slug, tier, status, branding, settings, created_at, updated_at`,
			t.Name, brandingJSON, settingsJSON, t.Tier, t.ID,
		).Scan(&rec.ID, &rec.Name, &rec.Slug, &rec.Tier, &rec.Status, &rec.Branding, &rec.Settings, &rec.CreatedAt, &rec.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		result = rec.toDomain()
		return nil
	})
	return result, err
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	return tenantdb.WithTenant(ctx, r.db, func(tx *sqlx.Tx) error {
		result, err := tx.ExecContext(ctx, `DELETE FROM tenants WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInternal, err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func (rec *tenantRecord) toDomain() *Tenant {
	var branding map[string]interface{}
	var settings map[string]interface{}

	if len(rec.Branding) > 0 {
		_ = json.Unmarshal(rec.Branding, &branding)
	}
	if len(rec.Settings) > 0 {
		_ = json.Unmarshal(rec.Settings, &settings)
	}

	if branding == nil {
		branding = map[string]interface{}{}
	}

	if settings == nil {
		settings = map[string]interface{}{}
	}

	return &Tenant{
		ID:        rec.ID,
		Name:      rec.Name,
		Slug:      rec.Slug,
		Tier:      rec.Tier,
		Status:    rec.Status,
		Branding:  branding,
		Settings:  settings,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}
