package service

import (
	"context"
	"fmt"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"

	"github.com/rs/zerolog/log"
)

// truncateToDay truncates a time to the start of the day (midnight)
// This ensures consistent date comparison for trend recording
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// DataCleaningService handles candidate pool management and promotion
type DataCleaningService struct {
	masterRepo     repository.MasterProductRepository
	candidateRepo  repository.CandidateRepository
	trendRepo      repository.TrendRepository
	priceValidator *PriceValidator
	titleCleaner   *TitleCleaner
}

// NewDataCleaningService creates a new data cleaning service
func NewDataCleaningService(
	masterRepo repository.MasterProductRepository,
	candidateRepo repository.CandidateRepository,
	trendRepo repository.TrendRepository,
) *DataCleaningService {
	return &DataCleaningService{
		masterRepo:     masterRepo,
		candidateRepo:  candidateRepo,
		trendRepo:      trendRepo,
		priceValidator: NewPriceValidator(),
		titleCleaner:   NewTitleCleaner(),
	}
}

// ProcessIncomingItem processes a new incoming item from DT platform
// Returns the promoted DTO if the item was matched to a master product
func (s *DataCleaningService) ProcessIncomingItem(
	ctx context.Context,
	item *entity.DTInputDTO,
	region string,
) (*entity.PlatformProductDTO, error) {
	if item == nil {
		return nil, fmt.Errorf("item cannot be nil")
	}

	rawTitle := item.Title
	cleanKey := s.titleCleaner.CleanTitleForID(rawTitle)

	// Try to match with existing master products
	masters, err := s.masterRepo.FindByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("find masters: %w", err)
	}

	// Strategy A: High confidence title match
	for _, master := range masters {
		if s.titleCleaner.IsHighSimilarity(rawTitle, master.StandardTitle) {
			return s.handleMasterMatch(ctx, master, item.Price, item.Status)
		}
	}

	// Strategy B: Mid confidence + price match (for typo correction)
	for _, master := range masters {
		if s.titleCleaner.IsMidSimilarity(rawTitle, master.StandardTitle) &&
			s.titleCleaner.IsPriceMatch(item.Price, master.Price) {
			return s.handleMasterMatch(ctx, master, item.Price, item.Status)
		}
	}

	// No match - add to candidate pool
	if err := s.handleCandidateLogic(ctx, region, rawTitle, cleanKey, item); err != nil {
		return nil, fmt.Errorf("handle candidate: %w", err)
	}

	return nil, nil
}

// handleMasterMatch handles when an item matches a master product
func (s *DataCleaningService) handleMasterMatch(
	ctx context.Context,
	master *entity.MasterProduct,
	price float64,
	status int,
) (*entity.PlatformProductDTO, error) {
	// Validate price update using Dutch auction model
	finalPrice, err := s.priceValidator.ValidateUpdate(
		master.Price,
		price,
		master.UpdateTime,
	)
	if err != nil {
		// Price anomaly detected - block update
		return nil, nil
	}

	// Update master
	master.Price = finalPrice
	master.Status = status
	master.IncrementTrustScore()

	if err := s.masterRepo.Update(ctx, master); err != nil {
		return nil, fmt.Errorf("update master: %w", err)
	}

	// Return DTO
	return &entity.PlatformProductDTO{
		ActivityID:         master.ID,
		Platform:           "DT",
		Region:             master.Region,
		Title:              master.StandardTitle,
		ShopName:           "DT生活精选",
		OriginalPrice:      price,
		CurrentPrice:       finalPrice,
		SalesStatus:        status,
		ActivityCreateTime: master.UpdateTime,
	}, nil
}

// handleCandidateLogic handles the candidate pool logic
func (s *DataCleaningService) handleCandidateLogic(
	ctx context.Context,
	region, rawTitle, cleanKey string,
	item *entity.DTInputDTO,
) error {
	// Validate input
	if rawTitle == "" {
		return fmt.Errorf("empty title")
	}
	if item.Price < 0 {
		return fmt.Errorf("invalid price: %f", item.Price)
	}

	// Check if candidate already exists
	candidates, err := s.candidateRepo.FindByRegion(ctx, region)
	if err != nil {
		return err
	}

	for _, candidate := range candidates {
		// Check for nil candidate
		if candidate == nil {
			continue
		}

		if s.titleCleaner.IsHighSimilarity(cleanKey, candidate.GroupKey) {
			// Update existing candidate
			candidate.AddTitleVote(rawTitle)
			candidate.UpdateLastSeen(item.Price, item.Status)

			return s.candidateRepo.Update(ctx, candidate)
		}
	}

	// Create new candidate
	candidate := &entity.CandidateItem{
		GroupKey:         cleanKey,
		Region:           region,
		TitleVotes:       map[string]int{rawTitle: 1},
		LastPrice:        item.Price,
		LastStatus:       item.Status,
		TotalOccurrences: 1,
		FirstSeenTime:    time.Now(),
		LastSeenTime:     time.Now(),
	}

	return s.candidateRepo.Create(ctx, candidate)
}

