package runfx

import (
	"bytes"
	"io"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
)

// RenderEngine provides double-buffered, flicker-free rendering with terminal awareness
type RenderEngine struct {
	output    io.Writer
	prevBuf   []byte
	mu        sync.Mutex
	detector  *terminal.Detector // Use same detection as writer package
	colorMode color.Mode
}

func NewRenderEngine(output io.Writer) *RenderEngine {
	detector := terminal.NewDetector(output)

	// Use same color mode detection logic as writer.ConsoleWriter
	var colorMode color.Mode
	if !terminal.IsTerminal(output) {
		colorMode = color.ModeNoColor
	} else {
		detectedMode := detector.GetMode()
		// Priority order: TrueColor > 256Color > ANSI > NoColor
		if detectedMode >= 3 {
			colorMode = color.ModeTrueColor
		} else if detectedMode >= 2 {
			colorMode = color.Mode256Color
		} else if detectedMode >= 1 {
			colorMode = color.ModeANSI
		} else {
			colorMode = color.ModeNoColor
		}
	}

	return &RenderEngine{
		output:    output,
		prevBuf:   []byte{},
		detector:  detector,
		colorMode: colorMode,
	}
}

// GetColorMode returns the detected color mode for visuals to use
func (r *RenderEngine) GetColorMode() color.Mode {
	return r.colorMode
}

// SupportsColor returns true if the terminal supports color output
func (r *RenderEngine) SupportsColor() bool {
	return r.colorMode != color.ModeNoColor
}

// Render writes only the diff between current and previous buffer
func (r *RenderEngine) Render(current []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !bytes.Equal(current, r.prevBuf) {
		r.output.Write(current)
		r.prevBuf = make([]byte, len(current))
		copy(r.prevBuf, current)
	}
}

// Clear resets the buffer and clears the screen
func (r *RenderEngine) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.output.Write([]byte("\033[2J\033[H")) // ANSI clear screen
	r.prevBuf = []byte{}
}
