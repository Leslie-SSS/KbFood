package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	db "kbfood/internal/infra/db/sqlc"
)

type candidateRepository struct {
	db *db.Queries
}

// NewCandidateRepository creates a new candidate repository
func NewCandidateRepository(db *db.Queries) repository.CandidateRepository {
	return &candidateRepository{db: db}
}

func (r *candidateRepository) FindByID(ctx context.Context, id int64) (*entity.CandidateItem, error) {
	candidate, err := r.db.GetCandidateByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get candidate: %w", err)
	}
	return convertDBCandidateToEntity(&candidate), nil
}

func (r *candidateRepository) FindByRegion(ctx context.Context, region string) ([]*entity.CandidateItem, error) {
	candidates, err := r.db.ListCandidatesByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("list candidates by region: %w", err)
	}

	result := make([]*entity.CandidateItem, len(candidates))
	for i, c := range candidates {
		result[i] = convertDBCandidateToEntity(&c)
	}
	return result, nil
}

func (r *candidateRepository) FindByGroupKey(ctx context.Context, groupKey, region string) (*entity.CandidateItem, error) {
	candidates, err := r.db.ListCandidatesByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("list candidates by region: %w", err)
	}

	for _, c := range candidates {
		if c.GroupKey == groupKey {
			return convertDBCandidateToEntity(&c), nil
		}
	}
	return nil, nil
}

func (r *candidateRepository) ListAll(ctx context.Context) ([]*entity.CandidateItem, error) {
	candidates, err := r.db.ListAllCandidates(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all candidates: %w", err)
	}

	result := make([]*entity.CandidateItem, len(candidates))
	for i, c := range candidates {
		result[i] = convertDBCandidateToEntity(&c)
	}
	return result, nil
}

func (r *candidateRepository) Create(ctx context.Context, candidate *entity.CandidateItem) error {
	titleVotesJSON, err := json.Marshal(candidate.TitleVotes)
	if err != nil {
		return fmt.Errorf("marshal title votes: %w", err)
	}

	params := db.CreateCandidateParams{
		GroupKey:         candidate.GroupKey,
		Region:           candidate.Region,
		TitleVotes:       string(titleVotesJSON),
		TotalOccurrences: sqlNullInt64FromInt(candidate.TotalOccurrences),
		LastPrice:        sqlNullFloat64FromFloat(candidate.LastPrice),
		LastStatus:       sqlNullInt64FromInt(candidate.LastStatus),
		FirstSeenTime:    timeToSQLite(candidate.FirstSeenTime),
		LastSeenTime:     timeToSQLite(candidate.LastSeenTime),
	}

	err = r.db.CreateCandidate(ctx, params)
	if err != nil {
		return fmt.Errorf("create candidate: %w", err)
	}
	return nil
}

func (r *candidateRepository) Update(ctx context.Context, candidate *entity.CandidateItem) error {
	titleVotesJSON, err := json.Marshal(candidate.TitleVotes)
	if err != nil {
		return fmt.Errorf("marshal title votes: %w", err)
	}

	params := db.UpdateCandidateParams{
		ID:               candidate.ID,
		TitleVotes:       string(titleVotesJSON),
		TotalOccurrences: sqlNullInt64FromInt(candidate.TotalOccurrences),
		LastPrice:        sqlNullFloat64FromFloat(candidate.LastPrice),
		LastStatus:       sqlNullInt64FromInt(candidate.LastStatus),
		LastSeenTime:     timeToSQLite(candidate.LastSeenTime),
	}

	err = r.db.UpdateCandidate(ctx, params)
	if err != nil {
		return fmt.Errorf("update candidate: %w", err)
	}
	return nil
}

func (r *candidateRepository) Delete(ctx context.Context, id int64) error {
	err := r.db.DeleteCandidate(ctx, id)
	if err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}
	return nil
}

func (r *candidateRepository) DeleteByIDs(ctx context.Context, ids []int64) error {
	// Note: DeleteCandidatesByIDs now deletes one at a time
	// This is a limitation of SQLite/sqlc without array support
	for _, id := range ids {
		if err := r.db.DeleteCandidate(ctx, id); err != nil {
			return fmt.Errorf("delete candidate %d: %w", id, err)
		}
	}
	return nil
}

// parseSQLiteTime parses an RFC3339 timestamp string to time.Time
func parseSQLiteTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// convertDBCandidateToEntity converts db.CandidateItem to entity.CandidateItem
func convertDBCandidateToEntity(c *db.CandidateItem) *entity.CandidateItem {
	titleVotes := make(map[string]int)
	if c.TitleVotes != "" && c.TitleVotes != "{}" {
		_ = json.Unmarshal([]byte(c.TitleVotes), &titleVotes)
	}

	return &entity.CandidateItem{
		ID:               c.ID,
		GroupKey:         c.GroupKey,
		Region:           c.Region,
		TitleVotes:       titleVotes,
		TotalOccurrences: int(int64FromNull(c.TotalOccurrences)),
		LastPrice:        float64FromNull(c.LastPrice),
		LastStatus:       int(int64FromNull(c.LastStatus)),
		FirstSeenTime:    parseSQLiteTime(c.FirstSeenTime),
		LastSeenTime:     parseSQLiteTime(c.LastSeenTime),
		CreateTime:       parseSQLiteTime(c.CreateTime),
		UpdateTime:       parseSQLiteTime(c.UpdateTime),
	}
}
