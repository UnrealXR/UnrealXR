package renderer

import (
	"image/color"
	"math"
	"time"
	"unsafe"

	libconfig "git.lunr.sh/UnrealXR/unrealxr/app/config"
	"git.lunr.sh/UnrealXR/unrealxr/app/edidtools"
	"git.lunr.sh/UnrealXR/unrealxr/ardriver"
	arcommons "git.lunr.sh/UnrealXR/unrealxr/ardriver/commons"
	"github.com/charmbracelet/log"
	"github.com/tebeka/atexit"

	rl "git.lunr.sh/UnrealXR/raylib-go/raylib"
)

type TextureModelPair struct {
	Texture               rl.Texture2D
	Model                 rl.Model
	CurrentAngle          float32
	CurrentDisplaySpacing float32
}

func findMaxVerticalSize(fovyDeg float32, distance float32) float32 {
	fovyRad := float64(fovyDeg * math.Pi / 180.0)
	return 2 * distance * float32(math.Tan(fovyRad/2))
}

func findOptimalHorizontalRes(verticalDisplayRes float32, horizontalDisplayRes float32, verticalSize float32) float32 {
	aspectRatio := horizontalDisplayRes / verticalDisplayRes
	horizontalSize := verticalSize * aspectRatio

	return horizontalSize
}

func findHfovFromVfov(vfovDeg, w, h float64) float64 {
	vfovRad := vfovDeg * math.Pi / 180

	ar := w / h
	hfovRad := 2 * math.Atan(math.Tan(vfovRad/2)*ar)

	return hfovRad * 180 / math.Pi
}

