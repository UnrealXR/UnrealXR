package edidpatcher

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// Calculates a checksum for a given EDID block (base EDID, extension blocks, etc.)
func CalculateEDIDChecksum(edidBlock []byte) byte {
	sum := 0

	for _, value := range edidBlock[:len(edidBlock)-1] {
		sum += int(value)
	}

	return byte((-sum) & 0xFF)
}

var MSFTPayloadSize = byte(22 + 4)

// Patch a given EDID to be a "specialized display", allowing for the display to be used by third-party window-managers/compositors/applications directly.
func PatchEDIDToBeSpecialized(edid []byte) ([]byte, error) {
	newEDID := make([]byte, len(edid))
	copy(newEDID, edid)

	isAnEnhancedEDID := len(newEDID) > 128

	foundExtensionBase := 0
	extensionBaseExists := false

	// Find an appropriate extension base
	if isAnEnhancedEDID {
		for currentExtensionPosition := 128; currentExtensionPosition < len(newEDID); currentExtensionPosition += 128 {
			if newEDID[currentExtensionPosition] != 0x02 {
				continue
			}

			if newEDID[currentExtensionPosition+1] != 0x03 {
				log.Warn("Incompatible version detected for ANSI CTA data section in EDID")
			}

			foundExtensionBase = currentExtensionPosition
			extensionBaseExists = true
		}

		if foundExtensionBase == 0 {
			foundExtensionBase = len(newEDID)
			newEDID = append(newEDID, make([]byte, 128)...)
		}
	} else {
		foundExtensionBase = 128
		newEDID = append(newEDID, make([]byte, 128)...)
	}

	newEDID[foundExtensionBase+2] = MSFTPayloadSize

	if !extensionBaseExists {
		// Add another extension to the original EDID
		if newEDID[126] == 255 {
			return nil, fmt.Errorf("EDID extension block limit reached, but we need to add another extension")
		}

		newEDID[126] += 1
		newEDID[127] = CalculateEDIDChecksum(newEDID[:128])

		newEDID[foundExtensionBase] = 0x02
		newEDID[foundExtensionBase+1] = 0x03
		newEDID[foundExtensionBase+3] = 0x00
	} else {
		if newEDID[foundExtensionBase+2] != MSFTPayloadSize && newEDID[foundExtensionBase+2] != 0 {
			currentBase := newEDID[foundExtensionBase+2]

			copy(newEDID[foundExtensionBase+4:foundExtensionBase+int(currentBase)-1], make([]byte, int(currentBase)-1))
			copy(newEDID[foundExtensionBase+int(MSFTPayloadSize):foundExtensionBase+127], newEDID[foundExtensionBase+int(currentBase):foundExtensionBase+127])
		}
	}

	generatedUUID := uuid.New()
	uuidBytes, err := generatedUUID.MarshalBinary()

	if err != nil {
		return nil, fmt.Errorf("failed to marshal UUID: %w", err)
	}

	// Implemented using https://learn.microsoft.com/en-us/windows-hardware/drivers/display/specialized-monitors-edid-extension
	// VST & Length
	newEDID[foundExtensionBase+4] = 0x3<<5 | 0x15 // 0x3: vendor specific tag; 0x15: length
	// Assigned IEEE OUI
	newEDID[foundExtensionBase+5] = 0x5C
	newEDID[foundExtensionBase+6] = 0x12
	newEDID[foundExtensionBase+7] = 0xCA
	// Actual data
	newEDID[foundExtensionBase+8] = 0x2 // Using version 0x2 for better compatibility
	newEDID[foundExtensionBase+9] = 0x7 // Using VR tag for better compatibility even though it probably doesn't matter
	copy(newEDID[foundExtensionBase+10:foundExtensionBase+10+16], uuidBytes)

	newEDID[foundExtensionBase+127] = CalculateEDIDChecksum(newEDID[foundExtensionBase : foundExtensionBase+127])

	return newEDID, nil
}
