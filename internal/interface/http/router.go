package http

import (
	"strings"

	"kbfood/internal/infra/db"
	"kbfood/internal/interface/http/handler"
	"kbfood/internal/interface/http/middleware"

	"github.com/labstack/echo/v4"
)

// Router returns the HTTP router
func Router(
	productHandler *handler.ProductHandler,
	externalHandler *handler.ExternalHandler,
	syncHandler *handler.SyncHandler,
	statusHandler *handler.StatusHandler,
	userHandler *handler.UserHandler,
	database *db.Pool,
) *echo.Echo {
	e := echo.New()

	// Create health handler with database
	healthHandler := handler.NewHealthHandler(database)

	// Hide Echo banner
	e.HideBanner = true

	// Global middleware
	e.Use(middleware.Recovery())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.UserExtractor())

	// Health check (no auth required)
	e.GET("/health", healthHandler.Health)
	e.GET("/ready", healthHandler.Ready)

	// API routes
	api := e.Group("/api")
	{
		// System status
		api.GET("/status", statusHandler.GetStatus)

		// User routes
		user := api.Group("/user")
		{
			user.GET("/settings", userHandler.GetSettings)
			user.POST("/settings", userHandler.SaveSettings)
		}

		// Admin routes (for manual operations)
		admin := api.Group("/admin")
		{
			admin.POST("/sync", syncHandler.TriggerSync)
			admin.GET("/test-api", syncHandler.TestAPI)
			admin.POST("/test-notification", syncHandler.TestNotification)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.QueryProducts)
			products.GET("/", productHandler.QueryProducts)
			products.GET("/blocked", productHandler.GetBlockedProducts)
			products.GET("/:activityId/trend", productHandler.GetPriceTrend)
			products.POST("/:activityId/block", productHandler.BlockProduct)
			products.POST("/unblock/:activityId", productHandler.UnblockProduct)
			products.DELETE("/platform/:platform", productHandler.ClearPlatform)

			// Notification routes - support both with and without trailing slash
			products.POST("/notifications", productHandler.CreateNotification)
			products.POST("/notifications/", productHandler.CreateNotification)
			products.PUT("/notifications/:activityId", productHandler.UpdateNotification)
			products.DELETE("/notifications/:activityId", productHandler.DeleteNotification)
		}

		// External platform routes (webhooks)
		external := api.Group("/external")
		{
			external.POST("/dt/push", externalHandler.HandleDTPush)
		}
	}

	// Serve frontend static files
	// "static" directory contains the built frontend (from Docker build or npm run build)
	staticDir := "static"

	// Serve static assets (JS, CSS, images, etc.)
	e.Static("/assets", staticDir+"/assets")

	// Serve specific static files
	e.File("/vite.svg", staticDir+"/vite.svg")

	// SPA fallback: serve index.html for all unmatched routes
	// This handles client-side routing (React Router, etc.)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if !ok {
			he = &echo.HTTPError{Code: 500, Message: err.Error()}
		}

		// Check if it's a 404 error for non-API routes
		if he.Code == 404 {
			path := c.Request().URL.Path
			// Don't serve index.html for API or health routes that don't exist
			if !strings.HasPrefix(path, "/api") &&
				!strings.HasPrefix(path, "/health") &&
				!strings.HasPrefix(path, "/ready") {
				// Serve index.html for SPA routing
				if fileErr := c.File(staticDir + "/index.html"); fileErr == nil {
					return
				}
			}
		}

		// Default error handler
		c.JSON(he.Code, map[string]interface{}{
			"code":    he.Code,
			"message": he.Message,
		})
	}

	return e
}
