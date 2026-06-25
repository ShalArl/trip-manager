package tenant

import "context"

type MetricsClient interface {
	QueryAPICallsByService(ctx context.Context, tenantID string) (map[string]int64, error)
}
