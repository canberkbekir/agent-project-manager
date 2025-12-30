package obs

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

var (
	metricProvider *metric.MeterProvider
)

// InitMetrics initializes OpenTelemetry metrics
func InitMetrics(serviceName, serviceVersion, endpoint string) (*metric.MeterProvider, error) {
	ctx := context.Background()

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter and reader
	var reader metric.Reader
	if endpoint == "" || endpoint == "none" {
		// No-op reader for when metrics is disabled
		reader = metric.NewManualReader() // Manual reader that doesn't export
	} else {
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(endpoint),
			otlpmetrichttp.WithInsecure(),
		}

		// Support OTLP_HTTP_ENDPOINT environment variable
		if envEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); envEndpoint != "" {
			opts = []otlpmetrichttp.Option{
				otlpmetrichttp.WithEndpoint(envEndpoint),
			}
		}

		exp, err := otlpmetrichttp.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create metric exporter: %w", err)
		}
		reader = metric.NewPeriodicReader(exp)
	}

	// Create meter provider
	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(mp)

	metricProvider = mp

	return mp, nil
}

// ShutdownMetrics gracefully shuts down the meter provider
func ShutdownMetrics(ctx context.Context) error {
	if metricProvider != nil {
		return metricProvider.Shutdown(ctx)
	}
	return nil
}

