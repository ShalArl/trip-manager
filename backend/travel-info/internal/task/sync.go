package task

import (
	"context"
	"log"
	"time"

	"github.com/ShalArl/trip-manager/backend/travel-info/config"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/db"
	"github.com/ShalArl/trip-manager/backend/travel-info/internal/fetcher"
)

func RunSyncTasks(ctx context.Context, sharedCache *cache.Cache, cfg *config.Config, meteoClient *fetcher.OpenMeteoClient) {
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

	log.Println("--- Starting Weather Info Sync ---")

	locationsDB, err := db.NewLocationsDB(cfg.LocationsDBURL)
	if err != nil {
		log.Printf("ERROR: failed to connect to locations db: %v", err)
		return // Weiche Landung: Task abbrechen, API lebt weiter
	}
	defer locationsDB.Close()

	log.Println("Fetching unique locations from DB...")
	locations, err := locationsDB.GetUniqueLocations(ctx)
	if err != nil {
		log.Printf("ERROR: failed to get locations from db: %v", err)
		return // Weiche Landung
	}
	log.Printf("Found %d unique locations", len(locations))

	if len(locations) == 0 {
		log.Println("No locations found, weather sync skipped")
	} else {
		success := 0
		failed := 0

		for _, loc := range locations {
			// Falls die API runterfährt, Loop sofort abbrechen
			select {
			case <-ctx.Done():
				log.Println("Sync aborted: context cancelled")
				return
			default:
			}

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
			time.Sleep(100 * time.Millisecond)
		}
		log.Printf("Weather cache updated: %d success, %d failed", success, failed)
	}

	log.Println("--- All worker sync tasks completed ---")
}
