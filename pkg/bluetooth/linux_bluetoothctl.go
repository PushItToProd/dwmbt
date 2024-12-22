package bluetooth

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type linuxBluetoothctlBluetoothManager struct{}

func newLinuxBluetoothctlBluetoothManager() linuxBluetoothctlBluetoothManager {
	// validate bluetoothctl is on the path
	if !onPath("bluetoothctl") {
		panic("couldn't find bluetoothctl on the path")
	}

	return linuxBluetoothctlBluetoothManager{}
}

func (m linuxBluetoothctlBluetoothManager) Connect(ctx context.Context, macAddr string) error {
	output, err := runCmd(ctx, "bluetoothctl", "connect", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

func (m linuxBluetoothctlBluetoothManager) Disconnect(ctx context.Context, macAddr string) error {
	output, err := runCmd(ctx, "bluetoothctl", "disconnect", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

var bluetoothctlDeviceListRegex = regexp.MustCompile(`Device ([0-9A-Za-z:]+) (.*)`)

func (m linuxBluetoothctlBluetoothManager) List(ctx context.Context) ([]BluetoothDevice, error) {
	output, err := runCmd(ctx, "bluetoothctl", "devices")
	if err != nil {
		// TODO: process the error to see what went wrong
		return nil, err
	}

	lines := bytes.Split(output, []byte("\n"))
	devices := []BluetoothDevice{}
	for _, line := range lines {
		ms := bluetoothctlDeviceListRegex.FindSubmatch(line)
		if ms == nil {
			lineStr := string(line)
			if strings.TrimSpace(lineStr) != "" {
				log.Printf("failed to match bluetoothctl device line: %s", string(line))
			}
			continue
		}
		macAddr := string(ms[1])
		name := string(ms[2])
		connected, err := m.IsConnected(ctx, macAddr)
		if err != nil {
			log.Printf("failed to get connection status for device %q with MAC %s: %v", name, macAddr, err)
			connected = false
		}
		devices = append(devices, BluetoothDevice{
			Name:      name,
			MacAddr:   macAddr,
			Connected: connected,
		})
	}

	return devices, nil
}

type bluetoothctlDeviceInfo struct {
	Name      *string
	MacAddr   *string
	Connected *bool
}

func (i bluetoothctlDeviceInfo) Validate() error {
	var missingFields []string
	if i.Name == nil {
		missingFields = append(missingFields, "Name")
	}
	if i.MacAddr == nil {
		missingFields = append(missingFields, "MacAddr")
	}
	if i.Connected == nil {
		missingFields = append(missingFields, "Connected")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

func parseDeviceInfo(output []byte) (bluetoothctlDeviceInfo, error) {
	var device bluetoothctlDeviceInfo
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch fields[0] {
		case "Name:":
			device.Name = &fields[1]
		case "Address:":
			device.MacAddr = &fields[1]
		case "Connected:":
			connected := fields[1] == "yes"
			device.Connected = &connected
		}
	}
	if err := device.Validate(); err != nil {
		return bluetoothctlDeviceInfo{}, fmt.Errorf("invalid bluetoothctl output: %w", err)
	}
	return device, nil
}

func (m linuxBluetoothctlBluetoothManager) Get(ctx context.Context, macAddr string) (BluetoothDevice, error) {
	output, err := runCmd(ctx, "bluetoothctl", "info", macAddr)
	if err != nil {
		// TODO: process the error to see what went wrong
		return BluetoothDevice{}, err
	}

	device, err := parseDeviceInfo(output)
	if err != nil {
		return BluetoothDevice{}, err
	}
	return BluetoothDevice{
		Name:      *device.Name,
		MacAddr:   *device.MacAddr,
		Connected: *device.Connected,
	}, nil
}

func (m linuxBluetoothctlBluetoothManager) IsConnected(ctx context.Context, macAddr string) (bool, error) {
	device, err := m.Get(ctx, macAddr)
	if err != nil {
		return false, err
	}
	return device.Connected, nil
}
