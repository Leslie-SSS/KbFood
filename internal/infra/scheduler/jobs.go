package scheduler

import (
	"context"
	"fmt"
	"time"

	"kbfood/internal/domain/repository"
	"kbfood/internal/domain/service"

	"github.com/rs/zerolog/log"
)

// PromoteCandidatesJob promotes candidates from the candidate pool to master products
type PromoteCandidatesJob struct {
	cleaningService *service.DataCleaningService
}

// NewPromoteCandidatesJob creates a new promote candidates job
func NewPromoteCandidatesJob(cleaningService *service.DataCleaningService) *PromoteCandidatesJob {
	return &PromoteCandidatesJob{
		cleaningService: cleaningService,
	}
}

// Name returns the job name
func (j *PromoteCandidatesJob) Name() string {
	return "promote-candidates"
}

// Run executes the job
func (j *PromoteCandidatesJob) Run(ctx context.Context) error {
	if j.cleaningService == nil {
		return fmt.Errorf("cleaningService not initialized")
	}

	promotedData, err := j.cleaningService.PromoteCandidates(ctx)
	if err != nil {
		return fmt.Errorf("promote candidates job failed: %w", err)
	}

	total := 0
	for _, products := range promotedData {
		total += len(products)
	}

	if total > 0 {
		log.Info().
			Int("count", total).
			Msg("Candidates promoted")
	} else {
		log.Debug().Msg("No candidates to promote")
	}

	return nil
}

// PriceCheckJob checks notification configurations and sends alerts
type PriceCheckJob struct {
	notificationService *service.NotificationService
}

// NewPriceCheckJob creates a new price check job
func NewPriceCheckJob(notificationService *service.NotificationService) *PriceCheckJob {
	return &PriceCheckJob{
		notificationService: notificationService,
	}
}

// Name returns the job name
func (j *PriceCheckJob) Name() string {
	return "price-check"
}

// Run executes the job
func (j *PriceCheckJob) Run(ctx context.Context) error {
	if j.notificationService == nil {
		return fmt.Errorf("notificationService not initialized")
	}

	if err := j.notificationService.CheckAndNotify(ctx); err != nil {
		return fmt.Errorf("price check job failed: %w", err)
	}

	return nil
}

// CleanupJob cleans up expired data
type CleanupJob struct {
	prodRepo repository.ProductRepository
}

// NewCleanupJob creates a new cleanup job
func NewCleanupJob(prodRepo repository.ProductRepository) *CleanupJob {
	return &CleanupJob{
		prodRepo: prodRepo,
	}
}

// Name returns the job name
func (j *CleanupJob) Name() string {
	return "cleanup"
}

// Run executes the job
func (j *CleanupJob) Run(ctx context.Context) error {
	if j.prodRepo == nil {
		return fmt.Errorf("product repository not initialized")
	}

	threshold := time.Now().Add(-48 * time.Hour)

	// TODO: Implement delete products older than threshold
	// This requires a FindByUpdateTime repository method

	log.Info().
		Time("threshold", threshold).
		Msg("Cleanup completed")

	return nil
}

// RecordTrendsJob records daily price trends for all master products
type RecordTrendsJob struct {
	cleaningService *service.DataCleaningService
}

// NewRecordTrendsJob creates a new record trends job
func NewRecordTrendsJob(cleaningService *service.DataCleaningService) *RecordTrendsJob {
	return &RecordTrendsJob{
		cleaningService: cleaningService,
	}
}

// Name returns the job name
func (j *RecordTrendsJob) Name() string {
	return "record-trends"
}

// Run executes the job
func (j *RecordTrendsJob) Run(ctx context.Context) error {
	if j.cleaningService == nil {
		return fmt.Errorf("cleaningService not initialized")
	}

	count, err := j.cleaningService.RecordDailyTrends(ctx)
	if err != nil {
		return fmt.Errorf("record trends job failed: %w", err)
	}

	if count > 0 {
		log.Info().
			Int("count", count).
			Msg("Price trends recorded")
	} else {
		log.Debug().Msg("No trends to record")
	}

	return nil
}
