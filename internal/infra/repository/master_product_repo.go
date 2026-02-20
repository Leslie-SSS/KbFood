package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type masterProductRepository struct {
	db *db.Queries
}

// NewMasterProductRepository creates a new master product repository
func NewMasterProductRepository(db *db.Queries) repository.MasterProductRepository {
	return &masterProductRepository{db: db}
}

func (r *masterProductRepository) FindByID(ctx context.Context, id string) (*entity.MasterProduct, error) {
	master, err := r.db.GetMasterProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get master product: %w", err)
	}
	return convertDBMasterProductToEntity(&master), nil
}

func (r *masterProductRepository) FindByRegion(ctx context.Context, region string) ([]*entity.MasterProduct, error) {
	masters, err := r.db.ListMasterProductsByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("list master products by region: %w", err)
	}

	result := make([]*entity.MasterProduct, len(masters))
	for i, m := range masters {
		result[i] = convertDBMasterProductToEntity(&m)
	}
	return result, nil
}

func (r *masterProductRepository) FindByPlatform(ctx context.Context, platform string) ([]*entity.MasterProduct, error) {
	masters, err := r.db.ListMasterProductsByPlatform(ctx, sqlNullString(platform))
	if err != nil {
		return nil, fmt.Errorf("list master products by platform: %w", err)
	}

	result := make([]*entity.MasterProduct, len(masters))
	for i, m := range masters {
		result[i] = convertDBMasterProductToEntity(&m)
	}
	return result, nil
}

func (r *masterProductRepository) FindByRegionAndPlatform(ctx context.Context, region, platform string) ([]*entity.MasterProduct, error) {
	params := db.ListMasterProductsByRegionAndPlatformParams{
		Region:   region,
		Platform: sqlNullString(platform),
	}
	masters, err := r.db.ListMasterProductsByRegionAndPlatform(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list master products by region and platform: %w", err)
	}

	result := make([]*entity.MasterProduct, len(masters))
	for i, m := range masters {
		result[i] = convertDBMasterProductToEntity(&m)
	}
	return result, nil
}

func (r *masterProductRepository) ListAll(ctx context.Context) ([]*entity.MasterProduct, error) {
	masters, err := r.db.ListAllMasterProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all master products: %w", err)
	}

	result := make([]*entity.MasterProduct, len(masters))
	for i, m := range masters {
		result[i] = convertDBMasterProductToEntity(&m)
	}
	return result, nil
}

func (r *masterProductRepository) Create(ctx context.Context, product *entity.MasterProduct) error {
	params := db.CreateMasterProductParams{
		ID:            product.ID,
		Region:        product.Region,
		Platform:      sqlNullString(product.Platform),
		StandardTitle: product.StandardTitle,
		Price:         sqlNullFloat64FromFloat(product.Price),
		Status:        sqlNullInt64FromInt(product.Status),
		TrustScore:    sqlNullInt64FromInt(product.TrustScore),
	}

	err := r.db.CreateMasterProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("create master product: %w", err)
	}
	return nil
}

func (r *masterProductRepository) Update(ctx context.Context, product *entity.MasterProduct) error {
	params := db.UpdateMasterProductParams{
		ID:         product.ID,
		Price:      sqlNullFloat64FromFloat(product.Price),
		Status:     sqlNullInt64FromInt(product.Status),
		TrustScore: sqlNullInt64FromInt(product.TrustScore),
	}

	err := r.db.UpdateMasterProduct(ctx, params)
	if err != nil {
		return fmt.Errorf("update master product: %w", err)
	}
	return nil
}

func (r *masterProductRepository) Delete(ctx context.Context, id string) error {
	err := r.db.DeleteMasterProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("delete master product: %w", err)
	}
	return nil
}

// convertDBMasterProductToEntity converts db.MasterProduct to entity.MasterProduct
func convertDBMasterProductToEntity(m *db.MasterProduct) *entity.MasterProduct {
	return &entity.MasterProduct{
		ID:            m.ID,
		Region:        m.Region,
		Platform:      m.Platform.String,
		StandardTitle: m.StandardTitle,
		Price:         float64FromNull(m.Price),
		Status:        int(int64FromNull(m.Status)),
		TrustScore:    int(int64FromNull(m.TrustScore)),
		CreateTime:    parseSQLiteTime(m.CreateTime),
		UpdateTime:    parseSQLiteTime(m.UpdateTime),
	}
}
