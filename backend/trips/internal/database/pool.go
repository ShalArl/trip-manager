package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PoolManager struct {
	mu              sync.RWMutex
	pools           map[string]*sqlx.DB
	defaultDB       *sqlx.DB
	usersServiceURL string
	internalSecret  string
}

func NewPoolManager(defaultDB *sqlx.DB, usersServiceURL, internalSecret string) *PoolManager {
	return &PoolManager{
		pools:           make(map[string]*sqlx.DB),
		defaultDB:       defaultDB,
		usersServiceURL: usersServiceURL,
		internalSecret:  internalSecret,
	}
}

func (p *PoolManager) GetDB(ctx context.Context, tenantID string) *sqlx.DB {
	if tenantID == "" || tenantID == "default" {
		return p.defaultDB
	}

	// Cache prüfen
	p.mu.RLock()
	if db, ok := p.pools[tenantID]; ok {
		p.mu.RUnlock()
		return db
	}
	p.mu.RUnlock()

	// Enterprise DB-URL vom users-Service holen
	dbURL, err := p.fetchEnterpriseDBURL(ctx, tenantID)
	if err != nil {
		return p.defaultDB
	}

	// Neue Verbindung erstellen
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return p.defaultDB
	}

	p.mu.Lock()
	p.pools[tenantID] = db
	p.mu.Unlock()

	return db
}

func (p *PoolManager) fetchEnterpriseDBURL(ctx context.Context, tenantID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/internal/tenants/%s/db-url", p.usersServiceURL, tenantID), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Internal-Secret", p.internalSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("not an enterprise tenant")
	}

	var result struct {
		DbURL string `json:"dbUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.DbURL, nil
}
