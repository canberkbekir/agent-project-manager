package obs

import (
	"context"
	"fmt"

	"agent-project-manager/internal/config"
)

// Init initializes OpenTelemetry tracing and metrics
// Returns the Prometheus metrics path if Prometheus is enabled, empty string otherwise
func Init(cfg config.Config) (string, error) {
	serviceName := "agentd"
	serviceVersion := "1.0.0"

	var endpoint string
	var prometheusPath string

	// Initialize tracing if enabled
	if cfg.Obs.Tracing.Enabled {
		endpoint = cfg.Obs.Tracing.Endpoint
		if endpoint == "" {
			endpoint = "none" // Disable if not configured
		}
		tp, err := InitTracing(serviceName, serviceVersion, endpoint)
		if err != nil {
			return "", fmt.Errorf("failed to initialize tracing: %w", err)
		}
		_ = tp // tracerProvider is stored globally
	}

	// Initialize metrics if enabled
	if cfg.Obs.Metrics.Enabled {
		// Initialize Prometheus metrics if enabled (for Grafana)
		if cfg.Obs.Metrics.PrometheusEnabled {
			_, err := InitPrometheusMetrics(serviceName, serviceVersion)
			if err != nil {
				return "", fmt.Errorf("failed to initialize Prometheus metrics: %w", err)
			}
			// Get Prometheus path from config, default to /metrics
			prometheusPath = cfg.Obs.Metrics.PrometheusPath
			if prometheusPath == "" {
				prometheusPath = "/metrics"
			}
		} else {
			// Initialize OTLP metrics if Prometheus is not enabled
			endpoint = cfg.Obs.Metrics.Endpoint
			if endpoint == "" {
				endpoint = "none" // Disable if not configured
			}
			mp, err := InitMetrics(serviceName, serviceVersion, endpoint)
			if err != nil {
				return "", fmt.Errorf("failed to initialize metrics: %w", err)
			}
			_ = mp // metricProvider is stored globally
		}
	}

	return prometheusPath, nil
}

// Shutdown gracefully shuts down OpenTelemetry
func Shutdown(ctx context.Context) error {
	var errs []error

	if err := ShutdownTracing(ctx); err != nil {
		errs = append(errs, fmt.Errorf("failed to shutdown tracing: %w", err))
	}

	if err := ShutdownMetrics(ctx); err != nil {
		errs = append(errs, fmt.Errorf("failed to shutdown metrics: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	return nil
}

