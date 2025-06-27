//go:build !xreal
// +build !xreal

package xreal

import (
	"fmt"

	"git.terah.dev/UnrealXR/unrealxr/ardriver/commons"
)

var IsXrealEnabled = false

// Implements commons.ARDevice
type XrealDevice struct {
}

func (device *XrealDevice) Initialize() error {
	return fmt.Errorf("xreal is not enabled")
}

func (device *XrealDevice) End() error {
	return fmt.Errorf("xreal is not enabled")
}

func (device *XrealDevice) IsPollingLibrary() bool {
	return false
}

func (device *XrealDevice) IsEventBasedLibrary() bool {
	return false
}

func (device *XrealDevice) Poll() error {
	return fmt.Errorf("xreal is not enabled")
}

func (device *XrealDevice) RegisterEventListeners(*commons.AREventListener) {}

func New() (*XrealDevice, error) {
	return nil, fmt.Errorf("xreal is not enabled")
}
