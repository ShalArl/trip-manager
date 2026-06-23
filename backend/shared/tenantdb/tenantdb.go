package tenantdb

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const tenantKey contextKey = "tenantId"

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

func GetTenantID(ctx context.Context) string {
	if v, ok := ctx.Value(tenantKey).(string); ok && v != "" {
		return v
	}
	return "default"
}

func WithTenant(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tenantID := GetTenantID(ctx)
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, execErr := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID)); execErr != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to set tenant_id: %w", execErr)
	}

	if fnErr := fn(tx); fnErr != nil {
		_ = tx.Rollback()
		return fnErr
	}

	return tx.Commit()
}
