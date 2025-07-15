//go:build darwin
// +build darwin

package platformtools

import "fmt"

// Attempts to do built in privilege escalation to admin
func PrivilegeEscalate(configDir string) error {
	return fmt.Errorf("privilege escalation not implemented on macOS")
}
