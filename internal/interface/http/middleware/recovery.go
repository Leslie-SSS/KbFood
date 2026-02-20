package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Recovery returns a panic recovery middleware for Echo
func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic
					log.Error().
						Str("path", c.Request().URL.Path).
						Str("method", c.Request().Method).
						Str("stack", string(debug.Stack())).
						Msg(fmt.Sprintf("panic: %v", r))

					// Return 500 error to client
					err = echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
				}
			}()

			return next(c)
		}
	}
}
