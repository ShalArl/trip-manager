package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/weather-info/internal/fetcher"
	"github.com/redis/go-redis/v9"
)

const weatherPrefix = "weather:"

type WeatherCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewWeatherCache(redisURL string, ttlHours int) (*WeatherCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &WeatherCache{
		client: client,
		ttl:    time.Duration(ttlHours) * time.Hour,
	}, nil
}

func cacheKey(lat, lng float64) string {
	return fmt.Sprintf("%s%.2f:%.2f", weatherPrefix, lat, lng)
}

func (c *WeatherCache) Set(ctx context.Context, weather *fetcher.WeatherResponse) error {
	data, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("marshal weather: %w", err)
	}
	return c.client.Set(ctx, cacheKey(weather.Latitude, weather.Longitude), data, c.ttl).Err()
}

func (c *WeatherCache) Get(ctx context.Context, lat, lng float64) (*fetcher.WeatherResponse, error) {
	data, err := c.client.Get(ctx, cacheKey(lat, lng)).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache Miss
	}
	if err != nil {
		return nil, fmt.Errorf("get weather: %w", err)
	}

	var weather fetcher.WeatherResponse
	if err := json.Unmarshal(data, &weather); err != nil {
		return nil, fmt.Errorf("unmarshal weather: %w", err)
	}
	return &weather, nil
}

func (c *WeatherCache) Close() error {
	return c.client.Close()
}
