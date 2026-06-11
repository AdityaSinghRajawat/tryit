// app.go is the bootstrap: load config, build the router (which wires
// integrations + services + handlers), and run the http.Server with a
// graceful drain on SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/routes"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func Run() error {
	if err := config.Init(); err != nil {
		return err
	}
	log := utils.NewLogger(config.GetLogLevel())

	mux, err := routes.NewRoutes(log)
	if err != nil {
		return fmt.Errorf("routes: %w", err)
	}

	srv := &http.Server{
		Addr:              config.GetListenAddr(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      config.GetExecTimeout() + 30*time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	log.Info("listening", "addr", srv.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		log.Info("shutdown signal received, draining")
		ctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownDrain())
		defer cancel()
		return srv.Shutdown(ctx)
	}
}
