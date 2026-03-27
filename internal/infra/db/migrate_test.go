package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestShouldSkipMigration_UserIDMigrationAlreadyApplied(t *testing.T) {
	pool := setupMigrationPool(t)

	mustExecMigrationSQL(t, pool.DB, `
		CREATE TABLE blocked_product (
			activity_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			create_time TEXT NOT NULL
		);
		CREATE TABLE notification_config (
			activity_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			target_price REAL NOT NULL,
			last_notify_time TEXT,
			create_time TEXT NOT NULL,
			update_time TEXT NOT NULL
		);
		CREATE TABLE user_settings (
			user_id TEXT PRIMARY KEY,
			bark_key TEXT NOT NULL,
			create_time TEXT NOT NULL,
			update_time TEXT NOT NULL
		);
	`)

	skip, err := pool.shouldSkipMigration("002_add_user_id.sql")
	if err != nil {
		t.Fatalf("shouldSkipMigration() error = %v", err)
	}
	if !skip {
		t.Fatal("expected 002_add_user_id.sql to be skipped when schema is already migrated")
	}
}

func TestShouldSkipMigration_UserIDMigrationPending(t *testing.T) {
	pool := setupMigrationPool(t)

	mustExecMigrationSQL(t, pool.DB, `
		CREATE TABLE blocked_product (
			activity_id TEXT NOT NULL,
			create_time TEXT NOT NULL
		);
		CREATE TABLE notification_config (
			activity_id TEXT NOT NULL,
			target_price REAL NOT NULL,
			last_notify_time TEXT,
			create_time TEXT NOT NULL,
			update_time TEXT NOT NULL
		);
	`)

	skip, err := pool.shouldSkipMigration("002_add_user_id.sql")
	if err != nil {
		t.Fatalf("shouldSkipMigration() error = %v", err)
	}
	if skip {
		t.Fatal("expected 002_add_user_id.sql to run when schema is still pending")
	}
}

func setupMigrationPool(t *testing.T) *Pool {
	t.Helper()

	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/migrate.db?mode=rwc")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &Pool{DB: db}
}

func mustExecMigrationSQL(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec migration sql failed: %v\nquery:\n%s", err, query)
	}
}
