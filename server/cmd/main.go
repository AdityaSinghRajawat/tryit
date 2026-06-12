package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := NewApp(ctx)
	if err != nil {
		utils.LogErrorWithoutCtx(err)
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		utils.LogErrorWithoutCtx(err)
		os.Exit(1)
	}
}
