package repository

import (
	"context"
	"time"

	"kbfood/internal/domain/entity"
)

// TrendRepository defines the interface for price trend data access
type TrendRepository interface {
	// FindByActivityIDAndDate finds a trend by activity ID and date
	FindByActivityIDAndDate(ctx context.Context, activityID string, date time.Time) (*entity.PriceTrend, error)

	// FindByActivityID finds all trends for an activity ID
	FindByActivityID(ctx context.Context, activityID string) ([]*entity.PriceTrend, error)

	// Create creates a new trend record
	Create(ctx context.Context, trend *entity.PriceTrend) error

	// Upsert creates or updates a trend record (keeps the lowest price)
	Upsert(ctx context.Context, trend *entity.PriceTrend) error

	// DeleteByActivityIDs deletes trends by activity IDs
	DeleteByActivityIDs(ctx context.Context, activityIDs []string) error
}

// BlockedRepository defines the interface for blocked product data access
type BlockedRepository interface {
	// Exists checks if a product is blocked for a user
	Exists(ctx context.Context, activityID string, userID string) (bool, error)

	// Create blocks a product for a user
	Create(ctx context.Context, activityID string, userID string) error

	// Delete unblocks a product for a user
	Delete(ctx context.Context, activityID string, userID string) error

	// List lists all blocked activity IDs for a user
	List(ctx context.Context, userID string) ([]string, error)
}
