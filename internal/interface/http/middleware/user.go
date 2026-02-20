package middleware

import "github.com/labstack/echo/v4"

const (
	// UserIDHeader is the HTTP header for user identification
	UserIDHeader = "X-User-ID"
	// UserIDContextKey is the context key for user ID
	UserIDContextKey = "userID"
)

// UserExtractor extracts user ID from request header and adds to context
func UserExtractor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Request().Header.Get(UserIDHeader)
			c.Set(UserIDContextKey, userID)
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
