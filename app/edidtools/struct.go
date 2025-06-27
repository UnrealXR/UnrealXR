package edidtools

type DisplayQuirks struct {
	MaxWidth          int
	MaxHeight         int
	MaxRefreshRate    int
	SensorInitDelay   int
	ZVectorDisabled   bool
	UsesMouseMovement bool
}

type DisplayMetadata struct {
	EDID              []byte
	DeviceVendor      string
	DeviceQuirks      DisplayQuirks
	MaxWidth          int
	MaxHeight         int
	MaxRefreshRate    int
	LinuxDRMCard      string
	LinuxDRMConnector string
}
