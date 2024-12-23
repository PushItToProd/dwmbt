package bluetoothctl

import (
	"testing"
)

func TestParseDeviceInfo(t *testing.T) {
	sampleOutput := `Device F8:4E:17:66:E8:55 (public)
        Name: WF-1000XM4
        Alias: WF-1000XM4
        Class: 0x00240404
        Icon: audio-card
        Paired: yes
        Trusted: no
        Blocked: no
        Connected: no
        LegacyPairing: no
        UUID: Vendor specific           (00000000-deca-fade-deca-deafdecacaff)
        UUID: Headset                   (00001108-0000-1000-8000-00805f9b34fb)
        UUID: Audio Sink                (0000110b-0000-1000-8000-00805f9b34fb)
        UUID: A/V Remote Control Target (0000110c-0000-1000-8000-00805f9b34fb)
        UUID: A/V Remote Control        (0000110e-0000-1000-8000-00805f9b34fb)
        UUID: Handsfree                 (0000111e-0000-1000-8000-00805f9b34fb)
        UUID: PnP Information           (00001200-0000-1000-8000-00805f9b34fb)
        UUID: Vendor specific           (764cbf0d-bbcb-438f-a8bb-6b92759d6053)
        UUID: Vendor specific           (81c2e72a-0591-443e-a1ff-05f988593351)
        UUID: Vendor specific           (8901dfa8-5c7e-4d8f-9f0c-c2b70683f5f0)
        UUID: Vendor specific           (931c7e8a-540f-4686-b798-e8df0a2ad9f7)
        UUID: Vendor specific           (956c7b26-d49a-4ba8-b03f-b17d393cb6e2)
        UUID: Vendor specific           (df21fe2c-2515-4fdb-8886-f12c4d67927c)
        UUID: Vendor specific           (f8d1fbe4-7966-4334-8024-ff96c9330e15)
        Modalias: usb:v054Cp0DE1d0201`

	devinfo, err := ParseDeviceInfo([]byte(sampleOutput))
	if err != nil {
		t.Logf("%#v", devinfo)
		t.Error(err)
	}
}
