package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"

	"github.com/rs/zerolog/log"
)

// NotificationService handles price notifications
type NotificationService struct {
	notiRepo       repository.NotificationRepository
	prodRepo       repository.ProductRepository
	userSettingsRepo repository.UserSettingsRepository
	barkURL        string
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notiRepo repository.NotificationRepository,
	prodRepo repository.ProductRepository,
	userSettingsRepo repository.UserSettingsRepository,
	barkURL string,
) *NotificationService {
	return &NotificationService{
		notiRepo:       notiRepo,
		prodRepo:       prodRepo,
		userSettingsRepo: userSettingsRepo,
		barkURL:        barkURL,
	}
}

// Create creates a new notification configuration
func (s *NotificationService) Create(ctx context.Context, userID, activityID string, targetPrice float64) error {
	config := &entity.NotificationConfig{
		ActivityID:  activityID,
		UserID:      userID,
		TargetPrice: targetPrice,
	}

	return s.notiRepo.Upsert(ctx, config)
}

// Update updates an existing notification configuration
func (s *NotificationService) Update(ctx context.Context, userID, activityID string, targetPrice float64) error {
	config, err := s.notiRepo.FindByActivityID(ctx, activityID, userID)
	if err != nil {
		return err
	}

	config.TargetPrice = targetPrice
	return s.notiRepo.Upsert(ctx, config)
}

// Upsert creates or updates a notification configuration
func (s *NotificationService) Upsert(ctx context.Context, userID, activityID string, targetPrice float64) error {
	config := &entity.NotificationConfig{
		ActivityID:  activityID,
		UserID:      userID,
		TargetPrice: targetPrice,
	}

	return s.notiRepo.Upsert(ctx, config)
}

// Delete deletes a notification configuration
func (s *NotificationService) Delete(ctx context.Context, userID, activityID string) error {
	return s.notiRepo.Delete(ctx, activityID, userID)
}

// CheckAndNotify checks all notifications and sends notifications for matching products
func (s *NotificationService) CheckAndNotify(ctx context.Context) error {
	configs, err := s.notiRepo.ListAll(ctx)
	if err != nil {
		return err
	}

	for _, config := range configs {
		if s.checkAndNotifySingle(ctx, config) {
			// Update last notify time
			if err := s.notiRepo.UpdateNotifyTime(ctx, config.ActivityID, config.UserID); err != nil {
				log.Error().Err(err).
					Str("activityId", config.ActivityID).
					Str("userId", config.UserID).
					Msg("failed to update notification time")
			}
		}
	}

	return nil
}

// checkAndNotifySingle checks a single notification and sends if conditions are met
func (s *NotificationService) checkAndNotifySingle(ctx context.Context, config *entity.NotificationConfig) bool {
	// Check if already notified today
	if config.HasNotifiedToday() {
		return false
	}

	// Get user's Bark Key
	userSettings, err := s.userSettingsRepo.Get(ctx, config.UserID)
	if err != nil || userSettings == nil {
		log.Error().Err(err).
			Str("userId", config.UserID).
			Msg("failed to get user settings")
		return false
	}

	// Get current product price
	product, err := s.prodRepo.FindByActivityID(ctx, config.ActivityID)
	if err != nil {
		log.Error().Err(err).
			Str("activityId", config.ActivityID).
			Msg("failed to find product")
		return false
	}

	// Check if product is nil
	if product == nil {
		log.Warn().
			Str("activityId", config.ActivityID).
			Msg("Product not found for notification")
		return false
	}

	// Check if product is on sale
	if !product.IsOnSale() {
		return false
	}

	// Check if price condition is met
	if !config.ShouldNotify(product.CurrentPrice) {
		return false
	}

	// Send notification using user's Bark Key
	return s.sendNotification(product, config.TargetPrice, userSettings.BarkKey)
}

// sendNotification sends a push notification
func (s *NotificationService) sendNotification(product *entity.Product, targetPrice float64, barkKey string) bool {
	message := fmt.Sprintf("【%s %s ¥%.2f】%s",
		product.Platform,
		product.Region,
		product.CurrentPrice,
		product.Title,
	)

	log.Info().
		Str("activityId", product.ActivityID).
		Str("message", message).
		Msg("Sending price notification")

	if barkKey == "" {
		log.Info().Msg("Bark key not configured, skipping actual send")
		return true
	}

	// Build Bark URL with user's bark key
	barkBaseURL := s.barkURL
	if barkBaseURL == "" {
		barkBaseURL = "https://api.day.app"
	}

	// Normalize barkKey - extract device key from URL if needed
	deviceKey := normalizeBarkKey(barkKey)

	// Send to Bark using HTTP client
	encodedMsg := url.QueryEscape(message)
	fullURL := fmt.Sprintf("%s/%s/%s?level=critical&volume=5", barkBaseURL, deviceKey, encodedMsg)

	// Create HTTP request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fullURL, nil)
	if err != nil {
		log.Error().Err(err).
			Str("url", fullURL).
			Msg("Failed to create bark request")
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).
			Str("url", fullURL).
			Msg("Failed to send bark notification")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error().
			Int("status", resp.StatusCode).
			Str("url", fullURL).
			Msg("Bark notification failed")
		return false
	}

	log.Info().
		Str("url", fullURL).
		Msg("Bark notification sent successfully")

	return true
}

// normalizeBarkKey extracts device key from full URL or returns key as-is
func normalizeBarkKey(input string) string {
	input = strings.TrimSpace(input)
	// If full URL, extract last segment as device key
	if strings.HasPrefix(input, "http") {
		parts := strings.Split(input, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return input
}

// GetByActivityID returns notification config for a product
func (s *NotificationService) GetByActivityID(ctx context.Context, userID, activityID string) (*entity.NotificationConfig, error) {
	return s.notiRepo.FindByActivityID(ctx, activityID, userID)
}
