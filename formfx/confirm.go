package formfx

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
	"github.com/garaekz/tfx/terminal"
)

// ConfirmConfig provides configuration for the Confirm prompt.
type ConfirmConfig struct {
	// Label is the prompt message displayed to the user.
	Label string
	// Default is the default value if the user just presses Enter.
	Default bool
	// Writer is the output writer for the prompt.
	Writer io.Writer
	// Reader is the input reader for the prompt.
	Reader Reader
	// Interactive enables RunFX-powered interactive mode with visual feedback.
	Interactive bool
}

// DefaultConfirmConfig returns the default configuration for Confirm.
func DefaultConfirmConfig() ConfirmConfig {
	// Detect if we're in an interactive environment
	ttyInfo := runfx.DetectTTY()

	return ConfirmConfig{
		Label:       "Are you sure?",
		Default:     false,
		Writer:      os.Stdout,
		Reader:      NewStdinReader(os.Stdin),
		Interactive: ttyInfo.IsTTY, // Enable interactive mode if TTY available
	}
}

// --- MULTIPATH API FUNCTIONS ---

// Confirm prompts the user for a yes/no confirmation with multipath configuration support.
// Supports two usage patterns:
//   - Confirm(label)                         // Express: simple label
//   - Confirm(config)                        // Instantiated: config struct
func Confirm(args ...any) (bool, error) {
	// Handle different argument patterns
	if len(args) == 0 {
		// No args: use default config
		cfg := DefaultConfirmConfig()
		return ConfirmWithConfig(cfg)
	}

	// Check if first arg is a string (Express API)
	if label, ok := args[0].(string); ok {
		cfg := DefaultConfirmConfig()
		cfg.Label = label
		return ConfirmWithConfig(cfg)
	}

	// Otherwise use Overload for config struct
	cfg := share.Overload(args, DefaultConfirmConfig())
	return ConfirmWithConfig(cfg)
}

// NewConfirm creates a new ConfirmBuilder for DSL chaining.
func NewConfirm() *ConfirmBuilder {
	return &ConfirmBuilder{config: DefaultConfirmConfig()}
}

// ConfirmWithConfig prompts the user for a yes/no confirmation with an explicit config.
func ConfirmWithConfig(cfg ConfirmConfig) (bool, error) {
	// Check for non-interactive environment
	if !terminal.IsTerminal(os.Stdin) && os.Getenv("FORM_NONINTERACTIVE") == "1" {
		return cfg.Default, nil
	}

	// Use interactive mode if enabled and available
	if cfg.Interactive {
		return confirmInteractive(cfg)
	}

	// Fall back to simple text mode
	return confirmSimple(cfg)
}

// confirmSimple provides a simple text-based confirmation prompt.
func confirmSimple(cfg ConfirmConfig) (bool, error) {
	// Determine prompt suffix based on default value
	suffix := " (y/N) "
	if cfg.Default {
		suffix = " (Y/n) "
	}

	prompt := cfg.Label + suffix

	// Loop until valid input
	for {
		fmt.Fprint(cfg.Writer, prompt)
		input, err := cfg.Reader.ReadLine(context.Background())
		if err != nil {
			return cfg.Default, err
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			return true, nil
		}
		if input == "n" || input == "no" {
			return false, nil
		}
		if input == "" {
			return cfg.Default, nil
		}
		fmt.Fprintln(cfg.Writer, "Please enter 'y' or 'n'.")
	}
}

// confirmInteractive provides an interactive confirmation prompt with visual feedback.
func confirmInteractive(cfg ConfirmConfig) (bool, error) {
	// Create an interactive confirm using RunFX
	confirmer := &ConfirmVisual{
		config:    cfg,
		selection: cfg.Default, // Start with default selection
		done:      make(chan bool, 1),
		canceled:  make(chan bool, 1),
	}

	// Create interactive loop and mount the visual component
	loop := runfx.StartInteractive()
	unmount, err := loop.MountInteractive(confirmer)
	if err != nil {
		// Fall back to simple mode if RunFX fails
		return confirmSimple(cfg)
	}
	defer unmount()

	// Start the main loop in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := loop.Run(ctx); err != nil {
			// Log error or handle it appropriately
		}
	}()

	// Wait for user confirmation or cancellation
	select {
	case result := <-confirmer.done:
		loop.Stop()
		return result, nil
	case <-confirmer.canceled:
		loop.Stop()
		return cfg.Default, ErrCanceled
	case <-ctx.Done():
		loop.Stop()
		return cfg.Default, ctx.Err()
	}
}

