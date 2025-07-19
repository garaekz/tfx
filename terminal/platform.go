package terminal

import (
	"io"
	"os"
)

// IsTerminal detects if the file descriptor is a terminal
// This is the single entry point for all platform-specific terminal detection
func IsTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return isTerminal(v.Fd())
	default:
		return false
	}
}

// TryEnableANSI attempts to enable ANSI support on platforms that need it
// Returns true if ANSI support is available (either natively or successfully enabled)
func TryEnableANSI() bool {
	return enableANSI()
}
