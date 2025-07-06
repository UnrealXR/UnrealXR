//go:build !dummy_ar
// +build !dummy_ar

package dummy

import (
	"fmt"

	"git.lunr.sh/UnrealXR/unrealxr/ardriver/commons"
)

var IsDummyDeviceEnabled = false

// Implements commons.ARDevice
type DummyDevice struct {
}

func (device *DummyDevice) Initialize() error {
	return fmt.Errorf("dummy device is not enabled")
}

func (device *DummyDevice) End() error {
	return fmt.Errorf("dummy device is not enabled")
}

func (device *DummyDevice) IsPollingLibrary() bool {
	return false
}

func (device *DummyDevice) IsEventBasedLibrary() bool {
	return false
}

func (device *DummyDevice) Poll() error {
	return fmt.Errorf("dummy device is not enabled")
}

func (device *DummyDevice) RegisterEventListeners(*commons.AREventListener) {}

func New() (*DummyDevice, error) {
	return nil, fmt.Errorf("dummy device is not enabled")
}