// ConfirmVisual implements runfx.Interactive for interactive confirmation with visual feedback.
type ConfirmVisual struct {
	config    ConfirmConfig
	selection bool // Current selection (true for Yes, false for No)
	done      chan bool
	canceled  chan bool
}

// Render implements runfx.Visual - displays the confirmation prompt.
func (cv *ConfirmVisual) Render(w share.Writer) {
	var output strings.Builder

	output.WriteString(cv.config.Label)
	output.WriteString("\n\n")

	// Show Yes/No options with visual indication
	yesPrefix := "  "
	noPrefix := "  "

	if cv.selection {
		yesPrefix = "▶ " // Arrow indicator for Yes
	} else {
		noPrefix = "▶ " // Arrow indicator for No
	}

	output.WriteString(fmt.Sprintf("%sYes\n", yesPrefix))
	output.WriteString(fmt.Sprintf("%sNo\n", noPrefix))

	// Show navigation hints
	output.WriteString("\n")
	output.WriteString("Use ↑↓ arrows, WASD, or Y/N to choose, Enter to confirm, Esc to cancel")

	// Write as a share.Entry
	entry := &share.Entry{
		Message: output.String(),
	}
	w.Write(entry)
}

// OnKey implements runfx.Interactive - handles keyboard input.
func (cv *ConfirmVisual) OnKey(key runfx.Key) bool {
	switch key {
	case runfx.KeyArrowUp, runfx.KeyW:
		cv.selection = true // Move to Yes
		return true
	case runfx.KeyArrowDown, runfx.KeyS:
		cv.selection = false // Move to No
		return true
	case runfx.KeyEnter:
		cv.done <- cv.selection
		return true
	case runfx.KeyEscape, runfx.KeyQ:
		cv.canceled <- true
		return true
	// Allow direct Y/N input
	case runfx.KeyY:
		cv.selection = true
		cv.done <- true
		return true
	case runfx.KeyN:
		cv.selection = false
		cv.done <- false
		return true
	}
	return false // Key not handled
}

// OnResize implements runfx.Visual - handles terminal resize.
func (cv *ConfirmVisual) OnResize(cols, rows int) {
	// ConfirmVisual doesn't need special resize handling
}

// Tick implements runfx.Visual - called on each render cycle.
func (cv *ConfirmVisual) Tick(now time.Time) {
	// ConfirmVisual doesn't need tick-based updates
}

// --- DSL BUILDER ---

// ConfirmBuilder provides a fluent API for building confirmation prompts.
type ConfirmBuilder struct {
	config ConfirmConfig
}

// Label sets the prompt label.
func (cb *ConfirmBuilder) Label(label string) *ConfirmBuilder {
	cb.config.Label = label
	return cb
}

// Default sets the default value.
func (cb *ConfirmBuilder) Default(value bool) *ConfirmBuilder {
	cb.config.Default = value
	return cb
}

// Writer sets the output writer.
func (cb *ConfirmBuilder) Writer(writer io.Writer) *ConfirmBuilder {
	cb.config.Writer = writer
	return cb
}

// Reader sets the input reader.
func (cb *ConfirmBuilder) Reader(reader Reader) *ConfirmBuilder {
	cb.config.Reader = reader
	return cb
}

// Interactive enables or disables interactive mode.
func (cb *ConfirmBuilder) Interactive(enabled bool) *ConfirmBuilder {
	cb.config.Interactive = enabled
	return cb
}

// Build creates a function that shows the confirmation prompt.
func (cb *ConfirmBuilder) Build() func() (bool, error) {
	config := cb.config
	return func() (bool, error) {
		return ConfirmWithConfig(config)
	}
}

// Show displays the confirmation prompt and returns the result.
func (cb *ConfirmBuilder) Show() (bool, error) {
	return ConfirmWithConfig(cb.config)
}
