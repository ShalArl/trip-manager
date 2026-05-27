package main

import (
	"context"
	"log"
	"time"

	"github.com/ShalArl/trip-manager/backend/travel-warning/config"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/cache"
	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/fetcher"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	warningCache, err := cache.NewWarningCache(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("failed to connect to redis %v", err)
	}

	defer func(warningCache *cache.WarningCache) {
		err := warningCache.Close()
		if err != nil {
			log.Fatalf("failed to close cache %v", err)
		}
	}(warningCache)

	client := fetcher.NewClient(cfg.APIUrl)

	log.Println("Fetching travel warnings...")
	start := time.Now()

	warnings, err := client.FetchAll(ctx)
	if err != nil {
		log.Fatalf("failed to fetch travel warnings: %v", err)
	}

	var cached []*cache.CachedWarning
	for _, w := range warnings {
		cached = append(cached, &cache.CachedWarning{
			CountryCode:    w.CountryCode,
			CountryName:    w.CountryName,
			Level:          w.Level(),
			Warning:        w.Warning,
			PartialWarning: w.PartialWarning,
			UpdatedAt:      time.Now(),
		})
	}

	if err := warningCache.SetBatch(ctx, cached); err != nil {
		log.Fatalf("failed to cache travel warnings: %v", err)
	}

	log.Printf("Fetched %d warnings in %s", len(warnings), time.Since(start))
}