package repository

import (
	"context"
	"fmt"

	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type blockedRepository struct {
	db *db.Queries
}

// NewBlockedRepository creates a new blocked repository
func NewBlockedRepository(db *db.Queries) repository.BlockedRepository {
	return &blockedRepository{db: db}
}

func (r *blockedRepository) Exists(ctx context.Context, activityID string, userID string) (bool, error) {
	exists, err := r.db.ExistsBlockedProduct(ctx, db.ExistsBlockedProductParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return false, fmt.Errorf("check blocked product exists: %w", err)
	}
	return exists, nil
}

func (r *blockedRepository) Create(ctx context.Context, activityID string, userID string) error {
	err := r.db.CreateBlockedProduct(ctx, db.CreateBlockedProductParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("create blocked product: %w", err)
	}
	return nil
}

func (r *blockedRepository) Delete(ctx context.Context, activityID string, userID string) error {
	err := r.db.DeleteBlockedProduct(ctx, db.DeleteBlockedProductParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("delete blocked product: %w", err)
	}
	return nil
}

func (r *blockedRepository) List(ctx context.Context, userID string) ([]string, error) {
	blocked, err := r.db.ListBlockedProductsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list blocked products: %w", err)
	}
	return blocked, nil
}
