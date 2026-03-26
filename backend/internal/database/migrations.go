package database

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
)

// RunMigrations executes all migration files in the migrations directory
func RunMigrations(db *sqlx.DB) error {
	// Get the migrations directory - relative to the executable or module root
	migrationsDir := "migrations"
	
	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try from backend directory if running from there
		migrationsDir = "backend/migrations"
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found")
		}
	}

	// Read all .sql files from migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort migration files to ensure consistent order
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		log.Println("No migration files found")
		return nil
	}

	log.Printf("Found %d migration files\n", len(migrationFiles))

	// Execute each migration file
	for _, migrationFile := range migrationFiles {
		migrationPath := filepath.Join(migrationsDir, migrationFile)
		log.Printf("Running migration: %s\n", migrationFile)

		// Read migration file
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
		}

		// Execute migration
		if _, err := db.Exec(string(migrationSQL)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationFile, err)
		}

		log.Printf("✓ Migration %s completed\n", migrationFile)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// RunMigrationsFromFS runs migrations from an embedded file system (useful for compiled binaries)
func RunMigrationsFromFS(db *sqlx.DB, migrationsFS fs.FS) error {
	// Get all .sql files from the embedded FS
	files, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations from FS: %w", err)
	}

	var migrationNames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationNames = append(migrationNames, file.Name())
		}
	}

	// Sort to ensure consistent order
	sort.Strings(migrationNames)

	if len(migrationNames) == 0 {
		log.Println("No migration files found in embedded FS")
		return nil
	}

	log.Printf("Found %d migration files\n", len(migrationNames))

	// Execute each migration
	for _, migrationName := range migrationNames {
		log.Printf("Running migration: %s\n", migrationName)

		// Read migration from FS
		migrationSQL, err := fs.ReadFile(migrationsFS, migrationName)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s from FS: %w", migrationName, err)
		}

		// Execute migration
		if _, err := db.Exec(string(migrationSQL)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
		}

		log.Printf("✓ Migration %s completed\n", migrationName)
	}

	log.Println("All migrations completed successfully")
	return nil
}

