package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agent-project-manager/internal/agentd"
	"agent-project-manager/internal/config"
	"agent-project-manager/internal/logger"

	_ "agent-project-manager/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// logger may not be initialized yet
		_, _ = os.Stderr.WriteString("agentd: failed to load config: " + err.Error() + "\n")
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		_, _ = os.Stderr.WriteString("agentd: invalid config: " + err.Error() + "\n")
		os.Exit(1)
	}

	app, err := agentd.Init(cfg, agentd.Options{
		ShutdownTimeout: 5 * time.Second,
		MaxDBRetries:    10,
		DBRetryDelay:    2 * time.Second,
		MigrationsDir:   "migrations",
	})
	if err != nil {
		logger.Init()
		logger.Fatalf("agentd: init failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Infof("agentd: starting HTTP server on %s", app.Server.Addr)
		errCh <- app.Server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("agentd: shutting down server...")
		if err := app.Shutdown(context.Background()); err != nil {
			logger.Fatalf("agentd: shutdown failed: %v", err)
		}

	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("agentd: server failed: %v", err)
		}
		_ = app.Shutdown(context.Background())
	}

	logger.Info("agentd: server exited")
}
