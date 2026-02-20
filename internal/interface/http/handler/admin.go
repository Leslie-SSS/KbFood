package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kbfood/internal/infra/platform"
	"kbfood/internal/interface/http/dto"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SyncHandler handles manual sync operations
type SyncHandler struct {
	syncJob   SyncJobRunner
	tttClient *platform.TanTanTangClient
}

// SyncJobRunner interface for jobs that can be manually triggered
type SyncJobRunner interface {
	Run(ctx context.Context) error
	Name() string
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(syncJob SyncJobRunner, tttClient *platform.TanTanTangClient) *SyncHandler {
	return &SyncHandler{
		syncJob:   syncJob,
		tttClient: tttClient,
	}
}

// TriggerSync handles POST /api/admin/sync
func (h *SyncHandler) TriggerSync(c echo.Context) error {
	ctx := c.Request().Context()

	log.Info().Msg("Manual sync triggered")

	if err := h.syncJob.Run(ctx); err != nil {
		log.Error().Err(err).Msg("Manual sync failed")
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "sync operation failed"))
	}

	return c.JSON(http.StatusOK, dto.Success(map[string]string{
		"message": "Sync completed successfully",
		"job":     h.syncJob.Name(),
	}))
}

// TestAPI handles GET /api/admin/test-api - test TanTanTang API directly
func (h *SyncHandler) TestAPI(c echo.Context) error {
	ctx := c.Request().Context()
	region := c.QueryParam("region")
	if region == "" {
		region = "广州"
	}

	// Get region config for coords
	cityCfg := platform.GetRegionConfig(region)

	// Test API call
	products, err := h.tttClient.FetchProducts(ctx, region)
	if err != nil {
		log.Error().Err(err).Str("region", region).Msg("API test failed")
		return c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
			"region":      region,
			"success":     false,
			"error":       err.Error(),
			"count":       0,
			"products":    []interface{}{},
			"serverTime":  time.Now().Format("2006-01-02 15:04:05"),
		}))
	}

	return c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
		"region":      region,
		"city":        cityCfg.CityName,
		"success":     true,
		"count":       len(products),
		"sampleCount": min(len(products), 3),
		"products":    products,
		"serverTime":  time.Now().Format("2006-01-02 15:04:05"),
	}))
}

// TestNotification handles POST /api/admin/test-notification
// Sends a test notification via Bark to verify the user's Bark key is valid
func (h *SyncHandler) TestNotification(c echo.Context) error {
	var params struct {
		BarkKey string `json:"barkKey"`
	}

	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request format"))
	}

	if params.BarkKey == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Bark key is required"))
	}

	// Build Bark URL
	// User can provide either:
	// 1. Full URL: https://api.day.app/XXXXXX
	// 2. Device key only: XXXXXX
	barkURL := strings.TrimSpace(params.BarkKey)
	if !strings.HasPrefix(barkURL, "http") {
		barkURL = fmt.Sprintf("https://api.day.app/%s", barkURL)
	}

	// Build notification URL with title and body
	title := url.QueryEscape("测试通知")
	body := url.QueryEscape("美食监控配置成功！您可以收到价格提醒了。")
	fullURL := fmt.Sprintf("%s/%s/%s", barkURL, title, body)

	log.Info().Str("barkURL", barkURL).Msg("Sending test notification")

	// Send test notification
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		log.Error().Err(err).Msg("Test notification request failed")
		return c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
			"success":  false,
			"error":    err.Error(),
			"testedAt": time.Now().Format("2006-01-02 15:04:05"),
		}))
	}
	defer resp.Body.Close()

	success := resp.StatusCode == 200
	log.Info().Bool("success", success).Int("statusCode", resp.StatusCode).Msg("Test notification completed")

	return c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
		"success":    success,
		"statusCode": resp.StatusCode,
		"testedAt":   time.Now().Format("2006-01-02 15:04:05"),
	}))
}
