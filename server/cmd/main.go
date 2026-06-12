package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := NewApp(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tryit: fatal:", err)
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "tryit: fatal:", err)
		os.Exit(1)
	}
}
