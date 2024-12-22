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
		devinfo := fmt.Sprintf("%s (%s)", device.Name, device.MacAddr)
		if device.Connected {
			fmt.Printf("\x1b[97m%s \x1b[1;97m[connected]\x1b[0m", devinfo)
		} else {
			fmt.Printf("%s", devinfo)
		}
		fmt.Println()
	}
}
