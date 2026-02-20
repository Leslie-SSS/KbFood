package handler

import (
	"net/http"
	"time"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	"kbfood/internal/interface/http/dto"

	"github.com/labstack/echo/v4"
)

// StatusHandler handles system status requests
type StatusHandler struct {
	syncStatusRepo repository.SyncStatusRepository
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(syncStatusRepo repository.SyncStatusRepository) *StatusHandler {
	return &StatusHandler{
		syncStatusRepo: syncStatusRepo,
	}
}

// SystemStatusResponse represents the system status response
type SystemStatusResponse struct {
	Sync       SyncStatus   `json:"sync"`
	ServerTime string       `json:"serverTime"`
}

// SyncStatus represents the sync job status
type SyncStatus struct {
	LastRunTime  string `json:"lastRunTime"`
	Status       string `json:"status"`
	ProductCount int    `json:"productCount"`
	IsHealthy    bool   `json:"isHealthy"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// GetStatus handles GET /api/status - returns system status
func (h *StatusHandler) GetStatus(c echo.Context) error {
	ctx := c.Request().Context()

	// Get latest sync status
	syncStatus, err := h.syncStatusRepo.GetLatest(ctx, "sync-tantantang")
	if err != nil {
		// Return default status if query fails
		return c.JSON(http.StatusOK, dto.Success(SystemStatusResponse{
			Sync: SyncStatus{
				LastRunTime:  "",
				Status:       entity.StatusFailed,
				ProductCount: 0,
				IsHealthy:    false,
				ErrorMessage: "Failed to retrieve sync status",
			},
			ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		}))
	}

	// Handle case where no sync has run yet
	if syncStatus == nil {
		return c.JSON(http.StatusOK, dto.Success(SystemStatusResponse{
			Sync: SyncStatus{
				LastRunTime:  "",
				Status:       "pending",
				ProductCount: 0,
				IsHealthy:    false,
				ErrorMessage: "No sync has been executed yet",
			},
			ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		}))
	}

	return c.JSON(http.StatusOK, dto.Success(SystemStatusResponse{
		Sync: SyncStatus{
			LastRunTime:  syncStatus.LastRunTime.Format("2006-01-02 15:04:05"),
			Status:       syncStatus.Status,
			ProductCount: syncStatus.ProductCount,
			IsHealthy:    syncStatus.IsHealthy(),
			ErrorMessage: syncStatus.ErrorMessage,
		},
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
	}))
}
