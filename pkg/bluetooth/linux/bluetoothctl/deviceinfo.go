package bluetoothctl

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// DeviceInfo represents a subset of data parsed from the output of `bluetoothctl info <macAddr>`.
type DeviceInfo struct {
	Name      *string
	MacAddr   *string
	Connected *bool
}

// Validate checks that all required fields are present, just in case the returned data somehow fails to match our
// expectations.
func (i DeviceInfo) Validate() error {
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

// ParseDeviceInfo parses the output of `bluetoothctl info <macAddr>`.
func ParseDeviceInfo(output []byte) (DeviceInfo, error) {
	var err error
	var device DeviceInfo
	// bluetoothctl doesn't provide structured output like JSON or whatever, so we have to parse it manually. We only
	// care about a few fields from the output, so we can just look for those lines and extract the values.
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

		// As soon as we have all the fields we care about, we can stop scanning.
		err = device.Validate()
		if err == nil {
			return device, nil
		}
	}
	return device, err
}
