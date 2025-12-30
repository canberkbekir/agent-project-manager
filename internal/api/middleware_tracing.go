package api

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TracingMiddleware returns an OpenTelemetry HTTP middleware
func TracingMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(
		next,
		"",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			if operation != "" {
				return operation
			}
			return r.Method + " " + r.URL.Path
		}),
	)
}

