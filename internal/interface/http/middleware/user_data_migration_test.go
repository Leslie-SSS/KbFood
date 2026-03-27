package middleware

import (
	"context"
	"database/sql"
	"testing"

	dbinfra "kbfood/internal/infra/db"

	_ "modernc.org/sqlite"
)

func TestMigrateLegacyUserData_MigratesLegacyRows(t *testing.T) {
	pool := setupMigrationTestDB(t)
	ctx := context.Background()

	mustExec(t, pool.DB, `
		INSERT INTO user_settings (user_id, bark_key, create_time, update_time)
		VALUES ('legacy-user', 'LEGACY_BARK', '2026-03-17 02:52:15', '2026-03-21 19:16:57')
	`)
	mustExec(t, pool.DB, `
		INSERT INTO notification_config (activity_id, user_id, target_price, last_notify_time, create_time, update_time)
		VALUES ('DT_1', 'legacy-user', 68.7, NULL, '2026-03-17 02:52:15', '2026-03-21 19:16:57')
	`)
	mustExec(t, pool.DB, `
		INSERT INTO blocked_product (activity_id, user_id, create_time)
		VALUES ('DT_2', 'legacy-user', '2026-03-17 02:52:15')
	`)

	if err := migrateLegacyUserData(ctx, pool, "client-123", "legacy-user"); err != nil {
		t.Fatalf("migrateLegacyUserData() error = %v", err)
	}

	assertCount(t, pool.DB, "SELECT COUNT(*) FROM user_settings WHERE user_id = 'legacy-user'", 0)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM notification_config WHERE user_id = 'legacy-user'", 0)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM blocked_product WHERE user_id = 'legacy-user'", 0)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM user_settings WHERE user_id = 'client-123' AND bark_key = 'LEGACY_BARK'", 1)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM notification_config WHERE user_id = 'client-123' AND activity_id = 'DT_1'", 1)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM blocked_product WHERE user_id = 'client-123' AND activity_id = 'DT_2'", 1)
}

func TestMigrateLegacyUserData_ResolvesConflictsByRecency(t *testing.T) {
	pool := setupMigrationTestDB(t)
	ctx := context.Background()

	mustExec(t, pool.DB, `
		INSERT INTO user_settings (user_id, bark_key, create_time, update_time)
		VALUES
			('client-123', '', '2026-03-20 00:00:00', '2026-03-20 00:00:00'),
			('legacy-user', 'LEGACY_BARK', '2026-03-17 00:00:00', '2026-03-21 00:00:00')
	`)
	mustExec(t, pool.DB, `
		INSERT INTO notification_config (activity_id, user_id, target_price, last_notify_time, create_time, update_time)
		VALUES
			('DT_SHARED', 'client-123', 99.0, NULL, '2026-03-20 00:00:00', '2026-03-20 00:00:00'),
			('DT_SHARED', 'legacy-user', 68.7, '2026-03-21 08:00:00', '2026-03-17 00:00:00', '2026-03-21 09:00:00'),
			('DT_ONLY_LEGACY', 'legacy-user', 12.3, NULL, '2026-03-17 00:00:00', '2026-03-21 09:00:00')
	`)
	mustExec(t, pool.DB, `
		INSERT INTO blocked_product (activity_id, user_id, create_time)
		VALUES
			('DT_BLOCK_CURRENT', 'client-123', '2026-03-20 00:00:00'),
			('DT_BLOCK_LEGACY', 'legacy-user', '2026-03-17 00:00:00')
	`)

	if err := migrateLegacyUserData(ctx, pool, "client-123", "legacy-user"); err != nil {
		t.Fatalf("migrateLegacyUserData() error = %v", err)
	}

	assertCount(t, pool.DB, "SELECT COUNT(*) FROM user_settings WHERE user_id = 'client-123' AND bark_key = 'LEGACY_BARK'", 1)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM notification_config WHERE user_id = 'client-123'", 2)
	assertCount(t, pool.DB, "SELECT COUNT(*) FROM blocked_product WHERE user_id = 'client-123'", 2)

	var targetPrice float64
	if err := pool.DB.QueryRowContext(ctx, `
		SELECT target_price FROM notification_config
		WHERE user_id = 'client-123' AND activity_id = 'DT_SHARED'
	`).Scan(&targetPrice); err != nil {
		t.Fatalf("scan merged target price: %v", err)
	}
	if targetPrice != 68.7 {
		t.Fatalf("expected merged target price 68.7, got %v", targetPrice)
	}
}

func setupMigrationTestDB(t *testing.T) *dbinfra.Pool {
	t.Helper()

	sqlDB, err := sql.Open("sqlite", "file:"+t.TempDir()+"/migration.db?mode=rwc")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	pool := &dbinfra.Pool{DB: sqlDB}
	mustExec(t, pool.DB, `
		CREATE TABLE user_settings (
			user_id TEXT PRIMARY KEY,
			bark_key TEXT NOT NULL,
			create_time TEXT NOT NULL DEFAULT (datetime('now')),
			update_time TEXT NOT NULL DEFAULT (datetime('now'))
		);

		CREATE TABLE notification_config (
			activity_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			target_price REAL NOT NULL,
			last_notify_time TEXT,
			create_time TEXT NOT NULL DEFAULT (datetime('now')),
			update_time TEXT NOT NULL DEFAULT (datetime('now')),
			PRIMARY KEY (activity_id, user_id)
		);

		CREATE TABLE blocked_product (
			activity_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			create_time TEXT NOT NULL DEFAULT (datetime('now')),
			PRIMARY KEY (activity_id, user_id)
		);
	`)

	return pool
}

func mustExec(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec query failed: %v\nquery:\n%s", err, query)
	}
}

func assertCount(t *testing.T, db *sql.DB, query string, want int) {
	t.Helper()

	var got int
	if err := db.QueryRow(query).Scan(&got); err != nil {
		t.Fatalf("query count failed: %v\nquery:\n%s", err, query)
	}
	if got != want {
		t.Fatalf("count mismatch for query %q: got %d want %d", query, got, want)
	}
}
