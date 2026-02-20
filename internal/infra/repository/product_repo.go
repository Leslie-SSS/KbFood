package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type productRepository struct {
	db *db.Queries
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *db.Queries) repository.ProductRepository {
	return &productRepository{db: db}
}

// FindByID finds a product by ID
func (r *productRepository) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	// Since we use activity_id as the primary key, this method is not supported
	// It's better to fail fast than return a misleading error
	return nil, nil
}

// FindByActivityID finds a product by activity ID
func (r *productRepository) FindByActivityID(ctx context.Context, activityID string) (*entity.Product, error) {
	product, err := r.db.GetProductByActivityID(ctx, activityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get product: %w", err)
	}
	return convertDBProductToEntity(&product), nil
}

// FindByFilter finds products by filter criteria
func (r *productRepository) FindByFilter(ctx context.Context, filter repository.ProductFilter) ([]*entity.Product, error) {
	var salesStatus sql.NullInt64
	if filter.SalesStatus != nil {
		salesStatus = sqlNullInt64FromInt(*filter.SalesStatus)
	}

	var recentDays sql.NullInt64
	if filter.RecentSevenDays {
		recentDays = sqlNullInt64FromInt(1)
	}

	products, err := r.db.ListProducts(ctx, db.ListProductsParams{
		Column1:     sqlNullString(filter.Platform),
		Platform:    sqlNullString(filter.Platform),
		Column3:     sqlNullString(filter.Region),
		Region:      sqlNullString(filter.Region),
		Column5:     sqlNullInt64FromInt(1), // IS NULL check placeholder
		SalesStatus: salesStatus,
		Column7:     sqlNullInt64FromInt(1), // IS NULL check placeholder
		Column8:     recentDays,
	})
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	result := make([]*entity.Product, len(products))
	for i, p := range products {
		result[i] = convertDBProductToEntity(&p)
	}
	return result, nil
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	params := db.CreateProductParams{
		ActivityID:         product.ActivityID,
		Platform:           sqlNullString(product.Platform),
		Region:             sqlNullString(product.Region),
		Title:              sqlNullString(product.Title),
		ShopName:           sqlNullString(product.ShopName),
		OriginalPrice:      sqlNullFloat64FromFloat(product.OriginalPrice),
		CurrentPrice:       sqlNullFloat64FromFloat(product.CurrentPrice),
		SalesStatus:        sqlNullInt64FromInt(product.SalesStatus),
		ActivityCreateTime: sqlNullStringFromTime(product.ActivityCreateTime),
	}

	err := r.db.CreateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	params := db.UpdateProductParams{
		ID:           product.ID,
		CurrentPrice: sqlNullFloat64FromFloat(product.CurrentPrice),
		SalesStatus:  sqlNullInt64FromInt(product.SalesStatus),
	}

	err := r.db.UpdateProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

// UpdateByActivityID updates a product by activity ID
func (r *productRepository) UpdateByActivityID(ctx context.Context, activityID string, product *entity.Product) error {
	params := db.UpdateProductByActivityIDParams{
		ActivityID:   activityID,
		CurrentPrice: sqlNullFloat64FromFloat(product.CurrentPrice),
		SalesStatus:  sqlNullInt64FromInt(product.SalesStatus),
	}

	err := r.db.UpdateProductByActivityID(ctx, params)
	if err != nil {
		return fmt.Errorf("update product by activity id: %w", err)
	}
	return nil
}

// Delete deletes a product by ID
func (r *productRepository) Delete(ctx context.Context, id int64) error {
	err := r.db.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}

// DeleteByActivityIDs deletes products by activity IDs
func (r *productRepository) DeleteByActivityIDs(ctx context.Context, activityIDs []string) error {
	// Delete one at a time since SQLite doesn't support array parameters
	for _, id := range activityIDs {
		if err := r.db.DeleteByActivityIDs(ctx, id); err != nil {
			return fmt.Errorf("delete product %s: %w", id, err)
		}
	}
	return nil
}

// DeleteByPlatform deletes all products for a platform
func (r *productRepository) DeleteByPlatform(ctx context.Context, platform string) error {
	err := r.db.DeleteByPlatform(ctx, sqlNullString(platform))
	if err != nil {
		return fmt.Errorf("delete products by platform: %w", err)
	}
	return nil
}

// CountByPlatform counts products by platform
func (r *productRepository) CountByPlatform(ctx context.Context, platform string) (int64, error) {
	count, err := r.db.CountByPlatform(ctx, sqlNullString(platform))
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}
	return count, nil
}

// ListBlocked lists all blocked products
func (r *productRepository) ListBlocked(ctx context.Context) ([]*entity.Product, error) {
	products, err := r.db.ListProductsWithBlockedStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("list blocked products: %w", err)
	}

	result := make([]*entity.Product, len(products))
	for i, p := range products {
		result[i] = convertDBProductToEntity(&p)
	}
	return result, nil
}

// convertDBProductToEntity converts db.Product to entity.Product
func convertDBProductToEntity(p *db.Product) *entity.Product {
	return &entity.Product{
		ID:                 p.ID,
		ActivityID:         p.ActivityID,
		Platform:           stringFromNull(p.Platform),
		Region:             stringFromNull(p.Region),
		Title:              stringFromNull(p.Title),
		ShopName:           stringFromNull(p.ShopName),
		OriginalPrice:      float64FromNull(p.OriginalPrice),
		CurrentPrice:       float64FromNull(p.CurrentPrice),
		SalesStatus:        int(int64FromNull(p.SalesStatus)),
		ActivityCreateTime: parseSQLiteTime(stringFromNull(p.ActivityCreateTime)),
		CreateTime:         parseSQLiteTime(p.CreateTime),
		UpdateTime:         parseSQLiteTime(p.UpdateTime),
	}
}
