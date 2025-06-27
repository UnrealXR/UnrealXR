//go:build dummy_ar
// +build dummy_ar

package dummy

import (
	"git.terah.dev/UnrealXR/unrealxr/ardriver/commons"
)

var IsDummyDeviceEnabled = true

// Implements commons.ARDevice
type DummyDevice struct {
}

func (device *DummyDevice) Initialize() error {
	return nil
}

func (device *DummyDevice) End() error {
	return nil
}

func (device *DummyDevice) IsPollingLibrary() bool {
	return false
}

func (device *DummyDevice) IsEventBasedLibrary() bool {
	return false
}

func (device *DummyDevice) Poll() error {
	return nil
}

func (device *DummyDevice) RegisterEventListeners(*commons.AREventListener) {}

func New() (*DummyDevice, error) {
	return &DummyDevice{}, nil
}
