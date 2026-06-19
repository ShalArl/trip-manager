package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ServiceMetrics struct {
	apiCallsCounter metric.Int64Counter
	serviceName     string
}

func NewServiceMetrics(meter metric.Meter, serviceName string) (*ServiceMetrics, error) {
	counter, err := meter.Int64Counter(
		"api_calls_total",
		metric.WithDescription("Total number of API calls"),
	)
	if err != nil {
		return nil, err
	}

	return &ServiceMetrics{
		apiCallsCounter: counter,
		serviceName:     serviceName,
	}, nil
}

func (m *ServiceMetrics) RecordAPICall(ctx context.Context, tenantID, endpoint, method string) {
	m.apiCallsCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("tenant_id", tenantID),
			attribute.String("service", m.serviceName),
			attribute.String("endpoint", endpoint),
			attribute.String("method", method),
		),
	)
}
