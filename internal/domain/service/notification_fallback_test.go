package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
)

type stubNotificationRepository struct {
	configs           []*entity.NotificationConfig
	updatedActivityID string
	updatedUserID     string
}

func (s *stubNotificationRepository) FindByActivityID(ctx context.Context, activityID string, userID string) (*entity.NotificationConfig, error) {
	for _, config := range s.configs {
		if config.ActivityID == activityID && config.UserID == userID {
			return config, nil
		}
	}
	return nil, nil
}

func (s *stubNotificationRepository) ListByUser(ctx context.Context, userID string) ([]*entity.NotificationConfig, error) {
	var result []*entity.NotificationConfig
	for _, config := range s.configs {
		if config.UserID == userID {
			result = append(result, config)
		}
	}
	return result, nil
}

func (s *stubNotificationRepository) ListAll(ctx context.Context) ([]*entity.NotificationConfig, error) {
	return s.configs, nil
}

func (s *stubNotificationRepository) Upsert(ctx context.Context, config *entity.NotificationConfig) error {
	return nil
}

func (s *stubNotificationRepository) Delete(ctx context.Context, activityID string, userID string) error {
	return nil
}

func (s *stubNotificationRepository) UpdateNotifyTime(ctx context.Context, activityID string, userID string) error {
	s.updatedActivityID = activityID
	s.updatedUserID = userID
	return nil
}

type stubProductRepository struct{}

func (s *stubProductRepository) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	return nil, nil
}

func (s *stubProductRepository) FindByActivityID(ctx context.Context, activityID string) (*entity.Product, error) {
	return nil, nil
}

func (s *stubProductRepository) FindByFilter(ctx context.Context, filter repository.ProductFilter) ([]*entity.Product, error) {
	return nil, nil
}

func (s *stubProductRepository) Create(ctx context.Context, product *entity.Product) error {
	return nil
}

func (s *stubProductRepository) Update(ctx context.Context, product *entity.Product) error {
	return nil
}

func (s *stubProductRepository) UpdateByActivityID(ctx context.Context, activityID string, product *entity.Product) error {
	return nil
}

func (s *stubProductRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *stubProductRepository) DeleteByActivityIDs(ctx context.Context, activityIDs []string) error {
	return nil
}

func (s *stubProductRepository) DeleteByPlatform(ctx context.Context, platform string) error {
	return nil
}

func (s *stubProductRepository) CountByPlatform(ctx context.Context, platform string) (int64, error) {
	return 0, nil
}

func (s *stubProductRepository) ListBlocked(ctx context.Context) ([]*entity.Product, error) {
	return nil, nil
}

type stubMasterProductRepository struct {
	product *entity.MasterProduct
}

func (s *stubMasterProductRepository) FindByID(ctx context.Context, id string) (*entity.MasterProduct, error) {
	if s.product != nil && s.product.ID == id {
		return s.product, nil
	}
	return nil, nil
}

func (s *stubMasterProductRepository) FindByRegion(ctx context.Context, region string) ([]*entity.MasterProduct, error) {
	return nil, nil
}

func (s *stubMasterProductRepository) FindByPlatform(ctx context.Context, platform string) ([]*entity.MasterProduct, error) {
	return nil, nil
}

func (s *stubMasterProductRepository) FindByRegionAndPlatform(ctx context.Context, region, platform string) ([]*entity.MasterProduct, error) {
	return nil, nil
}

func (s *stubMasterProductRepository) ListAll(ctx context.Context) ([]*entity.MasterProduct, error) {
	return nil, nil
}

func (s *stubMasterProductRepository) Create(ctx context.Context, product *entity.MasterProduct) error {
	return nil
}

func (s *stubMasterProductRepository) Update(ctx context.Context, product *entity.MasterProduct) error {
	return nil
}

func (s *stubMasterProductRepository) Delete(ctx context.Context, id string) error {
	return nil
}

type stubUserSettingsRepository struct {
	settings *entity.UserSettings
}

func (s *stubUserSettingsRepository) Get(ctx context.Context, userID string) (*entity.UserSettings, error) {
	if s.settings != nil && s.settings.UserID == userID {
		return s.settings, nil
	}
	return nil, nil
}

func (s *stubUserSettingsRepository) Upsert(ctx context.Context, settings *entity.UserSettings) error {
	s.settings = settings
	return nil
}

func TestNotificationService_CheckAndNotifyFallsBackToMasterProduct(t *testing.T) {
	ctx := context.Background()

	notiRepo := &stubNotificationRepository{
		configs: []*entity.NotificationConfig{
			{
				ActivityID:  "DT_ed87d13fb7fcf11aea6de4129038b765",
				UserID:      "client-123",
				TargetPrice: 75.05,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			},
		},
	}
	prodRepo := &stubProductRepository{}
	masterRepo := &stubMasterProductRepository{
		product: &entity.MasterProduct{
			ID:            "DT_ed87d13fb7fcf11aea6de4129038b765",
			Region:        "广州",
			Platform:      "DT",
			StandardTitle: "火锅四人餐",
			Price:         68.7,
			Status:        entity.SalesStatusOnSale,
		},
	}
	userSettingsRepo := &stubUserSettingsRepository{
		settings: &entity.UserSettings{
			UserID:  "client-123",
			BarkKey: "https://api.day.app/DEVICE123/?isArchive=1",
		},
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if !strings.Contains(r.URL.Path, "/DEVICE123/") {
			t.Fatalf("expected normalized Bark key in request path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := NewNotificationService(notiRepo, prodRepo, masterRepo, userSettingsRepo, server.URL)

	if err := service.CheckAndNotify(ctx); err != nil {
		t.Fatalf("CheckAndNotify() error = %v", err)
	}

	if requestCount != 1 {
		t.Fatalf("expected 1 Bark request, got %d", requestCount)
	}
	if notiRepo.updatedActivityID != "DT_ed87d13fb7fcf11aea6de4129038b765" || notiRepo.updatedUserID != "client-123" {
		t.Fatalf("expected notify time update for migrated notification, got activity=%q user=%q", notiRepo.updatedActivityID, notiRepo.updatedUserID)
	}
}
