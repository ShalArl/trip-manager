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

	if _, err := tx.ExecContext(ctx, "SET LOCAL app.tenant_id = $1", tenantID); err != nil {
		err := tx.Rollback()
		if err != nil {
			return fmt.Errorf("failed to rollback transaction after setting tenant_id: %w", err)
		}
		return fmt.Errorf("failed to set tenant_id: %w", err)
	}

	if err := fn(tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return err
	}

	return tx.Commit()
}
