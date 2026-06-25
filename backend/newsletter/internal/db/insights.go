package db

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TopDestination struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	TripCount   int64   `json:"tripCount"`
	AvgLikes    float64 `json:"avgLikes"`
}

type EngagementStats struct {
	TotalLikes      int64   `json:"totalLikes"`
	TotalComments   int64   `json:"totalComments"`
	AvgLikesPerTrip float64 `json:"avgLikesPerTrip"`
}

type SeasonalPattern struct {
	PeakMonth           string `json:"peakMonth"`
	AvgPlanningLeadDays int    `json:"avgPlanningLeadDays"`
}

type TenantInsights struct {
	TenantID        string           `json:"tenantId"`
	TopDestinations []TopDestination `json:"topDestinations"`
	Engagement      EngagementStats  `json:"engagement"`
	SeasonalPattern SeasonalPattern  `json:"seasonalPattern"`
	GeneratedAt     time.Time        `json:"generatedAt"`
}

func GenerateInsightsForTenant(ctx context.Context, driver neo4j.DriverWithContext, tenantID string) (*TenantInsights, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Top Destinationen
	topDest, err := queryTopDestinations(ctx, session, tenantID)
	if err != nil {
		return nil, fmt.Errorf("top destinations: %w", err)
	}

	// Engagement Stats
	engagement, err := queryEngagement(ctx, session, tenantID)
	if err != nil {
		return nil, fmt.Errorf("engagement: %w", err)
	}

	// Saisonale Muster
	seasonal, err := querySeasonalPattern(ctx, session, tenantID)
	if err != nil {
		return nil, fmt.Errorf("seasonal: %w", err)
	}

	return &TenantInsights{
		TenantID:        tenantID,
		TopDestinations: topDest,
		Engagement:      engagement,
		SeasonalPattern: seasonal,
		GeneratedAt:     time.Now(),
	}, nil
}

func queryTopDestinations(ctx context.Context, session neo4j.SessionWithContext, tenantID string) ([]TopDestination, error) {
	result, err := session.Run(ctx, `
        MATCH (t:Trip {tenantId: $tenantId})
        WHERE t.country IS NOT NULL
        WITH t.country AS country, t.countryCode AS countryCode, count(t) AS tripCount
        OPTIONAL MATCH (trip:Trip {country: country, tenantId: $tenantId})<-[:LIKES]-()
        WITH country, countryCode, tripCount, count(*) AS totalLikes
        RETURN country, countryCode, tripCount, 
               CASE WHEN tripCount > 0 THEN toFloat(totalLikes) / tripCount ELSE 0.0 END AS avgLikes
        ORDER BY tripCount DESC
        LIMIT 10
    `, map[string]any{"tenantId": tenantID})
	if err != nil {
		return nil, err
	}

	var destinations []TopDestination
	for result.Next(ctx) {
		record := result.Record()
		country, _ := record.Get("country")
		countryCode, _ := record.Get("countryCode")
		tripCount, _ := record.Get("tripCount")
		avgLikes, _ := record.Get("avgLikes")

		destinations = append(destinations, TopDestination{
			Country:     fmt.Sprintf("%v", country),
			CountryCode: fmt.Sprintf("%v", countryCode),
			TripCount:   tripCount.(int64),
			AvgLikes:    avgLikes.(float64),
		})
	}
	return destinations, nil
}

func queryEngagement(ctx context.Context, session neo4j.SessionWithContext, tenantID string) (EngagementStats, error) {
	result, err := session.Run(ctx, `
        MATCH (t:Trip {tenantId: $tenantId})
        OPTIONAL MATCH (t)<-[:LIKES]-(u)
        OPTIONAL MATCH (t)<-[:COMMENTED_ON]-(c)
        WITH count(DISTINCT t) AS tripCount,
             count(DISTINCT u) AS totalLikes,
             count(DISTINCT c) AS totalComments
        RETURN totalLikes, totalComments,
               CASE WHEN tripCount > 0 THEN toFloat(totalLikes) / tripCount ELSE 0.0 END AS avgLikesPerTrip
    `, map[string]any{"tenantId": tenantID})
	if err != nil {
		return EngagementStats{}, err
	}

	if result.Next(ctx) {
		record := result.Record()
		totalLikes, _ := record.Get("totalLikes")
		totalComments, _ := record.Get("totalComments")
		avgLikes, _ := record.Get("avgLikesPerTrip")
		return EngagementStats{
			TotalLikes:      totalLikes.(int64),
			TotalComments:   totalComments.(int64),
			AvgLikesPerTrip: avgLikes.(float64),
		}, nil
	}
	return EngagementStats{}, nil
}

func querySeasonalPattern(ctx context.Context, session neo4j.SessionWithContext, tenantID string) (SeasonalPattern, error) {
	result, err := session.Run(ctx, `
        MATCH (t:Trip {tenantId: $tenantId})
        WHERE t.startDate IS NOT NULL AND t.createdAt IS NOT NULL
        WITH t,
             duration.between(date(t.createdAt), date(t.startDate)).days AS leadDays,
             date(t.createdAt).month AS month
        WITH avg(leadDays) AS avgLead,
             month, count(*) AS monthCount
        ORDER BY monthCount DESC
        LIMIT 1
        RETURN month, toInteger(avgLead) AS avgLeadDays
    `, map[string]any{"tenantId": tenantID})
	if err != nil {
		return SeasonalPattern{}, err
	}

	months := []string{"", "Januar", "Februar", "März", "April", "Mai", "Juni",
		"Juli", "August", "September", "Oktober", "November", "Dezember"}

	if result.Next(ctx) {
		record := result.Record()
		month, _ := record.Get("month")
		leadDays, _ := record.Get("avgLeadDays")
		monthIdx := int(month.(int64))
		peakMonth := ""
		if monthIdx >= 1 && monthIdx <= 12 {
			peakMonth = months[monthIdx]
		}
		return SeasonalPattern{
			PeakMonth:           peakMonth,
			AvgPlanningLeadDays: int(leadDays.(int64)),
		}, nil
	}
	return SeasonalPattern{}, nil
}
