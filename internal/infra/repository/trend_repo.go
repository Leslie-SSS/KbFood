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

type trendRepository struct {
	db *db.Queries
}

// NewTrendRepository creates a new trend repository
func NewTrendRepository(db *db.Queries) repository.TrendRepository {
	return &trendRepository{db: db}
}

func (r *trendRepository) FindByActivityIDAndDate(ctx context.Context, activityID string, date time.Time) (*entity.PriceTrend, error) {
	params := db.GetTrendByActivityIDAndDateParams{
		ActivityID: activityID,
		RecordDate: dateToSQLite(date),
	}
	trend, err := r.db.GetTrendByActivityIDAndDate(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get trend: %w", err)
	}
	return convertDBTrendToEntity(&trend), nil
}

func (r *trendRepository) FindByActivityID(ctx context.Context, activityID string) ([]*entity.PriceTrend, error) {
	trends, err := r.db.ListTrendsByActivityID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("list trends: %w", err)
	}

	result := make([]*entity.PriceTrend, len(trends))
	for i, t := range trends {
		result[i] = convertDBTrendToEntity(&t)
	}
	return result, nil
}

func (r *trendRepository) Create(ctx context.Context, trend *entity.PriceTrend) error {
	return r.Upsert(ctx, trend)
}

func (r *trendRepository) Upsert(ctx context.Context, trend *entity.PriceTrend) error {
	params := db.CreateTrendParams{
		ActivityID: trend.ActivityID,
		Price:      trend.Price,
		RecordDate: dateToSQLite(trend.RecordDate),
	}

	err := r.db.CreateTrend(ctx, params)
	if err != nil {
		return fmt.Errorf("upsert trend: %w", err)
	}
	return nil
}

func (r *trendRepository) DeleteByActivityIDs(ctx context.Context, activityIDs []string) error {
	// Delete one at a time since SQLite doesn't support array parameters
	for _, id := range activityIDs {
		if err := r.db.DeleteTrendsByActivityIDs(ctx, id); err != nil {
			return fmt.Errorf("delete trend %s: %w", id, err)
		}
	}
	return nil
}

// convertDBTrendToEntity converts db.ProductPriceTrend to entity.PriceTrend
func convertDBTrendToEntity(t *db.ProductPriceTrend) *entity.PriceTrend {
	return &entity.PriceTrend{
		ID:         t.ID,
		ActivityID: t.ActivityID,
		Price:      t.Price,
		RecordDate: parseSQLiteDate(t.RecordDate),
		CreateTime: parseSQLiteTime(t.CreateTime),
	}
}
