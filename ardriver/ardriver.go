package ardriver

import (
	"fmt"

	"git.terah.dev/UnrealXR/unrealxr/ardriver/commons"
	"git.terah.dev/UnrealXR/unrealxr/ardriver/dummy"
	"git.terah.dev/UnrealXR/unrealxr/ardriver/xreal"
)

func GetDevice() (commons.ARDevice, error) {
	if xreal.IsXrealEnabled {
		device, err := xreal.New()

		if err != nil {
			fmt.Printf("failed to initialize xreal device: %w\n", err)
			return nil, err
		}

		return device, nil
	}

	if dummy.IsDummyDeviceEnabled {
		device, err := dummy.New()

		if err != nil {
			fmt.Printf("failed to initialize dummy device: %w\n", err)
			return nil, err
		}

		return device, nil
	}

	return nil, fmt.Errorf("failed to initialize any device")
}
