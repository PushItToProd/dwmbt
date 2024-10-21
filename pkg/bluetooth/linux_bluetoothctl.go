package bluetooth

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
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

func (m linuxBluetoothctlBluetoothManager) Connect(macAddr string) error {
	cmd := exec.Command("bluetoothctl", "connect", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

func (m linuxBluetoothctlBluetoothManager) Disconnect(macAddr string) error {
	cmd := exec.Command("bluetoothctl", "disconnect", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

var bluetoothctlDeviceListRegex = regexp.MustCompile(`Device ([0-9A-Za-z:]+) (.*)`)

func (m linuxBluetoothctlBluetoothManager) List() ([]BluetoothDevice, error) {
	cmd := exec.Command("bluetoothctl", "devices")
	output, err := cmd.Output()
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
		connected, err := m.IsConnected(macAddr)
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

func (m linuxBluetoothctlBluetoothManager) Get(macAddr string) (BluetoothDevice, error) {
	var device bluetoothctlDeviceInfo
	cmd := exec.Command("bluetoothctl", "info", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return BluetoothDevice{}, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch fields[0] {
		case "Device":
			device.MacAddr = &fields[1]
		case "Name:":
			device.Name = &fields[1]
		case "Connected:":
			connected := fields[1] == "yes"
			device.Connected = &connected
		}
	}
	if device.Name == nil || device.MacAddr == nil || device.Connected == nil {
		return BluetoothDevice{}, fmt.Errorf("failed to get all expected fields from bluetoothctl output: %q", output)
	}
	return BluetoothDevice{
		Name:      *device.Name,
		MacAddr:   *device.MacAddr,
		Connected: *device.Connected,
	}, nil
}

func (m linuxBluetoothctlBluetoothManager) IsConnected(macAddr string) (bool, error) {
	device, err := m.Get(macAddr)
	if err != nil {
		return false, err
	}
	return device.Connected, nil
}