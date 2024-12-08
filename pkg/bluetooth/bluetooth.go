package bluetooth

import (
	"fmt"
	"runtime"
)

type BluetoothDevice struct {
	Name      string
	MacAddr   string
	Connected bool
}

type BluetoothManager interface {
	Connect(macAddr string) error
	Disconnect(macAddr string) error
	List() ([]BluetoothDevice, error)
	Get(macAddr string) (BluetoothDevice, error)
	IsConnected(macAddr string) (bool, error)
}

func NewBluetoothManager() BluetoothManager {
	var manager BluetoothManager
	switch os := runtime.GOOS; os {
	case "darwin":
		manager = newMacosBlueutilBluetoothManager()
	case "linux":
		manager = newLinuxBluetoothctlBluetoothManager()
	default:
		panic(fmt.Sprintf("unsupported OS: %s", os))
	}
	// Wrap the manager in a normalizing manager to ensure all MAC addresses are validated and formatted consistently.
	// This saves us from having to do this in every implementation.
	return saferBluetoothManager{inner: manager}
}
