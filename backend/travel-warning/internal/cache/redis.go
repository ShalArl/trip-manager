package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ShalArl/trip-manager/backend/travel-warning/internal/fetcher"
	"github.com/redis/go-redis/v9"
)

const (
	warningPrefix = "warning:"
	defaultTTL    = 25 * time.Hour
)

type WarningCache struct {
	client *redis.Client
}

type CachedWarning struct {
	CountryCode    string               `json:"countryCode"`
	CountryName    string               `json:"countryName"`
	Level          fetcher.WarningLevel `json:"level"`
	Warning        bool                 `json:"warning"`
	PartialWarning bool                 `json:"partialWarning"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}

func NewWarningCache(redisURL string) (*WarningCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing Redis URL %s: %w", redisURL, err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis %s: %w", redisURL, err)
	}

	return &WarningCache{client: client}, nil
}

func (c *WarningCache) Set(ctx context.Context, w *CachedWarning) error {
	data, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("error serializing cached warning: %w", err)
	}

	key := warningPrefix + w.CountryCode
	return c.client.Set(ctx, key, data, defaultTTL).Err()
}

func (c *WarningCache) Get(ctx context.Context, countryCode string) (*CachedWarning, error) {
	key := warningPrefix + countryCode
	data, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting cached warning: %w", err)
	}

	var warning CachedWarning
	if err := json.Unmarshal(data, &warning); err != nil {
		return nil, fmt.Errorf("error deserializing cached warning: %w", err)
	}

	return &warning, nil
}

func (c *WarningCache) SetBatch(ctx context.Context, warnings []*CachedWarning) error {
	pipe := c.client.Pipeline()
	for _, warning := range warnings {
		data, err := json.Marshal(warning)
		if err != nil {
			continue
		}
		key := warningPrefix + warning.CountryCode
		pipe.Set(ctx, key, data, defaultTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (c *WarningCache) Close() error {
	return c.client.Close()
}
