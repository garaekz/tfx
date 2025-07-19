//go:build darwin

package platform

import "golang.org/x/sys/unix"

// DetectTerminal checks if fd is a terminal on macOS
// Uses TIOCGETA which is available in unix package for BSD-based systems
func DetectTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}

// tryEnableANSI is a no-op on macOS (ANSI is natively supported)
func tryEnableANSI() bool {
	return true
}
