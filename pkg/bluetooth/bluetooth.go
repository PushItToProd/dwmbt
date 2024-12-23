package bluetooth

import (
	"context"
	"fmt"
	"runtime"
)

// A BluetoothDevice represents a single Bluetooth device connected to or known by a host.
type BluetoothDevice struct {
	Name      string
	MacAddr   string
	Connected bool
}

// A BluetoothManager provides some means of managing Bluetooth devices connected to the host.
type BluetoothManager interface {
	Connect(ctx context.Context, macAddr string) error
	Disconnect(ctx context.Context, macAddr string) error
	List(ctx context.Context) ([]BluetoothDevice, error)
	Get(ctx context.Context, macAddr string) (BluetoothDevice, error)
	IsConnected(ctx context.Context, macAddr string) (bool, error)
}

// NewBluetoothManager returns an appropriate BluetoothManager for the current platform.
func NewBluetoothManager() BluetoothManager {
	// XXX: Instead of hardcoding the OS detection here, we could have a registry of managers with a function that
	// determines whether the manager is appropriate for the current platform. Then we could iterate over the registry
	// and return the first manager that is appropriate. This would let us have more flexibility, e.g. in case a
	// particular platform has multiple possible managers. That said, this current approach is easy to understand and
	// allows us to easily tell the user if they don't have the required command installed.
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
