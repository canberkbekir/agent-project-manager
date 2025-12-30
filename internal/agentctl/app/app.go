package app

import (
	"agent-project-manager/internal/config"
	"context"
	"errors"
	"io"
)

type App struct {
	Cfg config.Config
	Out io.Writer
	Err io.Writer
}

type ctxKey int

const appKey ctxKey = iota

func WithApp(ctx context.Context, a *App) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, appKey, a)
}

func FromContext(ctx context.Context) (*App, error) {
	if ctx == nil {
		return nil, errors.New("missing context")
	}
	v := ctx.Value(appKey)
	if v == nil {
		return nil, errors.New("app not initialized (context missing)")
	}
	a, ok := v.(*App)
	if !ok || a == nil {
		return nil, errors.New("invalid app in context")
	}
	return a, nil
}
