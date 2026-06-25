package tenant

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GCMMetricsClient struct {
	ProjectID string
}

func NewGCMMetricsClient(projectID string) *GCMMetricsClient {
	return &GCMMetricsClient{ProjectID: projectID}
}

func (c *GCMMetricsClient) QueryAPICallsByService(ctx context.Context, tenantID string) (map[string]int64, error) {
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring client: %w", err)
	}
	defer client.Close()

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", c.ProjectID),
		Filter: fmt.Sprintf(`metric.type="custom.googleapis.com/trip_manager/api_calls_total" AND metric.labels.tenant_id="%s"`, tenantID),
		Interval: &monitoringpb.TimeInterval{
			StartTime: timestamppb.New(startOfMonth),
			EndTime:   timestamppb.New(now),
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	services := map[string]int64{}
	it := client.ListTimeSeries(ctx, req)
	for {
		ts, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to query metrics: %w", err)
		}
		svc := ts.Metric.Labels["service"]
		for _, point := range ts.Points {
			services[svc] += point.Value.GetInt64Value()
		}
	}
	return services, nil
}
