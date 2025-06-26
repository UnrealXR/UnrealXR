package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	libconfig "git.terah.dev/UnrealXR/unrealxr/app/config"
	"git.terah.dev/UnrealXR/unrealxr/app/edidtools"
	"git.terah.dev/UnrealXR/unrealxr/app/renderer"
	"git.terah.dev/UnrealXR/unrealxr/edidpatcher"
	"git.terah.dev/imterah/goevdi/libevdi"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/kirsle/configdir"
	"github.com/tebeka/atexit"
	"github.com/urfave/cli/v3"

	rl "git.terah.dev/UnrealXR/raylib-go/raylib"
)

func mainEntrypoint(context.Context, *cli.Command) error {
	// Allow for clean exits
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info("Exiting...")
		atexit.Exit(1)
	}()

	// TODO: add built-in privesc
	if os.Geteuid() != 0 {
		return fmt.Errorf("this program must be run as root")
	}

	log.Info("Initializing UnrealXR")

	// Allow for overriding the config directory
	configDir := os.Getenv("UNREALXR_CONFIG_PATH")

	if configDir == "" {
		configDir = configdir.LocalConfig("unrealxr")
		err := configdir.MakePath(configDir)

		if err != nil {
			return fmt.Errorf("failed to ensure config directory exists: %w", err)
		}
	}

	_, err := os.Stat(path.Join(configDir, "config.yml"))

	if err != nil {
		log.Debug("Creating default config file")
		err := os.WriteFile(path.Join(configDir, "config.yml"), libconfig.InitialConfig, 0644)

		if err != nil {
			return fmt.Errorf("failed to create initial config file: %w", err)
		}
	}

	// Read and parse the config file
	configBytes, err := os.ReadFile(path.Join(configDir, "config.yml"))

	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &libconfig.Config{}
	err = yaml.Unmarshal(configBytes, config)

	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	libconfig.InitializePotentiallyMissingConfigValues(config)
	log.Debug("Attempting to read display EDID file and fetch metadata")

	displayMetadata, err := edidtools.FetchXRGlassEDID(*config.Overrides.AllowUnsupportedDevices)

	if err != nil {
		return fmt.Errorf("failed to fetch EDID or get metadata: %w", err)
	}

	log.Debug("Got EDID file and metadata")
	log.Debug("Patching EDID firmware to be specialized")

	patchedFirmware, err := edidpatcher.PatchEDIDToBeSpecialized(displayMetadata.EDID)

	if err != nil {
		return fmt.Errorf("failed to patch EDID firmware: %w", err)
	}

	log.Info("Uploading patched EDID firmware")
	err = edidtools.LoadCustomEDIDFirmware(displayMetadata, patchedFirmware)

	if err != nil {
		return fmt.Errorf("failed to upload patched EDID firmware: %w", err)
	}

	atexit.Register(func() {
		err := edidtools.UnloadCustomEDIDFirmware(displayMetadata)

		if err != nil {
			log.Errorf("Failed to unload custom EDID firmware: %s", err.Error())
		}

		log.Info("Please unplug and plug in your XR device to restore it back to normal settings.")
	})

	fmt.Print("Press the Enter key to continue loading after you unplug and plug in your XR device.")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // Wait for Enter key press before continuing

	log.Info("Initializing XR headset")
	rl.SetTargetFPS(int32(displayMetadata.MaxRefreshRate * 2))
	rl.InitWindow(int32(displayMetadata.MaxWidth), int32(displayMetadata.MaxHeight), "UnrealXR")

	atexit.Register(func() {
		rl.CloseWindow()
	})

	log.Info("Initializing virtual displays")

	libevdi.SetupLogger(&libevdi.EvdiLogger{
		Log: func(msg string) {
			log.Debugf("EVDI: %s", msg)
		},
	})

	evdiCards := make([]*renderer.EvdiDisplayMetadata, *config.DisplayConfig.Count)

	for currentDisplay := range *config.DisplayConfig.Count {
		openedDevice, err := libevdi.Open(nil)

		if err != nil {
			log.Errorf("Failed to open EVDI device: %s", err.Error())
		}

		openedDevice.Connect(displayMetadata.EDID, uint(displayMetadata.MaxWidth), uint(displayMetadata.MaxHeight), uint(displayMetadata.MaxRefreshRate))

		atexit.Register(func() {
			openedDevice.Disconnect()
		})

		displayRect := &libevdi.EvdiDisplayRect{
			X1: 0,
			Y1: 0,
			X2: displayMetadata.MaxWidth,
			Y2: displayMetadata.MaxHeight,
		}

		displayBuffer := openedDevice.CreateBuffer(displayMetadata.MaxWidth, displayMetadata.MaxHeight, libevdi.StridePixelFormatRGBA32, displayRect)

		displayMetadata := &renderer.EvdiDisplayMetadata{
			EvdiNode: openedDevice,
			Rect:     displayRect,
			Buffer:   displayBuffer,
		}

		displayMetadata.EventContext = &libevdi.EvdiEventContext{}
		openedDevice.RegisterEventHandler(displayMetadata.EventContext)

		evdiCards[currentDisplay] = displayMetadata
	}

	// HACK: sometimes the buffer doesn't get initialized properly if we don't wait a bit...
	time.Sleep(time.Millisecond * 100)

	log.Info("Initialized displays. Entering rendering loop")
	renderer.EnterRenderLoop(config, displayMetadata, evdiCards)

	atexit.Exit(0)
	return nil
}

func main() {
	logLevel := os.Getenv("UNREALXR_LOG_LEVEL")

	if logLevel != "" {
		switch logLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)

		case "info":
			log.SetLevel(log.InfoLevel)

		case "warn":
			log.SetLevel(log.WarnLevel)

		case "error":
			log.SetLevel(log.ErrorLevel)

		case "fatal":
			log.SetLevel(log.FatalLevel)
		}
	}

	// Initialize the CLI
	cmd := &cli.Command{
		Name:   "unrealxr",
		Usage:  "A spatial multi-display renderer for XR devices",
		Action: mainEntrypoint,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Fatal error during execution: %s", err.Error())
	}
}
