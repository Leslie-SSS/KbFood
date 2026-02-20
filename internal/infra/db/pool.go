package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbSqlc "kbfood/internal/infra/db/sqlc"

	_ "modernc.org/sqlite"
)

// Pool wraps sql.DB with additional methods
type Pool struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	// SQLite file path (e.g., "file:./food.db?mode=rwcb")
	Path            string        `envconfig:"PATH" default:"file:./data/food.db?mode=rwc"`
	MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"MAX_LIFETIME" default:"1h"`
	ConnMaxIdleTime time.Duration `envconfig:"MAX_IDLE_TIME" default:"30m"`
}

// NewPool creates a new database connection pool
func NewPool(ctx context.Context, cfg *Config) (*Pool, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Enable foreign keys and WAL mode for better performance
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("enable WAL mode: %w", err)
	}

	// Verify connection works
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	pool := &Pool{DB: db}

	// Run migrations
	if err := pool.RunMigrations(); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return pool, nil
}

// Ping checks if the database connection is alive
func (p *Pool) Ping(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

// Close closes the database connection pool
func (p *Pool) Close() {
	p.DB.Close()
}

// Queries returns a new Queries instance bound to this pool
func (p *Pool) Queries() *dbSqlc.Queries {
	return dbSqlc.New(p.DB)
}