func EnterRenderLoop(config *libconfig.Config, displayMetadata *edidtools.DisplayMetadata, evdiCards []*EvdiDisplayMetadata) {
	log.Info("Initializing AR driver")
	headset, err := ardriver.GetDevice()

	if err != nil {
		log.Errorf("Failed to get device: %s", err.Error())
		atexit.Exit(1)
	}

	log.Info("Initialized")

	var (
		currentPitch  float32
		previousPitch float32
		currentYaw    float32
		previousYaw   float32
		currentRoll   float32
		previousRoll  float32

		hasGottenPitchCallbackBefore bool
		hasGottenYawCallbackBefore   bool
		hasGottenRollCallbackBefore  bool
	)

	arEventListner := &arcommons.AREventListener{
		PitchCallback: func(newPitch float32) {
			if !hasGottenPitchCallbackBefore {
				hasGottenPitchCallbackBefore = true
				currentPitch = newPitch
				previousPitch = newPitch
			} else {
				previousPitch = currentPitch
				currentPitch = newPitch
			}
		},
		YawCallback: func(newYaw float32) {
			if !hasGottenYawCallbackBefore {
				hasGottenYawCallbackBefore = true
				currentYaw = newYaw
				previousYaw = newYaw
			} else {
				previousYaw = currentYaw
				currentYaw = newYaw
			}
		},
		RollCallback: func(newRoll float32) {
			if !hasGottenRollCallbackBefore {
				hasGottenRollCallbackBefore = true
				currentRoll = newRoll
				previousRoll = newRoll
			} else {
				previousRoll = currentRoll
				currentRoll = newRoll
			}
		},
	}

	if headset.IsPollingLibrary() {
		log.Error("Connected AR headset requires polling but polling is not implemented in the renderer!")
		atexit.Exit(1)
	}

	headset.RegisterEventListeners(arEventListner)

	fovY := float32(*config.DisplayConfig.FOV)
	fovX := findHfovFromVfov(float64(fovY), float64(displayMetadata.MaxWidth), float64(displayMetadata.MaxHeight))

	verticalSize := findMaxVerticalSize(fovY, 5.0)

	camera := rl.NewCamera3D(
		rl.Vector3{
			X: 0.0,
			Y: verticalSize / 2,
			Z: 5.0,
		},
		rl.Vector3{
			X: 0.0,
			Y: verticalSize / 2,
			Z: 0.0,
		},
		rl.Vector3{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
		},
		fovY,
		rl.CameraPerspective,
	)

	horizontalSize := findOptimalHorizontalRes(float32(displayMetadata.MaxHeight), float32(displayMetadata.MaxWidth), verticalSize)
	coreMesh := rl.GenMeshPlane(horizontalSize, verticalSize, 1, 1)

	var radius float32

	if *config.DisplayConfig.UseCircularSpacing == true {
		radiusX := (horizontalSize / 2) / float32(math.Tan((float64(fovX)*math.Pi/180.0)/2))
		radiusY := (verticalSize / 2) / float32(math.Tan((float64(fovY)*math.Pi/180.0)/2))

		if radiusY > radiusX {
			radius = radiusY
		} else {
			radius = radiusX
		}

		radius *= *config.DisplayConfig.RadiusMultiplier
	}

	movementVector := rl.Vector3{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}

	lookVector := rl.Vector3{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}

	hasZVectorDisabledQuirk := false
	hasSensorInitDelayQuirk := false
	sensorInitStartTime := time.Now()

	if displayMetadata.DeviceQuirks.ZVectorDisabled {
		log.Warn("QUIRK: The Z vector has been disabled for your specific device")
		hasZVectorDisabledQuirk = true
	}

	if displayMetadata.DeviceQuirks.SensorInitDelay != 0 {
		log.Warnf("QUIRK: Waiting %d second(s) before reading sensors", displayMetadata.DeviceQuirks.SensorInitDelay)
		log.Warn("|| MOVEMENT WILL NOT BE OPERATIONAL DURING THIS TIME. ||")
		hasSensorInitDelayQuirk = true
	}

	rects := make([]*TextureModelPair, len(evdiCards))

	displayAngle := float32(*config.DisplayConfig.Angle)
	displaySpacing := *config.DisplayConfig.Spacing + horizontalSize

	highestPossibleAngleOnBothSides := float32((*config.DisplayConfig.Count)-1) * displayAngle
	highestPossibleDisplaySpacingOnBothSides := float32((*config.DisplayConfig.Count)-1) * displaySpacing

	for i, card := range evdiCards {
		currentAngle := (-highestPossibleAngleOnBothSides) + (displayAngle * float32(i+1))
		currentDisplaySpacing := (-highestPossibleDisplaySpacingOnBothSides) + (displaySpacing * float32(i+1))

		log.Debugf("display #%d: currentAngle=%f, currentDisplaySpacing=%f", i, currentAngle, currentDisplaySpacing)

		image := rl.NewImage(card.Buffer.Buffer, int32(displayMetadata.MaxWidth), int32(displayMetadata.MaxHeight), 1, rl.UncompressedR8g8b8a8)

		texture := rl.LoadTextureFromImage(image)
		model := rl.LoadModelFromMesh(coreMesh)

		// spin up/down
		pitchRad := float32(-90 * rl.Deg2rad)
		// spin left/right
		yawRad := currentAngle * rl.Deg2rad

		rotX := rl.MatrixRotateX(pitchRad)
		rotY := rl.MatrixRotateY(yawRad)

		transform := rl.MatrixMultiply(rotX, rotY)
		model.Transform = transform

		rl.SetMaterialTexture(model.Materials, rl.MapAlbedo, texture)

		rects[i] = &TextureModelPair{
			Texture:               texture,
			Model:                 model,
			CurrentAngle:          currentAngle,
			CurrentDisplaySpacing: currentDisplaySpacing,
		}
	}

	eventTimeoutDuration := 0 * time.Millisecond

	for !rl.WindowShouldClose() {
		if !displayMetadata.DeviceQuirks.UsesMouseMovement {
			if hasSensorInitDelayQuirk {
				if time.Since(sensorInitStartTime) > time.Duration(displayMetadata.DeviceQuirks.SensorInitDelay)*time.Second {
					log.Info("Movement is now enabled.")
					hasSensorInitDelayQuirk = false
				}
			} else {
				lookVector.X = (currentYaw - previousYaw) * 6.5
				lookVector.Y = -(currentPitch - previousPitch) * 6.5

				if !hasZVectorDisabledQuirk {
					lookVector.Z = (currentRoll - previousRoll) * 6.5
				}

				// evil look hacks to not randomly break shit
				maxTrustedSize := float64(7)

				if math.Abs(float64(lookVector.X)) > maxTrustedSize {
					log.Errorf("WOAH. Ignoring extreme camera movement for vector X: %.02f", lookVector.X)
					lookVector.X = 0
				}

				if math.Abs(float64(lookVector.Y)) > maxTrustedSize {
					log.Errorf("WOAH. Ignoring extreme camera movement for vector Y: %.02f", lookVector.Y)
					lookVector.Y = 0
				}

				if !hasZVectorDisabledQuirk && math.Abs(float64(lookVector.Z)) > maxTrustedSize {
					log.Errorf("WOAH. Ignoring extreme camera movement for vector Z: %.02f", lookVector.Z)
					lookVector.Z = 0
				}

				rl.UpdateCameraPro(&camera, movementVector, lookVector, 0)
			}
		} else {
			rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode3D(camera)

		for rectPos, rect := range rects {
			card := evdiCards[rectPos]

			ready, err := card.EvdiNode.WaitUntilEventsAreReadyToHandle(eventTimeoutDuration)

			if err != nil {
				log.Errorf("Failed to wait for display events: %s", err.Error())
				continue
			}

			if ready {
				if err := card.EvdiNode.HandleEvents(card.EventContext); err != nil {
					log.Errorf("Failed to handle display events: %s", err.Error())
					continue
				}

				card.EvdiNode.GrabPixels(card.Rect)

				pixels := unsafe.Slice(
					(*color.RGBA)(unsafe.Pointer(&card.Buffer.Buffer[0])),
					len(card.Buffer.Buffer)/4,
				)

				rl.UpdateTexture(rect.Texture, pixels)
				card.EvdiNode.RequestUpdate(card.Buffer)
			}

			worldPos := rl.Vector3{
				X: 0,
				Y: verticalSize / 2,
				Z: 0,
			}

			if *config.DisplayConfig.UseCircularSpacing == true {
				yawRad := float32(rl.Deg2rad * rect.CurrentAngle)

				// WTF?
				posX := float32(math.Sin(float64(yawRad))) * radius
				posZ := -float32(math.Cos(float64(yawRad))) * radius

				worldPos.X = posX
				worldPos.Z = posZ + radius
			} else {
				worldPos.X = rect.CurrentDisplaySpacing
			}

			rl.DrawModelEx(
				rect.Model,
				worldPos,
				// rotate around X to make it vertical
				rl.Vector3{
					X: 0,
					Y: 0,
					Z: 0,
				},
				0,
				rl.Vector3{
					X: 1,
					Y: 1,
					Z: 1,
				},
				rl.White,
			)
		}

		rl.EndMode3D()
		rl.EndDrawing()
	}

	log.Info("Goodbye!")
	rl.CloseWindow()
}
