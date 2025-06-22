package main

import _ "embed"

//go:embed default_config.yml
var InitialConfig []byte

type DisplayConfig struct {
	Angle   *int `yaml:"angle"`
	FOV     *int `yaml:"fov"`
	Spacing *int `yaml:"spacing"`
	Count   *int `yaml:"count"`
}

type AppOverrides struct {
	AllowUnsupportedDevices *bool `yaml:"allow_unsupported_devices"`
	OverrideWidth           *int  `yaml:"width"`
	OverrideHeight          *int  `yaml:"height"`
	OverrideRefreshRate     *int  `yaml:"refresh_rate"`
}

type Config struct {
	DisplayConfig DisplayConfig `yaml:"display"`
	Overrides     AppOverrides  `yaml:"overrides"`
}

func getPtrToInt(int int) *int {
	return &int
}

func getPtrToBool(bool bool) *bool {
	return &bool
}

var DefaultConfig = &Config{
	DisplayConfig: DisplayConfig{
		Angle:   getPtrToInt(45),
		FOV:     getPtrToInt(45),
		Spacing: getPtrToInt(1),
		Count:   getPtrToInt(3),
	},
	Overrides: AppOverrides{
		AllowUnsupportedDevices: getPtrToBool(false),
	},
}

func InitializePotentiallyMissingConfigValues(config *Config) {
	// TODO: is there a better way to do this?
	if config.DisplayConfig.Angle == nil {
		config.DisplayConfig.Angle = DefaultConfig.DisplayConfig.Angle
	}

	if config.DisplayConfig.FOV == nil {
		config.DisplayConfig.FOV = DefaultConfig.DisplayConfig.FOV
	}

	if config.DisplayConfig.Spacing == nil {
		config.DisplayConfig.Spacing = DefaultConfig.DisplayConfig.Spacing
	}

	if config.DisplayConfig.Count == nil {
		config.DisplayConfig.Count = DefaultConfig.DisplayConfig.Count
	}

	if config.Overrides.AllowUnsupportedDevices == nil {
		config.Overrides.AllowUnsupportedDevices = DefaultConfig.Overrides.AllowUnsupportedDevices
	}

	if config.Overrides.OverrideWidth == nil {
		config.Overrides.OverrideWidth = DefaultConfig.Overrides.OverrideWidth
	}

	if config.Overrides.OverrideHeight == nil {
		config.Overrides.OverrideHeight = DefaultConfig.Overrides.OverrideHeight
	}

	if config.Overrides.OverrideRefreshRate == nil {
		config.Overrides.OverrideRefreshRate = DefaultConfig.Overrides.OverrideRefreshRate
	}
}
