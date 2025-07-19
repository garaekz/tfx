//go:build darwin

package terminal

import "golang.org/x/sys/unix"

// isTerminal checks if fd is a terminal on macOS
// Uses TIOCGETA which is available in unix package for BSD-based systems
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}

// enableANSI is a no-op on macOS (ANSI is natively supported)
func enableANSI() bool {
	return true
}
