//go:build linux || aix

package terminal

import (
	"runtime"

	"golang.org/x/sys/unix"
)

// Terminal ioctl constants for Linux-based systems
const (
	// Embedded Linux variants (only define if not in unix package)
	TCGETS_EMBEDDED = 0x5400 // Some embedded systems (routers, IoT)
)

// isTerminal checks if fd is a terminal on Linux-based systems
func isTerminal(fd uintptr) bool {
	// Try unix package constant first (most reliable)
	if _, err := unix.IoctlGetTermios(int(fd), unix.TCGETS); err == nil {
		return true
	}

	// Embedded Linux fallback for ARM/MIPS devices
	if isEmbeddedArch() {
		if _, err := unix.IoctlGetTermios(int(fd), TCGETS_EMBEDDED); err == nil {
			return true
		}
	}

	return false
}

// isEmbeddedArch checks if we're running on embedded architecture
func isEmbeddedArch() bool {
	return runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" ||
		runtime.GOARCH == "mips" || runtime.GOARCH == "mipsle" ||
		runtime.GOARCH == "mips64" || runtime.GOARCH == "mips64le"
}

// enableANSI is a no-op on Linux systems (ANSI is natively supported)
func enableANSI() bool {
	return true
}
