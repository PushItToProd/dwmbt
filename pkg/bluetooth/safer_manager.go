package bluetooth

import "fmt"

var ErrInvalidMac = fmt.Errorf("invalid MAC address")

// saferBluetoothManager wraps an underlying BluetoothManager, providing standardized validation of MAC addresses passed
// as arguments.
//
// This is primarily a security feature. We pass the MAC address as an argument to external commands, and while we use
// exec.Command to avoid shell injection, letting an attacker pass an arbitrary string creates the risk of buffer
// overflows or other attacks targeting vulnerabilities in the underlying commands that we use.
//
// TODO: we should also check that provided MAC addresses are in the list of devices we know about. This would eliminate
// the
type saferBluetoothManager struct {
	inner BluetoothManager
}

func (m saferBluetoothManager) Connect(macAddr string) error {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return ErrInvalidMac
	}
	return m.inner.Connect(mac)
}

func (m saferBluetoothManager) Disconnect(macAddr string) error {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return ErrInvalidMac
	}
	return m.inner.Disconnect(mac)
}

func (m saferBluetoothManager) List() ([]BluetoothDevice, error) {
	return m.inner.List()
}

func (m saferBluetoothManager) Get(macAddr string) (BluetoothDevice, error) {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return BluetoothDevice{}, ErrInvalidMac
	}
	return m.inner.Get(mac)
}

func (m saferBluetoothManager) IsConnected(macAddr string) (bool, error) {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return false, ErrInvalidMac
	}
	return m.inner.IsConnected(mac)
}
