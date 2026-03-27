package db

import (
	"database/sql"
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

		shouldSkip, err := p.shouldSkipMigration(entry.Name())
		if err != nil {
			return fmt.Errorf("check migration %s: %w", entry.Name(), err)
		}
		if shouldSkip {
			log.Info().Str("file", entry.Name()).Msg("Migration already satisfied by current schema, skipping")
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

func (p *Pool) shouldSkipMigration(filename string) (bool, error) {
	switch filename {
	case "002_add_user_id.sql":
		return p.multiUserSchemaReady()
	default:
		return false, nil
	}
}

func (p *Pool) multiUserSchemaReady() (bool, error) {
	blockedHasUserID, err := p.tableHasColumn("blocked_product", "user_id")
	if err != nil {
		return false, err
	}
	notificationHasUserID, err := p.tableHasColumn("notification_config", "user_id")
	if err != nil {
		return false, err
	}
	userSettingsExists, err := p.tableExists("user_settings")
	if err != nil {
		return false, err
	}

	return blockedHasUserID && notificationHasUserID && userSettingsExists, nil
}

func (p *Pool) tableExists(tableName string) (bool, error) {
	var exists int
	err := p.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM sqlite_master
			WHERE type = 'table' AND name = ?
		)
	`, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check table %s exists: %w", tableName, err)
	}
	return exists == 1, nil
}

func (p *Pool) tableHasColumn(tableName, columnName string) (bool, error) {
	if ok, err := p.tableExists(tableName); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	rows, err := p.Query(fmt.Sprintf("PRAGMA table_info(%s)", quoteIdentifier(tableName)))
	if err != nil {
		return false, fmt.Errorf("inspect table %s columns: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			dataType   string
			notNull    int
			defaultV   sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultV, &primaryKey); err != nil {
			return false, fmt.Errorf("scan table info for %s: %w", tableName, err)
		}
		if name == columnName {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate table info for %s: %w", tableName, err)
	}

	return false, nil
}

func quoteIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}
