package repository

import (
	"context"

	"kbfood/internal/domain/entity"
)

// CandidateRepository defines the interface for candidate item data access
type CandidateRepository interface {
	// FindByID finds a candidate by ID
	FindByID(ctx context.Context, id int64) (*entity.CandidateItem, error)

	// FindByRegion finds all candidates in a region
	FindByRegion(ctx context.Context, region string) ([]*entity.CandidateItem, error)

	// FindByGroupKey finds candidates by group key and region
	FindByGroupKey(ctx context.Context, groupKey, region string) (*entity.CandidateItem, error)

	// ListAll lists all candidates
	ListAll(ctx context.Context) ([]*entity.CandidateItem, error)

	// Create creates a new candidate
	Create(ctx context.Context, candidate *entity.CandidateItem) error

	// Update updates an existing candidate
	Update(ctx context.Context, candidate *entity.CandidateItem) error

	// Delete deletes a candidate by ID
	Delete(ctx context.Context, id int64) error

	// DeleteByIDs deletes candidates by IDs
	DeleteByIDs(ctx context.Context, ids []int64) error
}
