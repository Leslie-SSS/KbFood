package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// MasterProductRepository defines the interface for master product data access
type MasterProductRepository interface {
	// FindByID finds a master product by ID
	FindByID(ctx context.Context, id string) (*entity.MasterProduct, error)

	// FindByRegion finds all master products in a region
	FindByRegion(ctx context.Context, region string) ([]*entity.MasterProduct, error)

	// FindByPlatform finds all master products for a platform
	FindByPlatform(ctx context.Context, platform string) ([]*entity.MasterProduct, error)

	// FindByRegionAndPlatform finds all master products in a region for a platform
	FindByRegionAndPlatform(ctx context.Context, region, platform string) ([]*entity.MasterProduct, error)

	// ListAll lists all master products
	ListAll(ctx context.Context) ([]*entity.MasterProduct, error)

	// Create creates a new master product
	Create(ctx context.Context, product *entity.MasterProduct) error

	// Update updates an existing master product
	Update(ctx context.Context, product *entity.MasterProduct) error

	// Delete deletes a master product by ID
	Delete(ctx context.Context, id string) error
}
