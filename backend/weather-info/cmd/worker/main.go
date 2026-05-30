package main

import (
	"context"
	"log"
	"time"

	"github.com/ShalArl/trip-manager/backend/weather-info/config"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/db"
	"github.com/ShalArl/trip-manager/backend/weather-info/internal/fetcher"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := context.Background()

	// Redis Cache
	weatherCache, err := cache.NewWeatherCache(cfg.RedisUrl, cfg.CacheTTLHours)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer weatherCache.Close()

	// Locations DB
	locationsDB, err := db.NewLocationsDB(cfg.LocationsDBURL)
	if err != nil {
		log.Fatalf("failed to connect to locations db: %v", err)
	}
	defer locationsDB.Close()

	// Open-Meteo Client
	meteoClient := fetcher.NewClient(cfg.APIUrl, cfg.ForecastDays)

	// Alle unique Locations aus DB holen
	log.Println("Fetching unique locations from DB...")
	locations, err := locationsDB.GetUniqueLocations(ctx)
	if err != nil {
		log.Fatalf("failed to get locations: %v", err)
	}
	log.Printf("Found %d unique locations", len(locations))

	if len(locations) == 0 {
		log.Println("No locations found, nothing to cache")
		return
	}

	// Für jede Location Wetterdaten holen und cachen
	success := 0
	failed := 0

	for _, loc := range locations {
		weather, err := meteoClient.FetchForecast(ctx, loc.Lat, loc.Lng, "")
		if err != nil {
			log.Printf("failed to fetch weather for lat=%.2f lng=%.2f: %v", loc.Lat, loc.Lng, err)
			failed++
			continue
		}

		if err := weatherCache.Set(ctx, weather); err != nil {
			log.Printf("failed to cache weather for lat=%.2f lng=%.2f: %v", loc.Lat, loc.Lng, err)
			failed++
			continue
		}

		success++

		// Rate limiting – Open-Meteo erlaubt 10.000 Requests/Tag
		// Bei vielen Locations kurz warten
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Weather cache updated: %d success, %d failed", success, failed)
}
