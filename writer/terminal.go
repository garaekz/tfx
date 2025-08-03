package writer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
	"golang.org/x/term"
)

// TerminalWriter handles raw terminal output with double-buffering and color support
// This is used by RunFX for interactive terminal rendering
type TerminalWriter struct {
	output   io.Writer
	detector *terminal.Detector
	prevBuf  []byte
	mu       sync.Mutex
	options  TerminalOptions
}

// TerminalOptions configure the terminal writer behavior
type TerminalOptions struct {
	ForceColor   bool // Force color output even if detection fails
	DisableColor bool // Disable all color output
	DoubleBuffer bool // Enable double-buffering for flicker-free updates
}

// NewTerminalWriter creates a new terminal writer for raw output
func NewTerminalWriter(output io.Writer, opts TerminalOptions) *TerminalWriter {
	return &TerminalWriter{
		output:   output,
		detector: terminal.NewDetector(output),
		prevBuf:  []byte{},
		options:  opts,
	}
}

// Write writes raw bytes to terminal (implements io.Writer)
func (w *TerminalWriter) Write(p []byte) (n int, err error) {
	if w.options.DoubleBuffer {
		return w.writeBuffered(p)
	}
	return w.output.Write(p)
}

// writeBuffered implements double-buffering to prevent flicker
func (w *TerminalWriter) writeBuffered(current []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Only write if content changed
	if !bytes.Equal(current, w.prevBuf) {
		n, err := w.output.Write(current)
		if err != nil {
			return n, err
		}

		// Update previous buffer
		w.prevBuf = make([]byte, len(current))
		copy(w.prevBuf, current)
		return n, nil
	}

	// No change, return success without writing
	return len(current), nil
}

// SupportsColor returns true if terminal supports color output
// Uses the same logic as ConsoleWriter for consistency
func (w *TerminalWriter) SupportsColor() bool {
	if w.options.ForceColor {
		return true
	}
	if w.options.DisableColor {
		return false
	}
	return w.detector.SupportsANSI()
}

// GetColorMode returns the best color mode supported by the terminal
// Uses the same logic as ConsoleWriter for consistency
func (w *TerminalWriter) GetColorMode() color.Mode {
	if w.options.ForceColor {
		return color.ModeTrueColor
	}
	if !w.SupportsColor() {
		return color.ModeNoColor
	}

	// Priority order: TrueColor > 256Color > ANSI > NoColor
	detectedMode := w.detector.GetMode()

	// Try TrueColor first (24-bit)
	if detectedMode >= 3 {
		return color.ModeTrueColor
	}

	// Try 256 color (8-bit)
	if detectedMode >= 2 {
		return color.Mode256Color
	}

	// Fallback to ANSI (4-bit)
	if detectedMode >= 1 {
		return color.ModeANSI
	}

	// No color support
	return color.ModeNoColor
}

// IsTerminal returns true if the output is a terminal
func (w *TerminalWriter) IsTerminal() bool {
	return terminal.IsTerminal(w.output)
}

// Clear clears the terminal screen
func (w *TerminalWriter) Clear() error {
	if !w.IsTerminal() {
		return nil // Don't clear if not a terminal
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// ANSI escape sequence to clear screen and move cursor to top-left
	_, err := w.output.Write([]byte("\033[2J\033[H"))
	if err != nil {
		return err
	}

	// Reset buffer since screen is cleared
	w.prevBuf = []byte{}
	return nil
}

// MoveCursor moves the cursor to the specified position (1-based)
func (w *TerminalWriter) MoveCursor(row, col int) error {
	if !w.IsTerminal() {
		return nil // Don't move cursor if not a terminal
	}

	// ANSI escape sequence to move cursor
	sequence := []byte(fmt.Sprintf("\033[%d;%dH", row, col))
	_, err := w.output.Write(sequence)
	return err
}

// HideCursor hides the terminal cursor
func (w *TerminalWriter) HideCursor() error {
	if !w.IsTerminal() {
		return nil
	}
	_, err := w.output.Write([]byte("\033[?25l"))
	return err
}

// ShowCursor shows the terminal cursor
func (w *TerminalWriter) ShowCursor() error {
	if !w.IsTerminal() {
		return nil
	}
	_, err := w.output.Write([]byte("\033[?25h"))
	return err
}

// GetSize returns the terminal size (width, height)
func (w *TerminalWriter) GetSize() (int, int, error) {
	return terminal.GetSize()
}

// EnableRawMode enables raw terminal mode (for interactive applications)
func (w *TerminalWriter) EnableRawMode() (*term.State, error) {
	if f, ok := w.output.(*os.File); ok {
		return terminal.MakeRaw(f.Fd())
	}
	return nil, fmt.Errorf("cannot enable raw mode: output is not a file")
}

// RestoreMode restores terminal to previous mode
func (w *TerminalWriter) RestoreMode(state *term.State) error {
	if f, ok := w.output.(*os.File); ok {
		return terminal.RestoreTerminal(f.Fd(), state)
	}
	return fmt.Errorf("cannot restore mode: output is not a file")
}

// Flush flushes any buffered content (no-op for terminal, but implements interface)
func (w *TerminalWriter) Flush() error {
	return nil
}

// Close closes the writer (no-op for terminal)
func (w *TerminalWriter) Close() error {
	return nil
}
