package edidtools

// Vendor and devices names sourced from "https://uefi.org/uefi-pnp-export"
var QuirksRegistry = map[string]map[string]DisplayQuirks{
	"MRG": {
		"Air": {
			MaxWidth:        1920,
			MaxHeight:       1080,
			MaxRefreshRate:  120,
			SensorInitDelay: 10,
			ZVectorDisabled: true,
		},
	},
}
