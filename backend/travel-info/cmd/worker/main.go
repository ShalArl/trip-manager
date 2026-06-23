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

	// ==========================================
	// PART 2: WEATHER INFO WORKER
	// ==========================================

}
