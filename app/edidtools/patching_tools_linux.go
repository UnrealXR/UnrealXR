//go:build linux && !fake_edid_patching
// +build linux,!fake_edid_patching

package edidtools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// Attempts to fetch the EDID firmware for any supported XR glasses device
func FetchXRGlassEDID(allowUnsupportedDevices bool) (*DisplayMetadata, error) {
	// Implementation goes here
	pciDeviceCommand, err := exec.Command("lspci").Output()

	if err != nil {
		return nil, fmt.Errorf("failed to execute lspci command: %w", err)
	}

	pciDevices := strings.Split(string(pciDeviceCommand), "\n")
	pciDevices = pciDevices[:len(pciDevices)-1]

	vgaDevices := []string{}

	for _, pciDevice := range pciDevices {
		if strings.Contains(pciDevice, "VGA compatible controller:") {
			vgaDevices = append(vgaDevices, pciDevice[:strings.Index(pciDevice, " ")])
		}
	}

	for _, vgaDevice := range vgaDevices {
		cardDevices, err := os.ReadDir("/sys/devices/pci0000:00/0000:" + vgaDevice + "/drm/")

		if err != nil {
			return nil, fmt.Errorf("failed to read directory for device '%s': %w", vgaDevice, err)
		}

		for _, cardDevice := range cardDevices {
			if !strings.Contains(cardDevice.Name(), "card") {
				continue
			}

			monitors, err := os.ReadDir("/sys/devices/pci0000:00/0000:" + vgaDevice + "/drm/" + cardDevice.Name())

			if err != nil {
				return nil, fmt.Errorf("failed to read directory for card device '%s': %w", cardDevice.Name(), err)
			}

			for _, monitor := range monitors {
				if !strings.Contains(monitor.Name(), cardDevice.Name()) {
					continue
				}

				rawEDIDFile, err := os.ReadFile("/sys/devices/pci0000:00/0000:" + vgaDevice + "/drm/" + cardDevice.Name() + "/" + monitor.Name() + "/edid")

				if err != nil {
					return nil, fmt.Errorf("failed to read EDID file for monitor '%s': %w", monitor.Name(), err)
				}

				if len(rawEDIDFile) == 0 {
					continue
				}

				parsedEDID, err := ParseEDID(rawEDIDFile, allowUnsupportedDevices)

				if err != nil {
					if !strings.HasPrefix(err.Error(), "failed to match manufacturer for monitor vendor") {
						log.Warnf("Failed to parse EDID for monitor '%s': %s", monitor.Name(), err.Error())
					}
				} else {
					parsedEDID.LinuxDRMCard = cardDevice.Name()
					parsedEDID.LinuxDRMConnector = strings.Replace(monitor.Name(), cardDevice.Name()+"-", "", 1)

					return parsedEDID, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("could not find supported device! Check if the XR device is plugged in. If it is plugged in and working correctly, check the README or open an issue.")
}

// Loads custom firmware for a supported XR glass device
func LoadCustomEDIDFirmware(displayMetadata *DisplayMetadata, edidFirmware []byte) error {
	if displayMetadata.LinuxDRMCard == "" || displayMetadata.LinuxDRMConnector == "" {
		return fmt.Errorf("missing Linux DRM card or connector information")
	}

	drmFile, err := os.OpenFile("/sys/kernel/debug/dri/"+strings.Replace(displayMetadata.LinuxDRMCard, "card", "", 1)+"/"+displayMetadata.LinuxDRMConnector+"/edid_override", os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open EDID override file for monitor '%s': %w", displayMetadata.LinuxDRMConnector, err)
	}

	defer drmFile.Close()

	if _, err := drmFile.Write(edidFirmware); err != nil {
		return fmt.Errorf("failed to write EDID firmware for monitor '%s': %w", displayMetadata.LinuxDRMConnector, err)
	}

	return nil
}

// Unloads custom firmware for a supported XR glass device
func UnloadCustomEDIDFirmware(displayMetadata *DisplayMetadata) error {
	if displayMetadata.LinuxDRMCard == "" || displayMetadata.LinuxDRMConnector == "" {
		return fmt.Errorf("missing Linux DRM card or connector information")
	}

	drmFile, err := os.OpenFile("/sys/kernel/debug/dri/"+strings.Replace(displayMetadata.LinuxDRMCard, "card", "", 1)+"/"+displayMetadata.LinuxDRMConnector+"/edid_override", os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed to open EDID override file for monitor '%s': %w", displayMetadata.LinuxDRMConnector, err)
	}

	defer drmFile.Close()

	if _, err := drmFile.Write([]byte("reset")); err != nil {
		return fmt.Errorf("failed to unload EDID firmware for monitor '%s': %w", displayMetadata.LinuxDRMConnector, err)
	}

	return nil
}
