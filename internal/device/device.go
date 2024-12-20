package device

import (
	"sync"
)

type Device struct {
	SerialNumber string
	Connected    bool
}

var devices = struct {
	mu      sync.Mutex
	entries map[string]*Device
}{
	entries: make(map[string]*Device),
}

func AddOrUpdateDevice(serialNumber string, connected bool) {
	devices.mu.Lock()
	defer devices.mu.Unlock()

	if device, exists := devices.entries[serialNumber]; exists {
		device.Connected = connected
	} else {
		devices.entries[serialNumber] = &Device{
			SerialNumber: serialNumber,
			Connected:    connected,
		}
	}
}

func getConnectedDevices() []string {
	devices.mu.Lock()
	defer devices.mu.Unlock()

	var connectedDevices []string
	for _, device := range devices.entries {
		if device.Connected {
			connectedDevices = append(connectedDevices, device.SerialNumber)
		}
	}
	return connectedDevices
}

func IsDeviceConnected(serialNumber string) bool {
	devices.mu.Lock()
	defer devices.mu.Unlock()

	if device, exists := devices.entries[serialNumber]; exists && device.Connected {
		return true
	}
	return false
}
