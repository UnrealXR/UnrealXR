package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/kirsle/configdir"
	"github.com/urfave/cli/v3"
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
		err := os.WriteFile(path.Join(configDir, "config.yml"), InitialConfig, 0644)

		if err != nil {
			return fmt.Errorf("failed to create initial config file: %w", err)
		}
	}

	// Read and parse the config file
	configBytes, err := os.ReadFile(path.Join(configDir, "config.yml"))

	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(configBytes, config)

	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	InitializePotentiallyMissingConfigValues(config)
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
		log.Fatal(err)
	}
}
