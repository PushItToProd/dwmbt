package bluetooth

import (
	"fmt"
	"os/exec"
	"runtime"
)

// onPath returns true if the named executable is on the $PATH. This will return
// false even if the error is not ErrNotFound, so issues could potentially arise
// in edge cases.
func onPath(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

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
	switch os := runtime.GOOS; os {
	case "darwin":
		return newMacosBlueutilBluetoothManager()
	case "linux":
		return newLinuxBluetoothctlBluetoothManager()
	default:
		panic(fmt.Sprintf("unsupported OS: %s", os))
	}
}
