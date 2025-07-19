//go:build windows

package platform

import (
	"os"
	"syscall"
	"unsafe"
)

// Windows-specific syscalls
var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

// DetectTerminal checks if fd is a terminal on Windows
func DetectTerminal(fd uintptr) bool {
	var st uint32
	r, _, e := syscall.SyscallN(procGetConsoleMode.Addr(), fd, uintptr(unsafe.Pointer(&st)))
	return r != 0 && e == 0
}

// tryEnableANSI enables ANSI escape sequence support on Windows 10+
func tryEnableANSI() bool {
	stdout := os.Stdout.Fd()
	var mode uint32

	// Get current console mode
	r, _, _ := syscall.SyscallN(procGetConsoleMode.Addr(), stdout, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return false
	}

	// Enable virtual terminal processing (ANSI support)
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING

	// Set the new mode
	r, _, _ = syscall.SyscallN(procSetConsoleMode.Addr(), stdout, uintptr(mode))
	return r != 0
}
