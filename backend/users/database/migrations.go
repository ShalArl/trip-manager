package database

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func RunMigrations(db *sqlx.DB) error {
	subFS, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to access migrations: %w", err)
	}

	files, err := fs.ReadDir(subFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			names = append(names, f.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		log.Printf("Running migration: %s", name)
		sql, err := fs.ReadFile(subFS, name)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", name, err)
		}
		if _, err := db.Exec(string(sql)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", name, err)
		}
		log.Printf("✓ %s completed", name)
	}
	return nil
}
