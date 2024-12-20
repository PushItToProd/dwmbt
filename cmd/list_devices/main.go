package main

import (
	"fmt"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
)

func main() {
	btm := bluetooth.NewBluetoothManager()
	devices, err := btm.List()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for _, device := range devices {
		fmt.Printf("%+v\n", device)
	}
}
