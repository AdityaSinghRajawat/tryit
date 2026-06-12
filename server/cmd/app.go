package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/routes"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

type App struct {
	router http.Handler
	redis  *redis.Client
}

func NewApp(_ context.Context) (*App, error) {
	if err := config.Init(); err != nil {
		return nil, err
	}

	app := &App{}
	app.initClients()
	if err := app.loadRoutes(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) initClients() {
	if config.GetRedisAddr() == "" || !config.GetCacheEnabled() {
		return
	}
	a.redis = redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.GetRedisPassword(),
		DB:       config.GetRedisDB(),
	})
}

func (a *App) loadRoutes() error {
	r, err := routes.NewRoutes(a.redis)
	if err != nil {
		return err
	}
	a.router = r
	return nil
}

func (a *App) Start(ctx context.Context) error {
	if a.redis != nil {
		if err := a.redis.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("redis ping: %w", err)
		}
		utils.LogInfoWithoutCtx("redis connected", zap.String("addr", config.GetRedisAddr()))
		defer func() { _ = a.redis.Close() }()
	}

	srv := &http.Server{
		Addr:              config.GetListenAddr(),
		Handler:           a.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      config.GetExecTimeout() + 30*time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	utils.LogInfoWithoutCtx("listening", zap.String("addr", srv.Addr))

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		utils.LogInfoWithoutCtx("shutdown signal received, draining")
		sctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownDrain())
		defer cancel()
		return srv.Shutdown(sctx)
	}
}
