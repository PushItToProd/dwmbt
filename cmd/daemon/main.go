package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pushittoprod/bt-daemon/pkg/daemon"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// handle termination
	go func() {
		daemon.RunServer(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("got shutdown signal - cancelling server", "sig", sig)
	cancel()
}
