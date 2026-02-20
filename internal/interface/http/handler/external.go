package handler

import (
	"encoding/json"
	"net/http"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/service"
	"kbfood/internal/interface/http/dto"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ExternalHandler handles external platform requests
type ExternalHandler struct {
	cleaningService *service.DataCleaningService
}

// NewExternalHandler creates a new external handler
func NewExternalHandler(
	cleaningService *service.DataCleaningService,
) *ExternalHandler {
	return &ExternalHandler{
		cleaningService: cleaningService,
	}
}

// DTPlatformPushRequest represents DT platform push request
type DTPlatformPushRequest struct {
	Items []DTPlatformItem `json:"items"`
}

// DTPlatformItem represents a single DT platform item
type DTPlatformItem struct {
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Status    int     `json:"status"`
	CrawlTime int64   `json:"crawlTime"`
	Region    string  `json:"region"`
}

// HandleDTPush handles POST /api/external/dt/push
func (h *ExternalHandler) HandleDTPush(c echo.Context) error {
	ctx := c.Request().Context()

	var req DTPlatformPushRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode DT push request")
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request body"))
	}

	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "items array is required and must not be empty"))
	}

	if len(req.Items) > 1000 {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "too many items in request (max 1000)"))
	}

	promotedCount := 0
	for _, item := range req.Items {
		input := &entity.DTInputDTO{
			Title:     item.Title,
			Price:     item.Price,
			Status:    item.Status,
			CrawlTime: item.CrawlTime,
			Region:    item.Region,
		}

		promoted, err := h.cleaningService.ProcessIncomingItem(ctx, input, item.Region)
		if err != nil {
			log.Error().Err(err).
				Str("title", item.Title).
				Str("region", item.Region).
				Msg("Failed to process DT item")
			continue
		}

		if promoted != nil {
			promotedCount++
		}
	}

	log.Info().
		Int("total", len(req.Items)).
		Int("promoted", promotedCount).
		Msg("DT push processed")

	return c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
		"received": len(req.Items),
		"promoted": promotedCount,
	}))
}
