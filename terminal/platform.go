package terminal

import (
	"io"
	"os"

	"golang.org/x/term"
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

// MakeRaw puts the terminal into raw mode.
func MakeRaw(fd uintptr) (*term.State, error) {
	return term.MakeRaw(int(fd))
}

// RestoreTerminal restores the terminal to its original mode.
func RestoreTerminal(fd uintptr, state *term.State) error {
	return term.Restore(int(fd), state)
}

// GetSize returns the terminal size (columns, rows)
func GetSize() (int, int, error) {
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	return width, height, err
}
