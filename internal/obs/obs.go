package obs

import (
	"context"
	"fmt"

	"agent-project-manager/internal/config"
)

// Init initializes OpenTelemetry tracing and metrics
func Init(cfg config.Config) error {
	serviceName := "agentd"
	serviceVersion := "1.0.0"

	var endpoint string

	// Initialize tracing if enabled
	if cfg.Obs.Tracing.Enabled {
		endpoint = cfg.Obs.Tracing.Endpoint
		if endpoint == "" {
			endpoint = "none" // Disable if not configured
		}
		tp, err := InitTracing(serviceName, serviceVersion, endpoint)
		if err != nil {
			return fmt.Errorf("failed to initialize tracing: %w", err)
		}
		_ = tp // tracerProvider is stored globally
	}

	// Initialize metrics if enabled
	if cfg.Obs.Metrics.Enabled {
		endpoint = cfg.Obs.Metrics.Endpoint
		if endpoint == "" {
			endpoint = "none" // Disable if not configured
		}
		mp, err := InitMetrics(serviceName, serviceVersion, endpoint)
		if err != nil {
			return fmt.Errorf("failed to initialize metrics: %w", err)
		}
		_ = mp // metricProvider is stored globally
	}

	return nil
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

