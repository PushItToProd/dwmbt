package main

import (
	"context"
	"fmt"
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
	ctx := context.Background()
	btd, err := btm.Get(ctx, macAddr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(btd)
}
