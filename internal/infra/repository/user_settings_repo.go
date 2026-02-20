package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type userSettingsRepository struct {
	db *db.Queries
}

// NewUserSettingsRepository creates a new user settings repository
func NewUserSettingsRepository(db *db.Queries) repository.UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

func (r *userSettingsRepository) Get(ctx context.Context, userID string) (*entity.UserSettings, error) {
	settings, err := r.db.GetUserSettings(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user settings: %w", err)
	}

	return convertDBUserSettingsToEntity(&settings), nil
}

func (r *userSettingsRepository) Upsert(ctx context.Context, settings *entity.UserSettings) error {
	err := r.db.UpsertUserSettings(ctx, db.UpsertUserSettingsParams{
		UserID:  settings.UserID,
		BarkKey: settings.BarkKey,
	})
	if err != nil {
		return fmt.Errorf("upsert user settings: %w", err)
	}
	return nil
}

// convertDBUserSettingsToEntity converts db.UserSetting to entity.UserSettings
func convertDBUserSettingsToEntity(s *db.UserSetting) *entity.UserSettings {
	return &entity.UserSettings{
		UserID:     s.UserID,
		BarkKey:    s.BarkKey,
		CreateTime: parseSQLiteTime(s.CreateTime),
		UpdateTime: parseSQLiteTime(s.UpdateTime),
	}
}
