package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// ProductFilter holds filter parameters for product queries
type ProductFilter struct {
	Keyword         string
	Platform        string
	Region          string
	SalesStatus     *int
	MonitorStatus   *int
	RecentSevenDays bool
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	// FindByID finds a product by ID
	FindByID(ctx context.Context, id int64) (*entity.Product, error)

	// FindByActivityID finds a product by activity ID
	FindByActivityID(ctx context.Context, activityID string) (*entity.Product, error)

	// FindByFilter finds products by filter criteria
	FindByFilter(ctx context.Context, filter ProductFilter) ([]*entity.Product, error)

	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error

	// Update updates an existing product
	Update(ctx context.Context, product *entity.Product) error

	// UpdateByActivityID updates a product by activity ID
	UpdateByActivityID(ctx context.Context, activityID string, product *entity.Product) error

	// Delete deletes a product by ID
	Delete(ctx context.Context, id int64) error

	// DeleteByActivityIDs deletes products by activity IDs
	DeleteByActivityIDs(ctx context.Context, activityIDs []string) error

	// DeleteByPlatform deletes all products for a platform
	DeleteByPlatform(ctx context.Context, platform string) error

	// CountByPlatform counts products by platform
	CountByPlatform(ctx context.Context, platform string) (int64, error)

	// ListBlocked lists all blocked products
	ListBlocked(ctx context.Context) ([]*entity.Product, error)
}
