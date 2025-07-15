//go:build linux
// +build linux

package platformtools

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

// Checks if we're in a Nix shell (all Linux OSes), or if we're in a NixOS system
func isNixLikeEnvironment() bool {
	if os.Getenv("IN_NIX_SHELL") != "" {
		return true
	}

	_, err := os.Stat("/etc/nixos")
	return err == nil
}

// Attempts to do built in privilege escalation to admin
func PrivilegeEscalate(configDir string) error {
	executablePath, err := os.Executable()

	if err != nil {
		return fmt.Errorf("could not find own executable path")
	}

	logLevel := os.Getenv("UNREALXR_LOG_LEVEL")
	systemPath := os.Getenv("PATH")
	libraryPath := os.Getenv("LD_LIBRARY_PATH")
	waylandDisplay := path.Join(os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY"))
	rootXDGRuntimeDir := "/run/user/0"

	var command *exec.Cmd

	if isNixLikeEnvironment() {
		command = exec.Command(
			"pkexec",
			"--keep-cwd",
			"/usr/bin/env",
			"UXR_HAS_PRIVESC=1",
			"UNREALXR_LOG_LEVEL="+logLevel,
			"UNREALXR_CONFIG_PATH="+configDir,
			"WAYLAND_DISPLAY="+waylandDisplay,
			"XDG_RUNTIME_DIR="+rootXDGRuntimeDir,
			"LD_LIBRARY_PATH="+libraryPath,
			"PATH="+systemPath,
			executablePath,
		)
	} else {
		command = exec.Command(
			"pkexec",
			"--keep-cwd",
			"/usr/bin/env",
			"UXR_HAS_PRIVESC=1",
			"UNREALXR_LOG_LEVEL="+logLevel,
			"UNREALXR_CONFIG_PATH="+configDir,
			"WAYLAND_DISPLAY="+waylandDisplay,
			"XDG_RUNTIME_DIR="+rootXDGRuntimeDir,
			executablePath,
		)
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)

		if !ok {
			return fmt.Errorf("failed to execute command, and failed to typecast err to ExitError")
		}

		os.Exit(exitErr.ExitCode())
	}

	return nil
}
