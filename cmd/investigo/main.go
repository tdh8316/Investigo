package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tdh8316/Investigo/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(app.Run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
