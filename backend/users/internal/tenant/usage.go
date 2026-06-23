package tenant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tenantdb"
	"time"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
)

type UsageResponse struct {
	TenantID  string         `json:"tenantId"`
	Period    string         `json:"period"`
	APICalls  int64          `json:"apiCalls"`
	Breakdown []ServiceUsage `json:"breakdown"`
	Pricing   PricingInfo    `json:"pricing"`
}

type ServiceUsage struct {
	Service string `json:"service"`
	Calls   int64  `json:"calls"`
}

type PricingInfo struct {
	Tier        string  `json:"tier"`
	BasePrice   float64 `json:"basePrice"`
	APICallCost float64 `json:"apiCallCost"`
	TotalCost   float64 `json:"totalCost"`
	Currency    string  `json:"currency"`
}

func GetUsageHandler(repo Repository, prometheusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		if tenantID == "" || tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}

		// Prometheus abfragen – API Calls pro Service für diesen Tenant
		query := fmt.Sprintf(
			`sum by (service) (trip_manager_api_calls_total{tenant_id="%s"})`,
			tenantID,
		)

		result, err := queryPrometheus(prometheusURL, query)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to query metrics: %v", err))
			return
		}

		var totalCalls int64
		var breakdown []ServiceUsage
		for _, r := range result {
			svc := r.Metric["service"]
			calls := r.Value
			totalCalls += calls
			breakdown = append(breakdown, ServiceUsage{
				Service: svc,
				Calls:   calls,
			})
		}

		// Pricing berechnen
		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		tenant, err := repo.GetByID(ctx, tenantID)
		if err != nil {
			respondError(w, http.StatusNotFound, "tenant not found")
			return
		}

		pricing := calculatePricing(tenant.Tier, totalCalls)
		
		respondJSON(w, http.StatusOK, UsageResponse{
			TenantID:  tenantID,
			Period:    time.Now().Format("2006-01"),
			APICalls:  totalCalls,
			Breakdown: breakdown,
			Pricing:   pricing,
		})
	}
}

type prometheusResult struct {
	Metric map[string]string
	Value  int64
}

func queryPrometheus(baseURL, query string) ([]prometheusResult, error) {
	params := url.Values{}
	params.Set("query", query)

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/query?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var out []prometheusResult
	for _, r := range result.Data.Result {
		var calls int64
		if len(r.Value) > 1 {
			fmt.Sscanf(fmt.Sprintf("%v", r.Value[1]), "%d", &calls)
		}
		out = append(out, prometheusResult{
			Metric: r.Metric,
			Value:  calls,
		})
	}
	return out, nil
}

func calculatePricing(tier string, apiCalls int64) PricingInfo {
	const (
		basePrice      = 9.0
		freeCallsLimit = 10000
		pricePerCall   = 0.001
	)

	var apiCallCost float64
	if apiCalls > freeCallsLimit {
		apiCallCost = float64(apiCalls-freeCallsLimit) * pricePerCall
	}

	return PricingInfo{
		Tier:        tier,
		BasePrice:   basePrice,
		APICallCost: apiCallCost,
		TotalCost:   basePrice + apiCallCost,
		Currency:    "EUR",
	}
}
