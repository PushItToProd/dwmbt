package main

import (
	"context"
	"log"
	"os"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
)

func main() {
	// get device address from CLI args
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <device-id>", os.Args[0])
	}
	macAddr := os.Args[1]

	btm := bluetooth.NewBluetoothManager()

	log.Printf("disconnecting %q", macAddr)
	ctx := context.Background()
	if err := btm.Disconnect(ctx, macAddr); err != nil {
		log.Fatalf("failed to disconect: %v", err)
	}
}
