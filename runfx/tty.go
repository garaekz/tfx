package runfx

import (
	"fmt"
	"io"
	"os"

	"github.com/garaekz/tfx/writer"
)

// TTYInfo holds terminal capability details (uses writer.TerminalWriter for consistency)
type TTYInfo struct {
	IsTTY     bool
	TrueColor bool
	ANSI      bool
	NoColor   bool
}

// DetectTTY returns TTYInfo using the same detection logic as writer.TerminalWriter
func DetectTTY() TTYInfo {
	return DetectTTYForOutput(os.Stdout)
}

// DetectTTYForOutput returns TTYInfo for a specific output writer
func DetectTTYForOutput(output io.Writer) TTYInfo {
	if output == nil {
		output = os.Stdout
	}

	// Use TerminalWriter for consistent detection logic
	termWriter := writer.NewTerminalWriter(output, writer.TerminalOptions{})

	isTTY := termWriter.IsTerminal()
	supportsColor := termWriter.SupportsColor()
	colorMode := termWriter.GetColorMode()
	ansi := isTTY && supportsColor

	return TTYInfo{
		IsTTY:     isTTY,
		TrueColor: colorMode.String() == "TrueColor",
		ANSI:      ansi,
		NoColor:   !supportsColor,
	}
}

// FallbackOutput prints minimal output if not TTY
func FallbackOutput(msg string) {
	if !DetectTTY().IsTTY {
		fmt.Fprintln(os.Stdout, msg)
	}
}
