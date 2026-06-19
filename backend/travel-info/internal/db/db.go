package db

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Location struct {
	Lat float64 `db:"lat"`
	Lng float64 `db:"lng"`
}

type LocationsDB struct {
	db *sqlx.DB
}

func NewLocationsDB(dbURL string) (*LocationsDB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to locations db: %w", err)
	}
	return &LocationsDB{db: db}, nil
}

func (l *LocationsDB) GetUniqueLocations(ctx context.Context) ([]Location, error) {
	var locations []Location
	err := l.db.SelectContext(ctx, &locations, `
		SELECT DISTINCT
			ROUND(latitude::numeric, 2)::float8 AS lat,
			ROUND(longitude::numeric, 2)::float8 AS lng
		FROM locations
		WHERE latitude IS NOT NULL
		  AND longitude IS NOT NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("query locations: %w", err)
	}
	return locations, nil
}

func (l *LocationsDB) Close() error {
	return l.db.Close()
}

func RoundCoord(coord float64) float64 {
	return math.Round(coord*100) / 100
}
