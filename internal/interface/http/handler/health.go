package handler

import (
	"context"
	"net/http"
	"time"

	"kbfood/internal/infra/db"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *db.Pool
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(database *db.Pool) *HealthHandler {
	return &HealthHandler{db: database}
}

// Health handles GET /health - basic liveness check
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

// Ready handles GET /ready - readiness check with database connectivity
func (h *HealthHandler) Ready(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	status := map[string]interface{}{
		"status": "ready",
	}

	// Check database connectivity if pool is available
	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			status["status"] = "not_ready"
			status["database"] = "unreachable"
			status["error"] = err.Error()
			return c.JSON(http.StatusServiceUnavailable, status)
		}
		status["database"] = "ok"
	}

	return c.JSON(http.StatusOK, status)
}
