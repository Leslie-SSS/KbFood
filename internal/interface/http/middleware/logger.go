package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Logger returns a logger middleware
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Get request ID from header or generate one
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			if requestID == "" {
				requestID = "unknown"
			}

			// Call next handler
			err := next(c)

			// Log after request completes
			log.Info().
				Str("request_id", requestID).
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Str("query", c.Request().URL.RawQuery).
				Int("status", c.Response().Status).
				Int64("size", c.Response().Size).
				Dur("duration", time.Since(start)).
				Msg("HTTP request")

			return err
		}
	}
}
