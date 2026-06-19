package main

import (
	"context"
	"log"
	"time"

	"github.com/ShalArl/trip-manager/backend/travel-info/config"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/db"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/fetcher"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 1. Gemeinsamen Redis-Cache initialisieren
	sharedCache, err := cache.NewCache(cfg.RedisUrl, cfg.WeatherCacheTTLHours)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer sharedCache.Close()

	// ==========================================
	// PART 1: TRAVEL WARNINGS WORKER
	// ==========================================
	log.Println("--- Starting Travel Warnings Sync ---")
	travelClient := fetcher.NewWarningClient(cfg.TravelWarningUrl)

	startTravel := time.Now()
	warnings, err := travelClient.FetchAll(ctx)
	if err != nil {
		log.Printf("ERROR: failed to fetch travel warnings: %v", err)
	} else {
		var cachedWarnings []*cache.CachedWarning
		for _, w := range warnings {
			cachedWarnings = append(cachedWarnings, &cache.CachedWarning{
				CountryCode:    w.CountryCode,
				CountryName:    w.CountryName,
				Level:          w.Level(),
				Warning:        w.Warning,
				PartialWarning: w.PartialWarning,
				UpdatedAt:      time.Now(),
			})
		}

		if err := sharedCache.SetWarningBatch(ctx, cachedWarnings); err != nil {
			log.Printf("ERROR: failed to cache travel warnings: %v", err)
		} else {
			log.Printf("SUCCESS: Fetched and cached %d travel warnings in %s", len(warnings), time.Since(startTravel))
		}
	}

	// ==========================================
	// PART 2: WEATHER INFO WORKER
	// ==========================================
	log.Println("--- Starting Weather Info Sync ---")

	locationsDB, err := db.NewLocationsDB(cfg.LocationsDBURL)
	if err != nil {
		log.Fatalf("failed to connect to locations db: %v", err)
	}
	defer locationsDB.Close()

	meteoClient := fetcher.NewOpenMeteoClient(cfg.WeatherAPIUrl, cfg.WeatherForecastDays)

	log.Println("Fetching unique locations from DB...")
	locations, err := locationsDB.GetUniqueLocations(ctx)
	if err != nil {
		log.Fatalf("failed to get locations: %v", err)
	}
	log.Printf("Found %d unique locations", len(locations))

	if len(locations) == 0 {
		log.Println("No locations found, weather sync skipped")
	} else {
		success := 0
		failed := 0

		for _, loc := range locations {
			weather, err := meteoClient.FetchForecast(ctx, loc.Lat, loc.Lng, "")
			if err != nil {
				log.Printf("failed to fetch weather for lat=%.2f lng=%.2f: %v", loc.Lat, loc.Lng, err)
				failed++
				continue
			}

			if err := sharedCache.SetWeather(ctx, weather); err != nil {
				log.Printf("failed to cache weather for lat=%.2f lng=%.2f: %v", loc.Lat, loc.Lng, err)
				failed++
				continue
			}

			success++
			time.Sleep(100 * time.Millisecond) // Rate limiting Schutz
		}
		log.Printf("Weather cache updated: %d success, %d failed", success, failed)
	}

	log.Println("--- All worker sync tasks completed ---")
}
