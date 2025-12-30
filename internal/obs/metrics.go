package obs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

var (
	metricProvider      *metric.MeterProvider
	prometheusExporter  *otelprom.Exporter
	prometheusHandler   http.Handler
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

// InitPrometheusMetrics initializes Prometheus metrics exporter for Grafana
func InitPrometheusMetrics(serviceName, serviceVersion string) (http.Handler, error) {
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

	// Create Prometheus registry
	reg := prometheus.NewRegistry()

	// Create Prometheus exporter with namespace
	exporter, err := otelprom.New(
		otelprom.WithNamespace("agentd"),
		otelprom.WithRegisterer(reg),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider with Prometheus reader
	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	// Set global meter provider
	otel.SetMeterProvider(mp)

	// Create HTTP handler for Prometheus metrics
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	metricProvider = mp
	prometheusExporter = exporter
	prometheusHandler = handler

	return handler, nil
}

// GetPrometheusHandler returns the Prometheus metrics HTTP handler
func GetPrometheusHandler() http.Handler {
	if prometheusHandler == nil {
		return http.NotFoundHandler()
	}
	return prometheusHandler
}

// ShutdownMetrics gracefully shuts down the meter provider
func ShutdownMetrics(ctx context.Context) error {
	if metricProvider != nil {
		return metricProvider.Shutdown(ctx)
	}
	return nil
}

