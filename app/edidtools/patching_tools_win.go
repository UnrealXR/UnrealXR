//go:build windows && !fake_edid_patching
// +build windows,!fake_edid_patching

package edidtools

import "fmt"

// Attempts to fetch the EDID firmware for any supported XR glasses device
func FetchXRGlassEDID(allowUnsupportedDevices bool) (*DisplayMetadata, error) {
	return nil, fmt.Errorf("automatic fetching of EDID data is not supported on Windows")
}

// Loads custom firmware for a supported XR glass device
func LoadCustomEDIDFirmware(displayMetadata *DisplayMetadata, edidFirmware []byte) error {
	return fmt.Errorf("loading custom EDID firmware is not supported on Windows")
}

// Unloads custom firmware for a supported XR glass device
func UnloadCustomEDIDFirmware(displayMetadata *DisplayMetadata) error {
	return fmt.Errorf("unloading custom EDID firmware is not supported on Windows")
}
