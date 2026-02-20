package scheduler

import (
	"context"
	"fmt"
	"time"

	"kbfood/internal/config"
	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	"kbfood/internal/domain/service"
	"kbfood/internal/infra/platform"

	"github.com/rs/zerolog/log"
)

// SyncJob synchronizes products from TanTanTang platform
type SyncJob struct {
	cfg             *config.Config
	client          *platform.TanTanTangClient
	cleaningService *service.DataCleaningService
	syncStatusRepo  repository.SyncStatusRepository
}

// NewSyncJob creates a new sync job
func NewSyncJob(
	cfg *config.Config,
	client *platform.TanTanTangClient,
	cleaningService *service.DataCleaningService,
	syncStatusRepo repository.SyncStatusRepository,
) *SyncJob {
	return &SyncJob{
		cfg:             cfg,
		client:          client,
		cleaningService: cleaningService,
		syncStatusRepo:  syncStatusRepo,
	}
}

// Name returns the job name
func (j *SyncJob) Name() string {
	return "sync-tantantang"
}

// Run executes the sync job
func (j *SyncJob) Run(ctx context.Context) error {
	startTime := time.Now()

	if j.client == nil {
		j.recordStatus(ctx, startTime, 0, fmt.Errorf("client not initialized"))
		return fmt.Errorf("client not initialized")
	}
	if j.cleaningService == nil {
		j.recordStatus(ctx, startTime, 0, fmt.Errorf("cleaningService not initialized"))
		return fmt.Errorf("cleaningService not initialized")
	}

	regions := []string{"广州", "佛山"}

	totalProducts := 0
	var lastErr error

	for _, region := range regions {
		// Check for context cancellation before each region
		if ctx.Err() != nil {
			lastErr = ctx.Err()
			break
		}

		// TanTanTangClient.FetchProducts handles pagination internally
		products, err := j.client.FetchProducts(ctx, region)
		if err != nil {
			log.Error().Err(err).
				Str("region", region).
				Msg("Failed to fetch products")
			lastErr = err
			continue
		}

		if len(products) == 0 {
			log.Warn().
				Str("region", region).
				Msg("No products fetched from API")
			continue
		}

		if products == nil {
			log.Warn().
				Str("region", region).
				Msg("No products fetched")
			continue
		}

		for _, p := range products {
			// Nil check for individual products
			if p == nil {
				log.Warn().Msg("Skipping nil product")
				continue
			}

			// Convert PlatformProductDTO to DTInputDTO for processing
			input := &entity.DTInputDTO{
				Title:     p.Title,
				Price:     p.CurrentPrice,
				Status:    p.SalesStatus,
				CrawlTime: p.ActivityCreateTime.Unix(),
				Region:    region,
			}

			_, err := j.cleaningService.ProcessIncomingItem(ctx, input, region)
			if err != nil {
				log.Error().Err(err).
					Str("title", p.Title).
					Str("region", region).
					Msg("Failed to process product")
			} else {
				totalProducts++
			}
		}
	}

	log.Info().
		Int("totalProducts", totalProducts).
		Msg("Sync job completed")

	// Record sync status
	j.recordStatus(ctx, startTime, totalProducts, lastErr)

	return lastErr
}

// recordStatus records the sync status to the database
func (j *SyncJob) recordStatus(ctx context.Context, startTime time.Time, productCount int, err error) {
	if j.syncStatusRepo == nil {
		return
	}

	status := &entity.SyncStatus{
		JobName:      j.Name(),
		LastRunTime:  startTime,
		ProductCount: productCount,
		Status:       entity.StatusSuccess,
	}

	if err != nil {
		status.Status = entity.StatusFailed
		status.ErrorMessage = err.Error()
	}

	if recordErr := j.syncStatusRepo.Upsert(ctx, status); recordErr != nil {
		log.Error().Err(recordErr).Msg("Failed to record sync status")
	}
}
