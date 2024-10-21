package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
	"github.com/pushittoprod/bt-daemon/pkg/daemon"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// handle termination
	go func() {
		d := daemon.Daemon{
			ServeAddr:        ":8080",
			BluetoothManager: bluetooth.NewBluetoothManager(),
		}
		d.RunServer(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("got shutdown signal - stopping server", "sig", sig)
	cancel()
}
