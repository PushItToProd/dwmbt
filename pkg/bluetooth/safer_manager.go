package bluetooth

import (
	"context"
	"fmt"
)

var ErrInvalidMac = fmt.Errorf("invalid MAC address")

// saferBluetoothManager wraps an underlying BluetoothManager, providing standardized validation of MAC addresses passed
// as arguments.
//
// This is primarily a security feature. We pass the MAC address as an argument to external commands, and while we use
// exec.Command to avoid shell injection, letting an attacker pass an arbitrary string creates the risk of buffer
// overflows or other attacks targeting vulnerabilities in the underlying commands that we use.
//
// TODO: To be extra safe, we should also check that provided MAC addresses are in the list of devices we know about.
// This would all but ensure that a MAC address passed to us is valid and safe to use. While I doubt an attacker could
// fit a real attack into a string that looks like a valid MAC address, it's better to be safe than sorry.
type saferBluetoothManager struct {
	inner BluetoothManager
}

func (m saferBluetoothManager) Connect(ctx context.Context, macAddr string) error {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return ErrInvalidMac
	}
	return m.inner.Connect(ctx, mac)
}

func (m saferBluetoothManager) Disconnect(ctx context.Context, macAddr string) error {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return ErrInvalidMac
	}
	return m.inner.Disconnect(ctx, mac)
}

func (m saferBluetoothManager) List(ctx context.Context) ([]BluetoothDevice, error) {
	return m.inner.List(ctx)
}

func (m saferBluetoothManager) Get(ctx context.Context, macAddr string) (BluetoothDevice, error) {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return BluetoothDevice{}, ErrInvalidMac
	}
	return m.inner.Get(ctx, mac)
}

func (m saferBluetoothManager) IsConnected(ctx context.Context, macAddr string) (bool, error) {
	mac, ok := NormalizeMac(macAddr)
	if !ok {
		return false, ErrInvalidMac
	}
	return m.inner.IsConnected(ctx, mac)
}
