package agentd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"agent-project-manager/internal/api"
	"agent-project-manager/internal/config"
	"agent-project-manager/internal/logger"
	"agent-project-manager/internal/obs"
	"agent-project-manager/internal/state"
)

type App struct {
	Store    state.Store
	Server   *http.Server
	Shutdown func(ctx context.Context) error
}

type Options struct {
	ShutdownTimeout time.Duration
	MaxDBRetries    int
	DBRetryDelay    time.Duration
	MigrationsDir   string
}

// Init wires logger, DB store (with retry), migrations, OpenTelemetry, and HTTP server.
func Init(cfg config.Config, opts Options) (*App, error) {
	// Defaults
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 5 * time.Second
	}
	if opts.MaxDBRetries == 0 {
		opts.MaxDBRetries = 5
	}
	if opts.DBRetryDelay == 0 {
		opts.DBRetryDelay = 2 * time.Second
	}
	if opts.MigrationsDir == "" {
		opts.MigrationsDir = "migrations"
	}

	logger.Init()

	// DB connect with retry
	var store state.Store
	var lastErr error
	for i := 0; i < opts.MaxDBRetries; i++ {
		s, err := state.NewStore(cfg.State.ConnectionString)
		if err == nil {
			store = s
			logger.Info("agentd: successfully connected to database")
			lastErr = nil
			break
		}
		lastErr = err
		if i < opts.MaxDBRetries-1 {
			logger.Warnf("agentd: failed to connect to database (attempt %d/%d): %v, retrying in %v...",
				i+1, opts.MaxDBRetries, err, opts.DBRetryDelay)
			time.Sleep(opts.DBRetryDelay)
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}

	// Migrations
	if err := state.RunMigrations(store, opts.MigrationsDir); err != nil {
		_ = store.Close()
		return nil, err
	}
	logger.Info("agentd: database migrations completed successfully")

	// OpenTelemetry
	prometheusPath, err := obs.Init(cfg)
	if err != nil {
		logger.Warnf("agentd: failed to initialize OpenTelemetry: %v", err)
	} else {
		if cfg.Obs.Tracing.Enabled {
			api.EnableTracing = true
		}
		if prometheusPath != "" {
			api.PrometheusMetricsPath = prometheusPath
			logger.Infof("agentd: Prometheus metrics endpoint enabled at %s", prometheusPath)
		}
	}

	srv := &http.Server{
		Addr:    cfg.API.Addr,
		Handler: api.Router(store),
	}

	app := &App{
		Store:  store,
		Server: srv,
		Shutdown: func(ctx context.Context) error {
			// stop HTTP server first
			shutdownCtx, cancel := context.WithTimeout(ctx, opts.ShutdownTimeout)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}

			// shutdown OTel (if it was initialized, obs.Shutdown should be safe/no-op per your impl)
			if err := obs.Shutdown(ctx); err != nil {
				logger.Errorf("agentd: failed to shutdown OpenTelemetry: %v", err)
			}

			// close DB
			if err := store.Close(); err != nil {
				logger.Errorf("agentd: failed to close database store: %v", err)
			}
			return nil
		},
	}

	return app, nil
}
