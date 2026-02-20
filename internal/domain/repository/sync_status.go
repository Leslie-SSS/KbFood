package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// SyncStatusRepository handles sync status persistence
type SyncStatusRepository interface {
	// Upsert creates or updates a sync status record
	Upsert(ctx context.Context, status *entity.SyncStatus) error
	// GetLatest retrieves the latest sync status for a job
	GetLatest(ctx context.Context, jobName string) (*entity.SyncStatus, error)
}
