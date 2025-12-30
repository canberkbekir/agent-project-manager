// @title           Agent Project Manager API
// @version         1.0
// @description     REST API for the Agent Project Manager daemon. Orchestrates agent-based LLM workflows.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3333
// @BasePath  /v1

// @schemes   http https

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agent-project-manager/internal/api"
	"agent-project-manager/internal/config"
	"agent-project-manager/internal/logger"
	"agent-project-manager/internal/obs"
	"agent-project-manager/internal/state"

	_ "agent-project-manager/docs" // Swagger docs
)

const shutdownTimeout = 5 * time.Second

func main() {
	// Log config path before loading (for debugging)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	_, _ = fmt.Fprintf(os.Stdout, "agentd: loading config from: %s\n", configPath)

	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "agentd: failed to load config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "agentd: invalid config: %v\n", err)
		os.Exit(1)
	}

	// Log connection string before logger init (for debugging)
	_, _ = fmt.Fprintf(os.Stdout, "agentd: database connection string: %s\n", maskConnectionString(cfg.State.ConnectionString))
	_, _ = fmt.Fprintf(os.Stdout, "agentd: config file loaded successfully\n")
	_, _ = fmt.Fprintf(os.Stdout, "agentd: API addr: %s\n", cfg.API.Addr)

	logger.Init()

	logger.Infof("agentd: loaded config from %s", configPath)
	logger.Infof("agentd: database connection string: %s", maskConnectionString(cfg.State.ConnectionString))

	// Initialize database store with retry mechanism
	var store state.Store
	maxRetries := 10
	retryDelay := 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		var err error
		store, err = state.NewStore(cfg.State.ConnectionString)
		if err == nil {
			logger.Info("agentd: successfully connected to database")
			break
		}
		if i < maxRetries-1 {
			logger.Warnf("agentd: failed to connect to database (attempt %d/%d): %v, retrying in %v...", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		} else {
			logger.Fatalf("agentd: failed to initialize database store after %d attempts: %v", maxRetries, err)
		}
	}
	defer func() {
		if err := store.Close(); err != nil {
			logger.Errorf("agentd: failed to close database store: %v", err)
		}
	}()

	// Run database migrations
	if err := state.RunMigrations(store, "migrations"); err != nil {
		logger.Fatalf("agentd: failed to run database migrations: %v", err)
	}
	logger.Info("agentd: database migrations completed successfully")

	// Initialize OpenTelemetry
	prometheusPath, err := obs.Init(cfg)
	if err != nil {
		logger.Warnf("agentd: failed to initialize OpenTelemetry: %v", err)
	} else {
		// Enable tracing middleware if tracing is enabled
		if cfg.Obs.Tracing.Enabled {
			api.EnableTracing = true
		}
		// Set Prometheus metrics path if enabled
		if prometheusPath != "" {
			api.PrometheusMetricsPath = prometheusPath
			logger.Infof("agentd: Prometheus metrics endpoint enabled at %s", prometheusPath)
		}
		defer func() {
			if err := obs.Shutdown(context.Background()); err != nil {
				logger.Errorf("agentd: failed to shutdown OpenTelemetry: %v", err)
			}
		}()
	}

	srv := &http.Server{
		Addr:    cfg.API.Addr,
		Handler: api.Router(store),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Infof("agentd: starting HTTP server on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("agentd: shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Fatalf("agentd: shutdown failed: %v", err)
		}

	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("agentd: server failed: %v", err)
		}
	}

	logger.Info("agentd: server exited")
}

// maskConnectionString masks the password in connection string for logging
func maskConnectionString(connStr string) string {
	// Simple masking: replace password part
	// Format: postgres://user:password@host:port/db
	// We'll just show user@host:port/db
	if len(connStr) < 10 {
		return "***"
	}
	// Find @ symbol and show only user@host part
	atIdx := -1
	for i := 0; i < len(connStr); i++ {
		if connStr[i] == '@' {
			atIdx = i
			break
		}
	}
	if atIdx > 0 {
		// Find : before @ to mask password
		colonIdx := -1
		for i := 0; i < atIdx; i++ {
			if connStr[i] == ':' {
				colonIdx = i
				break
			}
		}
		if colonIdx > 0 {
			return connStr[:colonIdx+1] + "***" + connStr[atIdx:]
		}
	}
	return "***"
}
