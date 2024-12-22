package bluetooth

import (
	"context"
	"encoding/json"
)

type blueutilDeviceInfo struct {
	Address   string `json:"address"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	Paired    bool   `json:"paired"`
}

type macosBlueutilBluetoothManager struct{}

func newMacosBlueutilBluetoothManager() macosBlueutilBluetoothManager {
	// validate blueutil is on the path
	if !onPath("blueutil") {
		panic("couldn't find blueutil on the path")
	}

	return macosBlueutilBluetoothManager{}
}

func (m macosBlueutilBluetoothManager) Connect(ctx context.Context, macAddr string) error {
	output, err := execcmd(ctx, "blueutil", "--connect", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

func (m macosBlueutilBluetoothManager) Disconnect(ctx context.Context, macAddr string) error {
	output, err := execcmd(ctx, "blueutil", "--disconnect", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

func (m macosBlueutilBluetoothManager) List(ctx context.Context) ([]BluetoothDevice, error) {
	output, err := execcmd(ctx, "blueutil", "--paired", "--format", "json")
	if err != nil {
		// TODO: process the error to see what went wrong
		return nil, err
	}

	var rawDevices []blueutilDeviceInfo
	err = json.Unmarshal(output, &rawDevices)
	if err != nil {
		return nil, err
	}

	devices := []BluetoothDevice{}
	for _, rawDevice := range rawDevices {
		devices = append(devices, BluetoothDevice{
			Name:      rawDevice.Name,
			MacAddr:   rawDevice.Address,
			Connected: rawDevice.Connected,
		})
	}

	return devices, nil
}

func (m macosBlueutilBluetoothManager) Get(ctx context.Context, macAddr string) (BluetoothDevice, error) {
	var device BluetoothDevice
	output, err := execcmd(ctx, "blueutil", "--info", macAddr, "--format", "json")
	if err != nil {
		// TODO: process the error to see what went wrong
		return device, err
	}

	var rawDevice blueutilDeviceInfo
	err = json.Unmarshal(output, &rawDevice)
	if err != nil {
		return device, err
	}

	device.Name = rawDevice.Name
	device.MacAddr = rawDevice.Address
	device.Connected = rawDevice.Connected
	return device, nil
}

func (m macosBlueutilBluetoothManager) IsConnected(ctx context.Context, macAddr string) (bool, error) {
	output, err := execcmd(ctx, "blueutil", "--is-connected", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return false, err
	}

	// --is-connected returns '1' if the device is connected and 0 if not
	return output[0] == '1', nil
}
