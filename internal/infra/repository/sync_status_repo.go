package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"

	db "kbfood/internal/infra/db"
)

type syncStatusRepository struct {
	db *db.Pool
}

// NewSyncStatusRepository creates a new sync status repository
func NewSyncStatusRepository(db *db.Pool) repository.SyncStatusRepository {
	return &syncStatusRepository{db: db}
}

func (r *syncStatusRepository) Upsert(ctx context.Context, status *entity.SyncStatus) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	lastRunTime := status.LastRunTime.Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO sync_status (job_name, last_run_time, status, product_count, error_message, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(job_name) DO UPDATE SET
			last_run_time = excluded.last_run_time,
			status = excluded.status,
			product_count = excluded.product_count,
			error_message = excluded.error_message,
			updated_at = excluded.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		status.JobName,
		lastRunTime,
		status.Status,
		status.ProductCount,
		status.ErrorMessage,
		now,
	)
	if err != nil {
		return fmt.Errorf("upsert sync status: %w", err)
	}

	return nil
}

func (r *syncStatusRepository) GetLatest(ctx context.Context, jobName string) (*entity.SyncStatus, error) {
	query := `
		SELECT id, job_name, last_run_time, status, product_count, error_message, created_at, updated_at
		FROM sync_status
		WHERE job_name = ?
	`

	row := r.db.QueryRowContext(ctx, query, jobName)

	var status entity.SyncStatus
	var lastRunTime, createdAt, updatedAt string

	err := row.Scan(
		&status.ID,
		&status.JobName,
		&lastRunTime,
		&status.Status,
		&status.ProductCount,
		&status.ErrorMessage,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get sync status: %w", err)
	}

	// Parse timestamps
	status.LastRunTime, _ = time.Parse("2006-01-02 15:04:05", lastRunTime)
	status.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	status.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &status, nil
}
