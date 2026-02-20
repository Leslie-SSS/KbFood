package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	"kbfood/internal/interface/http/dto"
	"kbfood/internal/interface/http/middleware"

	"github.com/labstack/echo/v4"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	prodRepo    repository.ProductRepository
	masterRepo  repository.MasterProductRepository
	notiRepo    repository.NotificationRepository
	blockedRepo repository.BlockedRepository
	trendRepo   repository.TrendRepository
}

// NewProductHandler creates a new product handler
func NewProductHandler(
	prodRepo repository.ProductRepository,
	masterRepo repository.MasterProductRepository,
	notiRepo repository.NotificationRepository,
	blockedRepo repository.BlockedRepository,
	trendRepo repository.TrendRepository,
) *ProductHandler {
	return &ProductHandler{
		prodRepo:    prodRepo,
		masterRepo:  masterRepo,
		notiRepo:    notiRepo,
		blockedRepo: blockedRepo,
		trendRepo:   trendRepo,
	}
}

// QueryProducts handles GET /api/products
func (h *ProductHandler) QueryProducts(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	// Parse query parameters
	region := c.QueryParam("region")
	platform := c.QueryParam("platform")
	keyword := c.QueryParam("keyword")
	salesStatusStr := c.QueryParam("salesStatus")
	monitorStatus := c.QueryParam("monitorStatus")

	var salesStatus *int
	if salesStatusStr != "" {
		val := 0
		if salesStatusStr == "1" {
			val = 1
		}
		salesStatus = &val
	}

	// Get blocked products list for user
	var blockedSet map[string]bool
	if userID != "" {
		blockedIDs, _ := h.blockedRepo.List(ctx, userID)
		blockedSet = make(map[string]bool)
		for _, id := range blockedIDs {
			blockedSet[id] = true
		}
	} else {
		blockedSet = make(map[string]bool)
	}

	// Get notification configs for user
	var notificationMap map[string]*entity.NotificationConfig
	if userID != "" {
		notis, _ := h.notiRepo.ListByUser(ctx, userID)
		notificationMap = make(map[string]*entity.NotificationConfig)
		for _, n := range notis {
			notificationMap[n.ActivityID] = n
		}
	} else {
		notificationMap = make(map[string]*entity.NotificationConfig)
	}

	// Fetch master products with platform and region filters
	var masterProducts []*entity.MasterProduct
	var err error

	if region != "" && platform != "" {
		masterProducts, err = h.masterRepo.FindByRegionAndPlatform(ctx, region, platform)
	} else if region != "" {
		masterProducts, err = h.masterRepo.FindByRegion(ctx, region)
	} else if platform != "" {
		masterProducts, err = h.masterRepo.FindByPlatform(ctx, platform)
	} else {
		masterProducts, err = h.masterRepo.ListAll(ctx)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to fetch products"))
	}

	// Apply filters
	filtered := make([]*entity.MasterProduct, 0)
	for _, p := range masterProducts {
		// Skip blocked products
		if blockedSet[p.ID] {
			continue
		}

		// Filter by keyword
		if keyword != "" && !containsIgnoreCase(p.StandardTitle, keyword) {
			continue
		}

		// Filter by sales status
		if salesStatus != nil && p.Status != *salesStatus {
			continue
		}

		// Filter by monitor status
		_, hasNotification := notificationMap[p.ID]
		if monitorStatus == "1" && !hasNotification {
			continue
		}
		if monitorStatus == "0" && hasNotification {
			continue
		}

		filtered = append(filtered, p)
	}

	// Convert to DTOs with notification info
	result := make([]dto.ProductDTO, 0, len(filtered))
	for _, p := range filtered {
		productDTO := dto.FromMasterEntity(p)
		if noti, exists := notificationMap[p.ID]; exists {
			productDTO.HasNotification = true
			productDTO.TargetPrice = &noti.TargetPrice
		}
		result = append(result, productDTO)
	}

	return c.JSON(http.StatusOK, dto.Success(result))
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
// Uses strings.ToLower which properly handles UTF-8 multi-byte characters
func containsIgnoreCase(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// GetPriceTrend handles GET /api/products/:activityId/trend
func (h *ProductHandler) GetPriceTrend(c echo.Context) error {
	ctx := c.Request().Context()
	activityID := c.Param("activityId")

	if activityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}

	trends, err := h.trendRepo.FindByActivityID(ctx, activityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to fetch price trends"))
	}

	result := dto.FromTrendEntities(trends)
	return c.JSON(http.StatusOK, dto.Success(result))
}

// BlockProduct handles POST /api/products/:activityId/block
func (h *ProductHandler) BlockProduct(c echo.Context) error {
	ctx := c.Request().Context()
	activityID := c.Param("activityId")
	userID := middleware.GetUserID(c)

	if activityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请先设置 Bark Key"))
	}

	if err := h.blockedRepo.Create(ctx, activityID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to block product"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}

// UnblockProduct handles POST /api/products/unblock/:activityId
func (h *ProductHandler) UnblockProduct(c echo.Context) error {
	ctx := c.Request().Context()
	activityID := c.Param("activityId")
	userID := middleware.GetUserID(c)

	if activityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请先设置 Bark Key"))
	}

	if err := h.blockedRepo.Delete(ctx, activityID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to unblock product"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}

// GetBlockedProducts handles GET /api/products/blocked
func (h *ProductHandler) GetBlockedProducts(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	if userID == "" {
		return c.JSON(http.StatusOK, dto.Success([]dto.ProductDTO{}))
	}

	// Get blocked activity IDs for user
	blockedIDs, err := h.blockedRepo.List(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to fetch blocked products"))
	}

	if len(blockedIDs) == 0 {
		return c.JSON(http.StatusOK, dto.Success([]dto.ProductDTO{}))
	}

	// Get master products for blocked IDs
	result := make([]dto.ProductDTO, 0)
	for _, activityID := range blockedIDs {
		master, err := h.masterRepo.FindByID(ctx, activityID)
		if err == nil && master != nil {
			result = append(result, dto.FromMasterEntity(master))
		}
	}

	return c.JSON(http.StatusOK, dto.Success(result))
}

// ClearPlatform handles DELETE /api/products/platform/:platform
func (h *ProductHandler) ClearPlatform(c echo.Context) error {
	ctx := c.Request().Context()
	platform := c.Param("platform")

	if err := h.prodRepo.DeleteByPlatform(ctx, platform); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to clear platform data"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}

// CreateNotification handles POST /api/products/notifications
func (h *ProductHandler) CreateNotification(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请先设置 Bark Key"))
	}

	var params struct {
		ActivityID  string  `json:"activityId"`
		TargetPrice float64 `json:"targetPrice"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request body"))
	}

	if params.ActivityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}
	if params.TargetPrice <= 0 {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "targetPrice must be positive"))
	}

	config := &entity.NotificationConfig{
		ActivityID:  params.ActivityID,
		UserID:      userID,
		TargetPrice: params.TargetPrice,
	}

	if err := h.notiRepo.Upsert(ctx, config); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to create notification"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}

// UpdateNotification handles PUT /api/products/notifications/:activityId
func (h *ProductHandler) UpdateNotification(c echo.Context) error {
	ctx := c.Request().Context()
	activityID := c.Param("activityId")
	userID := middleware.GetUserID(c)

	if activityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请先设置 Bark Key"))
	}

	var params struct {
		TargetPrice float64 `json:"targetPrice"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request body"))
	}

	if params.TargetPrice <= 0 {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "targetPrice must be positive"))
	}

	config, err := h.notiRepo.FindByActivityID(ctx, activityID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.Error(404, "Notification not found"))
	}
	if config == nil {
		return c.JSON(http.StatusNotFound, dto.Error(404, "Notification not found"))
	}

	config.TargetPrice = params.TargetPrice
	if err := h.notiRepo.Upsert(ctx, config); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update notification"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}

// DeleteNotification handles DELETE /api/products/notifications/:activityId
func (h *ProductHandler) DeleteNotification(c echo.Context) error {
	ctx := c.Request().Context()
	activityID := c.Param("activityId")
	userID := middleware.GetUserID(c)

	if activityID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "activityId is required"))
	}
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请先设置 Bark Key"))
	}

	if err := h.notiRepo.Delete(ctx, activityID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to delete notification"))
	}

	return c.JSON(http.StatusOK, dto.Success(nil))
}
