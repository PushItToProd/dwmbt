package bluetooth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
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
	switch os := runtime.GOOS; os {
	case "darwin":
		return newMacosBlueutilBluetoothManager()
	case "linux":
		return newLinuxBluetoothctlBluetoothManager()
	default:
		panic(fmt.Sprintf("unsupported OS: %s", os))
	}
}

type macosBlueutilBluetoothManager struct{}

func newMacosBlueutilBluetoothManager() macosBlueutilBluetoothManager {
	// validate blueutil is on the path
	_, err := exec.LookPath("blueutil")
	if err != nil {
		panic("couldn't find blueutil on the path")
	}

	return macosBlueutilBluetoothManager{}
}

func (m macosBlueutilBluetoothManager) Connect(macAddr string) error {
	cmd := exec.Command("blueutil", "--connect", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

func (m macosBlueutilBluetoothManager) Disconnect(macAddr string) error {
	cmd := exec.Command("blueutil", "--disconnect", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return err
	}
	// TODO: validate output
	_ = output
	return nil
}

type blueutilDeviceInfo struct {
	Address   string `json:"address"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	Paired    bool   `json:"paired"`
}

func (m macosBlueutilBluetoothManager) List() ([]BluetoothDevice, error) {
	cmd := exec.Command("blueutil", "--paired", "--format", "json")
	output, err := cmd.Output()
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

func (m macosBlueutilBluetoothManager) Get(macAddr string) (BluetoothDevice, error) {
	var device BluetoothDevice
	cmd := exec.Command("blueutil", "--info", macAddr, "--format", "json")
	output, err := cmd.Output()
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

func (m macosBlueutilBluetoothManager) IsConnected(macAddr string) (bool, error) {
	cmd := exec.Command("blueutil", "--is-connected", macAddr)
	output, err := cmd.Output()
	if err != nil {
		// TODO: process the error to see what went wrong
		return false, err
	}

	// --is-connected returns '1' if the device is connected and 0 if not
	return output[0] == '1', nil
}

type linuxBluetoothctlBluetoothManager struct{}

func newLinuxBluetoothctlBluetoothManager() linuxBluetoothctlBluetoothManager {
	// validate bluetoothctl is on the path
	_, err := exec.LookPath("bluetoothctl")
	if err != nil {
		panic("couldn't find blueutil on the path")
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
	devices := make([]BluetoothDevice, len(lines))
	for _, line := range lines {
		ms := bluetoothctlDeviceListRegex.FindSubmatch(line)
		if ms == nil {
			log.Printf("failed to match bluetoothctl device line: %s", string(line))
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

type bluetoothctlDevice struct {
	Name      *string
	MacAddr   *string
	Connected *bool
}

func (m linuxBluetoothctlBluetoothManager) Get(macAddr string) (BluetoothDevice, error) {
	var device bluetoothctlDevice
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
