package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type notificationRepository struct {
	db *db.Queries
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *db.Queries) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) FindByActivityID(ctx context.Context, activityID string, userID string) (*entity.NotificationConfig, error) {
	config, err := r.db.GetNotification(ctx, db.GetNotificationParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return convertDBNotificationToEntity(&config), nil
}

func (r *notificationRepository) ListByUser(ctx context.Context, userID string) ([]*entity.NotificationConfig, error) {
	configs, err := r.db.ListNotificationsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list notifications by user: %w", err)
	}

	result := make([]*entity.NotificationConfig, len(configs))
	for i, c := range configs {
		result[i] = convertDBNotificationToEntity(&c)
	}
	return result, nil
}

func (r *notificationRepository) ListAll(ctx context.Context) ([]*entity.NotificationConfig, error) {
	configs, err := r.db.ListAllNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all notifications: %w", err)
	}

	result := make([]*entity.NotificationConfig, len(configs))
	for i, c := range configs {
		result[i] = convertDBNotificationToEntity(&c)
	}
	return result, nil
}

func (r *notificationRepository) Upsert(ctx context.Context, config *entity.NotificationConfig) error {
	params := db.UpsertNotificationParams{
		ActivityID:     config.ActivityID,
		UserID:         config.UserID,
		TargetPrice:    config.TargetPrice,
		LastNotifyTime: sqlNullStringFromTimePtr(config.LastNotifyTime),
	}

	err := r.db.UpsertNotification(ctx, params)
	if err != nil {
		return fmt.Errorf("upsert notification: %w", err)
	}
	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, activityID string, userID string) error {
	err := r.db.DeleteNotification(ctx, db.DeleteNotificationParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

func (r *notificationRepository) UpdateNotifyTime(ctx context.Context, activityID string, userID string) error {
	err := r.db.UpdateNotificationNotifyTime(ctx, db.UpdateNotificationNotifyTimeParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("update notify time: %w", err)
	}
	return nil
}

// convertDBNotificationToEntity converts db.NotificationConfig to entity.NotificationConfig
func convertDBNotificationToEntity(c *db.NotificationConfig) *entity.NotificationConfig {
	var lastNotifyTime *time.Time
	if c.LastNotifyTime.Valid && c.LastNotifyTime.String != "" {
		if t := parseSQLiteTime(c.LastNotifyTime.String); !t.IsZero() {
			lastNotifyTime = &t
		}
	}

	return &entity.NotificationConfig{
		ActivityID:     c.ActivityID,
		UserID:         c.UserID,
		TargetPrice:    c.TargetPrice,
		LastNotifyTime: lastNotifyTime,
		CreateTime:     parseSQLiteTime(c.CreateTime),
		UpdateTime:     parseSQLiteTime(c.UpdateTime),
	}
}
