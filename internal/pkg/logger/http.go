package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// RequestIDHeader is the header name for request ID
const RequestIDHeader = "X-Request-ID"

// Middleware returns a logging middleware for chi
func Middleware(next http.Handler) http.Handler {
	return hlog.NewHandler(Logger)(next)
}

// AccessHandler returns an access log handler
func AccessHandler(next http.Handler) http.Handler {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("HTTP request")
	})(next)
}

// RequestIDHandler adds a request ID to the context
func RequestIDHandler(next http.Handler) http.Handler {
	return hlog.RequestIDHandler(RequestIDHeader, "req-id")(next)
}

// RequestIDFromContext returns the request ID from the context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := hlog.IDFromCtx(ctx); ok {
		return id.String()
	}
	return ""
}

// LoggerFromContext returns a logger from the context
func LoggerFromContext(ctx context.Context) zerolog.Logger {
	l := zerolog.Ctx(ctx)
	if l == nil {
		return Logger
	}
	return *l
}
