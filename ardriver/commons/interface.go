package commons

type AREventListener struct {
	PitchCallback func(float32)
	YawCallback   func(float32)
	RollCallback  func(float32)
}

type ARDevice interface {
	// Initializes the AR device's sensors.
	Initialize() error
	// Ends the AR device's sensors.
	End() error
	// Polls the AR device's sensors.
	Poll() error
	// Checks if the underlying AR library is polling-based.
	IsPollingLibrary() bool
	// Checks if the underlying AR library is event-based.
	IsEventBasedLibrary() bool
	// Registers event listeners for the AR device.
	RegisterEventListeners(eventListener *AREventListener)
}
