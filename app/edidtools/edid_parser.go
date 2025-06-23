package edidtools

import (
	"fmt"

	edidparser "github.com/anoopengineer/edidparser/edid"
)

func ParseEDID(rawEDIDFile []byte, allowUnsupportedDevices bool) (*DisplayMetadata, error) {
	parsedEDID, err := edidparser.NewEdid(rawEDIDFile)

	if err != nil {
		return nil, fmt.Errorf("failed to parse EDID file: %w", err)
	}

	for manufacturer, manufacturerSupportedDevices := range QuirksRegistry {
		if parsedEDID.ManufacturerId == manufacturer {
			if deviceQuirks, ok := manufacturerSupportedDevices[parsedEDID.MonitorName]; ok || allowUnsupportedDevices {
				maxWidth := 0
				maxHeight := 0
				maxRefreshRate := 0

				for _, resolution := range parsedEDID.DetailedTimingDescriptors {
					if int(resolution.HorizontalActive) > maxWidth && int(resolution.VerticalActive) > maxHeight {
						maxWidth = int(resolution.HorizontalActive)
						maxHeight = int(resolution.VerticalActive)
					}

					// Convert pixel clock to refresh rate
					// Refresh Rate = Pixel Clock / ((Horizontal Active + Horizontal Blanking) * (Vertical Active + Vertical Blanking))
					hTotal := int(resolution.HorizontalActive + resolution.HorizontalBlanking)
					vTotal := int(resolution.VerticalActive + resolution.VerticalBlanking)
					refreshRate := int(int(resolution.PixelClock*1000) / (hTotal * vTotal))

					if refreshRate > maxRefreshRate {
						maxRefreshRate = refreshRate
					}
				}

				if maxWidth == 0 || maxHeight == 0 {
					if deviceQuirks.MaxWidth == 0 || deviceQuirks.MaxHeight == 0 {
						return nil, fmt.Errorf("failed to determine maximum resolution for monitor '%s'", parsedEDID.MonitorName)
					}

					maxWidth = deviceQuirks.MaxWidth
					maxHeight = deviceQuirks.MaxHeight
				}

				if maxRefreshRate == 0 {
					if deviceQuirks.MaxRefreshRate == 0 {
						return nil, fmt.Errorf("failed to determine maximum refresh rate for monitor '%s'", parsedEDID.MonitorName)
					}

					maxRefreshRate = deviceQuirks.MaxRefreshRate
				}

				displayMetadata := &DisplayMetadata{
					EDID:           rawEDIDFile,
					DeviceVendor:   parsedEDID.ManufacturerId,
					DeviceQuirks:   deviceQuirks,
					MaxWidth:       maxWidth,
					MaxHeight:      maxHeight,
					MaxRefreshRate: maxRefreshRate,
				}

				return displayMetadata, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to match manufacturer for monitor vendor: '%s'", parsedEDID.ManufacturerId)
}
