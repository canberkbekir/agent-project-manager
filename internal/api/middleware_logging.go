package api

import (
	"net/http"
	"time"

	"agent-project-manager/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger is a middleware that logs HTTP requests using the application's logger
func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code and response size
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Call the next handler
			next.ServeHTTP(ww, r)

			// Calculate duration
			duration := time.Since(start)

			// Get request ID from context (set by middleware.RequestID)
			requestID := middleware.GetReqID(r.Context())

			// Get real IP (set by middleware.RealIP)
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			} else if realIP := r.Header.Get("X-Real-Ip"); realIP != "" {
				ip = realIP
			}

			// Build log fields
			fields := logger.WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      ww.Status(),
				"duration":    duration.String(),
				"duration_ms": duration.Milliseconds(),
				"bytes":       ww.BytesWritten(),
				"ip":          ip,
			})

			if requestID != "" {
				fields = fields.WithField("request_id", requestID)
			}

			if r.URL.RawQuery != "" {
				fields = fields.WithField("query", r.URL.RawQuery)
			}

			if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
				fields = fields.WithField("user_agent", userAgent)
			}

			// Log based on status code
			status := ww.Status()
			switch {
			case status >= 500:
				fields.Error("HTTP request failed")
			case status >= 400:
				fields.Warn("HTTP request client error")
			default:
				fields.Info("HTTP request")
			}
		})
	}
}

