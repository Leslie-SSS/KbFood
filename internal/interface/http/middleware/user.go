package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// UserIDHeader is the HTTP header for user identification
	UserIDHeader = "X-User-ID"
	// LegacyUserIDHeader carries the legacy Bark-key-derived user identifier.
	LegacyUserIDHeader = "X-Legacy-User-ID"
	// UserIDContextKey is the context key for user ID
	UserIDContextKey = "userID"
	// LegacyUserIDContextKey is the context key for the legacy user ID.
	LegacyUserIDContextKey = "legacyUserID"
)

// UserExtractor extracts user ID from request header and adds to context
func UserExtractor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := strings.TrimSpace(c.Request().Header.Get(UserIDHeader))
			legacyUserID := strings.TrimSpace(c.Request().Header.Get(LegacyUserIDHeader))
			c.Set(UserIDContextKey, userID)
			c.Set(LegacyUserIDContextKey, legacyUserID)
			return next(c)
		}
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c echo.Context) string {
	if id, ok := c.Get(UserIDContextKey).(string); ok {
		return id
	}
	return ""
}

// GetLegacyUserID retrieves the legacy user ID from context.
func GetLegacyUserID(c echo.Context) string {
	if id, ok := c.Get(LegacyUserIDContextKey).(string); ok {
		return id
	}
	return ""
}