// PromoteCandidates promotes candidates that meet the threshold to master products
func (s *DataCleaningService) PromoteCandidates(
	ctx context.Context,
) (map[string][]*entity.PlatformProductDTO, error) {
	promotedData := make(map[string][]*entity.PlatformProductDTO)

	candidates, err := s.candidateRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list candidates: %w", err)
	}

	var toDeleteIDs []int64

	for _, candidate := range candidates {
		// Check for nil candidate
		if candidate == nil {
			continue
		}

		if !s.titleCleaner.ShouldPromote(candidate.TotalOccurrences) {
			continue
		}

		// Elect standard title
		winnerTitle := s.titleCleaner.ElectStandardTitle(candidate.TitleVotes)
		if winnerTitle == "" {
			log.Warn().
				Int64("id", candidate.ID).
				Msg("Skipping candidate with empty title")
			continue
		}

		uniqueID := s.titleCleaner.GenerateID("DT", winnerTitle)

		// Check if master already exists
		master, err := s.masterRepo.FindByID(ctx, uniqueID)
		if err != nil {
			return nil, fmt.Errorf("find master: %w", err)
		}

		if master == nil {
			// Create new master
			master = &entity.MasterProduct{
				ID:            uniqueID,
				Region:        candidate.Region,
				Platform:      "探探糖", // Default platform for DT data
				StandardTitle: winnerTitle,
				Price:         candidate.LastPrice,
				Status:        candidate.LastStatus,
				TrustScore:    candidate.TotalOccurrences,
			}

			if err := s.masterRepo.Create(ctx, master); err != nil {
				return nil, fmt.Errorf("create master: %w", err)
			}

			// Record initial price trend
			s.recordPriceTrend(ctx, master.ID, master.Price)
		} else {
			// Update existing master
			oldPrice := master.Price
			finalPrice, err := s.priceValidator.ValidateUpdate(
				master.Price,
				candidate.LastPrice,
				master.UpdateTime,
			)
			if err != nil {
				log.Error().Err(err).
					Str("masterId", master.ID).
					Msg("Failed to validate price during promotion")
				finalPrice = master.Price // Keep existing price on validation failure
			}
			master.Price = finalPrice
			master.Status = candidate.LastStatus

			if err := s.masterRepo.Update(ctx, master); err != nil {
				return nil, fmt.Errorf("update master: %w", err)
			}

			// Record price trend if price changed
			if finalPrice != oldPrice {
				s.recordPriceTrend(ctx, master.ID, finalPrice)
			}
		}

		// Add to promoted data
		dto := &entity.PlatformProductDTO{
			ActivityID:         master.ID,
			Platform:           "DT",
			Region:             master.Region,
			Title:              master.StandardTitle,
			ShopName:           "DT生活精选",
			OriginalPrice:      master.Price,
			CurrentPrice:       master.Price,
			SalesStatus:        master.Status,
			ActivityCreateTime: master.UpdateTime,
		}

		promotedData[candidate.Region] = append(promotedData[candidate.Region], dto)
		toDeleteIDs = append(toDeleteIDs, candidate.ID)
	}

	// Delete promoted candidates
	if len(toDeleteIDs) > 0 {
		if err := s.candidateRepo.DeleteByIDs(ctx, toDeleteIDs); err != nil {
			return nil, fmt.Errorf("delete candidates: %w", err)
		}
	}

	return promotedData, nil
}

// recordPriceTrend records a price trend for a master product
func (s *DataCleaningService) recordPriceTrend(ctx context.Context, activityID string, price float64) {
	if s.trendRepo == nil {
		return
	}

	// Truncate to day to ensure consistent date for ON CONFLICT clause
	trend, err := entity.NewPriceTrend(activityID, price, truncateToDay(time.Now()))
	if err != nil {
		log.Error().Err(err).
			Str("activityId", activityID).
			Float64("price", price).
			Msg("Failed to create price trend")
		return
	}

	if err := s.trendRepo.Upsert(ctx, trend); err != nil {
		log.Error().Err(err).
			Str("activityId", activityID).
			Float64("price", price).
			Msg("Failed to record price trend")
	}
}

// RecordDailyTrends records price trends for all master products
// This should be called daily to track price history
func (s *DataCleaningService) RecordDailyTrends(ctx context.Context) (int, error) {
	if s.trendRepo == nil {
		return 0, nil
	}

	// Get all master products
	masters, err := s.masterRepo.ListAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("list master products: %w", err)
	}

	count := 0
	// Truncate to day to ensure consistent date for ON CONFLICT clause
	now := truncateToDay(time.Now())

	for _, master := range masters {
		if master == nil {
			continue
		}

		trend, err := entity.NewPriceTrend(master.ID, master.Price, now)
		if err != nil {
			log.Error().Err(err).
				Str("masterId", master.ID).
				Float64("price", master.Price).
				Msg("Failed to create price trend entity")
			continue
		}

		if err := s.trendRepo.Upsert(ctx, trend); err != nil {
			log.Error().Err(err).
				Str("masterId", master.ID).
				Float64("price", master.Price).
				Msg("Failed to record price trend")
			continue
		}

		count++
	}

	return count, nil
}
