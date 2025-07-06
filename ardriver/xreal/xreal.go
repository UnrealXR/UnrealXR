//go:build xreal
// +build xreal

package xreal

import (
	xreal "git.lunr.sh/UnrealXR/unrealxr/ardriver/xreal/xrealsrc"
)

var IsXrealEnabled = true

type XrealDevice struct {
	*xreal.XrealDevice
}

func New() (*XrealDevice, error) {
	device := &XrealDevice{
		XrealDevice: &xreal.XrealDevice{},
	}

	err := device.Initialize()

	if err != nil {
		return nil, err
	}

	return device, nil
}
