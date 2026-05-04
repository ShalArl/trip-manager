package service

import (
	"sync"
	"time"
)

// urlCacheEntry holds a URL with its expiration time
type urlCacheEntry struct {
	url       string
	expiresAt time.Time
}

// URLCache is a thread-safe cache for download URLs with automatic cleanup
type URLCache struct {
	mu       sync.RWMutex
	cache    map[string]urlCacheEntry
	ttl      time.Duration
	stopChan chan struct{}
}

// NewURLCache creates a new URL cache with the specified TTL and starts cleanup goroutine
func NewURLCache(ttl time.Duration) *URLCache {
	uc := &URLCache{
		cache:    make(map[string]urlCacheEntry),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine that runs every ttl/4
	go uc.cleanupLoop(ttl / 4)

	return uc
}

// Get retrieves a URL from cache if it exists and hasn't expired
func (uc *URLCache) Get(key string) (string, bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if entry, ok := uc.cache[key]; ok && time.Now().Before(entry.expiresAt) {
		return entry.url, true
	}
	return "", false
}

// Set stores a URL in cache with expiration
func (uc *URLCache) Set(key string, url string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.cache[key] = urlCacheEntry{
		url:       url,
		expiresAt: time.Now().Add(uc.ttl),
	}
}

// Clear removes all entries from cache
func (uc *URLCache) Clear() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.cache = make(map[string]urlCacheEntry)
}

// Stop stops the cleanup goroutine
func (uc *URLCache) Stop() {
	close(uc.stopChan)
}

// cleanupLoop periodically removes expired entries
func (uc *URLCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			uc.cleanup()
		case <-uc.stopChan:
			return
		}
	}
}

// cleanup removes all expired entries from the cache
func (uc *URLCache) cleanup() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	now := time.Now()
	for key, entry := range uc.cache {
		if now.After(entry.expiresAt) {
			delete(uc.cache, key)
		}
	}
}
