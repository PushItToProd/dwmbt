package main

import (
	"log"
	"os"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
)

func main() {
	// get device address from CLI args
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <device-id>", os.Args[0])
	}
	mac := os.Args[1]

	btm := bluetooth.NewBluetoothManager()

	log.Printf("disconnecting %q", mac)
	if err := btm.Disconnect(mac); err != nil {
		log.Fatalf("failed to disconect: %v", err)
	}
}
