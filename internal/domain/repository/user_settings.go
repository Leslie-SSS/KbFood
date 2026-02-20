package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// UserSettingsRepository defines the interface for user settings data access
type UserSettingsRepository interface {
	// Get retrieves user settings by user ID
	Get(ctx context.Context, userID string) (*entity.UserSettings, error)

	// Upsert creates or updates user settings
	Upsert(ctx context.Context, settings *entity.UserSettings) error
}
