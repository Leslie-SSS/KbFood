// Package test provides test helpers for E2E and integration tests
package test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"kbfood/internal/config"
	"kbfood/internal/infra/db"
)

// TestDB wraps an in-memory SQLite database
type TestDB struct {
	Pool    *db.Pool
	DataDir string
}

// SetupTestDB creates a new in-memory SQLite test database
func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	// Create a temporary directory for the test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create connection pool
	dbCfg := &db.Config{
		Path:            "file:" + dbPath + "?mode=rwc",
		MaxOpenConns:    5,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}

	pool, err := db.NewPool(ctx, dbCfg)
	if err != nil {
		t.Fatalf("Failed to create db pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	// Run migrations
	runMigrations(t, pool)

	return &TestDB{
		Pool:    pool,
		DataDir: tmpDir,
	}
}

// runMigrations runs the database schema migrations
func runMigrations(t *testing.T, pool *db.Pool) {
	schemaSQL := `
	-- 商品主表
	CREATE TABLE IF NOT EXISTS product (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		activity_id TEXT NOT NULL UNIQUE,
		platform TEXT,
		region TEXT,
		title TEXT,
		shop_name TEXT,
		original_price REAL,
		current_price REAL,
		sales_status INTEGER,
		activity_create_time TEXT,
		create_time TEXT NOT NULL DEFAULT (datetime('now')),
		update_time TEXT NOT NULL DEFAULT (datetime('now'))
	);

	-- 标准商品库 (DT专用)
	CREATE TABLE IF NOT EXISTS master_product (
		id TEXT PRIMARY KEY,
		region TEXT NOT NULL,
		standard_title TEXT NOT NULL,
		price REAL,
		status INTEGER,
		trust_score INTEGER DEFAULT 0,
		create_time TEXT NOT NULL DEFAULT (datetime('now')),
		update_time TEXT NOT NULL DEFAULT (datetime('now'))
	);

	-- 候选商品池
	CREATE TABLE IF NOT EXISTS candidate_item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_key TEXT NOT NULL,
		region TEXT NOT NULL,
		title_votes TEXT NOT NULL DEFAULT '{}',
		total_occurrences INTEGER DEFAULT 0,
		last_price REAL,
		last_status INTEGER,
		first_seen_time TEXT NOT NULL DEFAULT (datetime('now')),
		last_seen_time TEXT NOT NULL DEFAULT (datetime('now')),
		create_time TEXT NOT NULL DEFAULT (datetime('now')),
		update_time TEXT NOT NULL DEFAULT (datetime('now'))
	);

	-- 屏蔽商品
	CREATE TABLE IF NOT EXISTS blocked_product (
		activity_id TEXT PRIMARY KEY,
		create_time TEXT NOT NULL DEFAULT (datetime('now'))
	);

	-- 通知配置
	CREATE TABLE IF NOT EXISTS notification_config (
		activity_id TEXT PRIMARY KEY,
		target_price REAL NOT NULL,
		last_notify_time TEXT,
		create_time TEXT NOT NULL DEFAULT (datetime('now')),
		update_time TEXT NOT NULL DEFAULT (datetime('now'))
	);

	-- 价格趋势
	CREATE TABLE IF NOT EXISTS product_price_trend (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		activity_id TEXT NOT NULL,
		price REAL NOT NULL,
		record_date TEXT NOT NULL,
		create_time TEXT NOT NULL DEFAULT (datetime('now')),
		UNIQUE(activity_id, record_date)
	);
	`

	if _, err := pool.Exec(schemaSQL); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// SetupTestConfig creates a test configuration
func SetupTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port:            8080,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Database: config.DatabaseConfig{
			Path:            ":memory:",
			MaxOpenConns:    5,
			MaxIdleConns:    1,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
		},
		Log: config.LogConfig{
			Level:  "debug",
			Format: "console",
		},
	}
}

// TruncateTables truncates all tables for a clean test state
func TruncateTables(t *testing.T, pool *db.Pool) {
	tables := []string{
		"product_price_trend",
		"notification_config",
		"blocked_product",
		"product",
		"master_product",
		"candidate_item",
	}

	for _, table := range tables {
		_, err := pool.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// AssertNoDBError checks if error is a "no rows" error (expected in some tests)
func AssertNoDBError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected database error: %v", err)
	}
}
