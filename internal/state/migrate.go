package state

import (
	"fmt"
	"os"
	"path/filepath"
)

// RunMigrations runs all migrations in the migrations directory
func RunMigrations(store Store, migrationsPath string) error {
	// Check if migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	// Find all migration files and sort them
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsPath)
	}

	// Run the initial migration (for now we only have 0001_init.sql)
	// The store.Migrate method will read and execute it
	if err := store.Migrate(migrationsPath); err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	return nil
}

