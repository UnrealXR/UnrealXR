package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"

	libconfig "git.terah.dev/UnrealXR/unrealxr/app/config"
	"git.terah.dev/UnrealXR/unrealxr/app/edidtools"
	"git.terah.dev/UnrealXR/unrealxr/edidpatcher"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/kirsle/configdir"
	"github.com/urfave/cli/v3"

	rl "git.terah.dev/UnrealXR/raylib-go/raylib"
)

func mainEntrypoint(context.Context, *cli.Command) error {
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
	log.Info("Attempting to read display EDID file and fetch metadata")

	displayMetadata, err := edidtools.FetchXRGlassEDID(*config.Overrides.AllowUnsupportedDevices)

	if err != nil {
		return fmt.Errorf("failed to fetch EDID or get metadata: %w", err)
	}

	log.Info("Got EDID file and metadata")
	log.Info("Patching EDID firmware to be specialized")

	patchedFirmware, err := edidpatcher.PatchEDIDToBeSpecialized(displayMetadata.EDID)

	if err != nil {
		return fmt.Errorf("failed to patch EDID firmware: %w", err)
	}

	log.Info("Uploading patched EDID firmware")
	err = edidtools.LoadCustomEDIDFirmware(displayMetadata, patchedFirmware)

	if err != nil {
		return fmt.Errorf("failed to upload patched EDID firmware: %w", err)
	}

	defer func() {
		err := edidtools.UnloadCustomEDIDFirmware(displayMetadata)

		if err != nil {
			log.Errorf("Failed to unload custom EDID firmware: %s", err.Error())
		}

		log.Info("Please unplug and plug in your XR device to restore it back to normal settings.")
	}()

	fmt.Print("Press the Enter key to continue loading after you unplug and plug in your XR device.")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // Wait for Enter key press before continuing

	log.Info("Initializing XR headset")

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

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
