// tryit local server entry point. Composition root per IMPL §11.
//
// Phase 1: envStore (secret), file-backed pairStore, /health + /pair + /execute.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/handler"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integration/storage"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integration/target"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/execute"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/pair"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/httpclient"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "tryit: fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New(cfg.LogLevel)
	hc := httpclient.New(cfg)

	pairStore, freshlyGenerated, err := storage.NewPairStore(cfg.PairFile)
	if err != nil {
		return fmt.Errorf("pair store: %w", err)
	}

	secretStore, err := storage.NewSecretStore(cfg)
	if err != nil {
		return fmt.Errorf("secret store: %w", err)
	}

	secretSvc := secret.New(secretStore)
	pairSvc := pair.New(pairStore)
	targetClient := target.New(hc, cfg)
	execSvc := execute.New(secretSvc, targetClient)

	mux := handler.NewRouter(handler.Deps{
		Pair:    pairStore,
		Health:  handler.NewHealthHandler(pairStore),
		Pairing: handler.NewPairHandler(pairSvc),
		Execute: handler.NewExecuteHandler(execSvc),
		Logger:  log,
		Host:    cfg.HostHeader(),
	})
	srv := handler.NewServer(cfg, log, mux)

	announcePairingToken(pairStore.Token(), pairStore.BoundOrigin(), freshlyGenerated)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		log.Info("shutdown signal received, draining")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}
}

func announcePairingToken(token, boundOrigin string, freshlyGenerated bool) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "──────────────────────── tryit pairing ────────────────────────")
	if boundOrigin != "" {
		fmt.Fprintf(os.Stdout, " Already paired with: %s\n", boundOrigin)
	} else if freshlyGenerated {
		fmt.Fprintln(os.Stdout, " Fresh pairing token (paste into the extension panel):")
	} else {
		fmt.Fprintln(os.Stdout, " Existing pairing token (paste into the extension panel):")
	}
	fmt.Fprintf(os.Stdout, "   %s\n", token)
	fmt.Fprintln(os.Stdout, " Reset:  make reset-pairing")
	fmt.Fprintln(os.Stdout, "───────────────────────────────────────────────────────────────")
}
