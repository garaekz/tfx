package platform

// IsTerminal detects if the file descriptor is a terminal
// This is the single entry point for all platform-specific terminal detection
func IsTerminal(fd uintptr) bool {
	return DetectTerminal(fd)
}

// CanEnableANSI returns true if the platform supports ANSI color codes
func CanEnableANSI() bool {
	// Only Windows needs explicit ANSI enabling
	// Unix systems support ANSI natively if they're terminals
	return true
}

// TryEnableANSI attempts to enable ANSI support on platforms that need it
// Returns true if ANSI support is available (either natively or successfully enabled)
func TryEnableANSI() bool {
	// This will be implemented per platform
	return tryEnableANSI()
}
