package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// NotificationRepository defines the interface for notification config data access
type NotificationRepository interface {
	// FindByActivityID finds a notification config by activity ID and user ID
	FindByActivityID(ctx context.Context, activityID string, userID string) (*entity.NotificationConfig, error)

	// ListByUser lists all notification configs for a user
	ListByUser(ctx context.Context, userID string) ([]*entity.NotificationConfig, error)

	// ListAll lists all notification configs (for background notification checker)
	ListAll(ctx context.Context) ([]*entity.NotificationConfig, error)

	// Upsert creates or updates a notification config
	Upsert(ctx context.Context, config *entity.NotificationConfig) error

	// Delete deletes a notification config by activity ID and user ID
	Delete(ctx context.Context, activityID string, userID string) error

	// UpdateNotifyTime updates the last notification time
	UpdateNotifyTime(ctx context.Context, activityID string, userID string) error
}
