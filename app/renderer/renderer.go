package renderer

import (
	"time"

	libconfig "git.terah.dev/UnrealXR/unrealxr/app/config"
	"git.terah.dev/UnrealXR/unrealxr/app/edidtools"
	"git.terah.dev/UnrealXR/unrealxr/ardriver"
	arcommons "git.terah.dev/UnrealXR/unrealxr/ardriver/commons"
	"github.com/charmbracelet/log"
	"github.com/tebeka/atexit"
)

func EnterRenderLoop(config *libconfig.Config, displayMetadata *edidtools.DisplayMetadata, evdiCards []*EvdiDisplayMetadata) {
	log.Info("Initializing AR driver")
	headset, err := ardriver.GetDevice()

	if err != nil {
		log.Errorf("Failed to get device: %s", err.Error())
		atexit.Exit(1)
	}
	log.Info("Initialized")

	var pitch float32
	var yaw float32
	var roll float32

	arEventListner := &arcommons.AREventListener{
		PitchCallback: func(newPitch float32) {
			pitch = newPitch
		},
		YawCallback: func(newYaw float32) {
			yaw = newYaw
		},
		RollCallback: func(newRoll float32) {
			roll = newRoll
		},
	}

	if headset.IsPollingLibrary() {
		log.Error("Connected AR headset requires polling but polling is not implemented in the renderer!")
		atexit.Exit(1)
	}

	headset.RegisterEventListeners(arEventListner)

	for {
		log.Debugf("pitch: %f, yaw: %f, roll: %f", pitch, yaw, roll)
		time.Sleep(10 * time.Millisecond)
	}
}
