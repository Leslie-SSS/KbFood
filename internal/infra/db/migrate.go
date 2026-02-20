package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// RunMigrations executes all migration files from the migrations directory
func (p *Pool) RunMigrations() error {
	migrationsDir := "/app/migrations"
	
	log.Info().Str("dir", migrationsDir).Msg("Running database migrations")

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Warn().Msg("Migrations directory not found, skipping")
		return nil
	}

	// Read all migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	// Sort by filename (001_, 002_, etc.)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	executed := 0
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		log.Info().Str("file", entry.Name()).Msg("Executing migration")
		if _, err := p.Exec(string(content)); err != nil {
			// Log but don't fail on duplicate column/table errors (migration already run)
			if strings.Contains(err.Error(), "duplicate column") || 
			   strings.Contains(err.Error(), "already exists") {
				log.Warn().Str("file", entry.Name()).Err(err).Msg("Migration already applied, skipping")
				continue
			}
			return fmt.Errorf("execute migration %s: %w", entry.Name(), err)
		}
		executed++
	}

	log.Info().Int("count", executed).Msg("Migrations completed")
	return nil
}
