//go:build fake_edid_patching
// +build fake_edid_patching

package edidtools

import (
	_ "embed"
	"fmt"

	"github.com/charmbracelet/log"
)

//go:embed bin/xreal-air-edid.bin
var edidFirmware []byte

// Attempts to fetch the EDID firmware for any supported XR glasses device
func FetchXRGlassEDID(allowUnsupportedDevices bool) (*DisplayMetadata, error) {
	log.Warn("Not actually fetching EDID firmware in fake patching build -- using embedded firmware")
	parsedEDID, err := ParseEDID(edidFirmware, allowUnsupportedDevices)

	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded EDID firmware: %w", err)
	}

	parsedEDID.DeviceQuirks.ZVectorDisabled = false
	parsedEDID.DeviceQuirks.SensorInitDelay = 0
	parsedEDID.DeviceQuirks.UsesMouseMovement = true

	return parsedEDID, nil
}

// Loads custom firmware for a supported XR glass device
func LoadCustomEDIDFirmware(displayMetadata *DisplayMetadata, edidFirmware []byte) error {
	log.Warn("Not actually patching EDID firmware in fake patching build -- ignoring")
	return nil
}

// Unloads custom firmware for a supported XR glass device
func UnloadCustomEDIDFirmware(displayMetadata *DisplayMetadata) error {
	log.Warn("Not actually unloading EDID firmware in fake patching build -- ignoring")
	return nil
}
