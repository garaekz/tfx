package formfx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// SelectConfig provides configuration for the Select prompt.
type SelectConfig struct {
	// Label is the prompt message displayed to the user.
	Label string
	// Options are the choices presented to the user.
	Options []string
	// Default is the index of the default selected option.
	Default int
	// Writer is the output writer for the prompt.
	Writer io.Writer
	// Reader is the input reader for the prompt.
	Reader Reader
	// PageSize is the number of options to display at once.
	PageSize int
	// Interactive enables RunFX-powered interactive mode with arrow key navigation.
	Interactive bool
}

// DefaultSelectConfig returns the default configuration for Select.
func DefaultSelectConfig() SelectConfig {
	// Detect if we're in an interactive environment
	ttyInfo := runfx.DetectTTY()

	return SelectConfig{
		Label:       "Choose an option:",
		Options:     []string{},
		Default:     0,
		Writer:      os.Stdout,
		Reader:      NewStdinReader(os.Stdin),
		PageSize:    10,
		Interactive: ttyInfo.IsTTY, // Enable interactive mode if TTY available
	}
}

// --- MULTIPATH API FUNCTIONS ---

// Select prompts the user to choose from a list of options with multipath configuration support.
// Supports multiple usage patterns:
//   - Select(label, options)                     // Express: simple label and options
//   - Select(config)                             // Instantiated: config struct
func Select(args ...any) (int, error) {
	// Handle different argument patterns
	if len(args) == 0 {
		// No args: use default config (will fail with no options)
		cfg := DefaultSelectConfig()
		return SelectWithConfig(cfg)
	}

	// Check if first two args are string and []string (Express API)
	if len(args) >= 2 {
		if label, ok := args[0].(string); ok {
			if options, ok := args[1].([]string); ok {
				cfg := DefaultSelectConfig()
				cfg.Label = label
				cfg.Options = options
				return SelectWithConfig(cfg)
			}
		}
	}

	// Otherwise use Overload for config struct
	cfg := share.Overload(args, DefaultSelectConfig())
	return SelectWithConfig(cfg)
}

// NewSelect creates a new SelectBuilder for DSL chaining.
func NewSelect() *SelectBuilder {
	return &SelectBuilder{config: DefaultSelectConfig()}
}

// SelectWithConfig prompts the user to choose from a list of options with an explicit config.
func SelectWithConfig(cfg SelectConfig) (int, error) {
	if len(cfg.Options) == 0 {
		return -1, errors.New("select: no options provided")
	}

	// Check for non-interactive environment
	ttyInfo := runfx.DetectTTY()
	if !ttyInfo.IsTTY && os.Getenv("FORM_NONINTERACTIVE") == "1" {
		return cfg.Default, nil
	}

	// Use interactive mode if enabled and available
	if cfg.Interactive && ttyInfo.IsTTY {
		return selectInteractive(cfg)
	}

	// Fall back to simple text mode
	return selectSimple(cfg)
}

// selectSimple provides a simple text-based selection prompt.
func selectSimple(cfg SelectConfig) (int, error) {
	for {
		fmt.Fprintln(cfg.Writer, cfg.Label)
		for i, opt := range cfg.Options {
			fmt.Fprintf(cfg.Writer, "%d) %s\n", i+1, opt)
		}
		fmt.Fprint(cfg.Writer, "Enter choice (number): ")
		input, err := cfg.Reader.ReadLine(context.Background())
		if err != nil {
			if err == io.EOF {
				return -1, io.EOF
			}
			return -1, err
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(cfg.Options) {
			fmt.Fprintf(cfg.Writer, "Invalid choice: %s\n", input)
			continue
		}
		return choice - 1, nil
	}
}

// selectInteractive provides an interactive selection prompt with arrow key navigation.
func selectInteractive(cfg SelectConfig) (int, error) {
	// Create an interactive select using RunFX
	selector := &SelectVisual{
		config:   cfg,
		selected: cfg.Default,
		done:     make(chan int, 1),
		canceled: make(chan bool, 1),
	}

	// Create interactive loop and mount the visual component
	loop := runfx.StartInteractive()
	unmount, err := loop.MountInteractive(selector)
	if err != nil {
		// Fall back to simple mode if RunFX fails
		return selectSimple(cfg)
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

	// Wait for user selection or cancellation
	select {
	case selectedIndex := <-selector.done:
		loop.Stop()
		return selectedIndex, nil
	case <-selector.canceled:
		loop.Stop()
		return -1, ErrCanceled
	case <-ctx.Done():
		loop.Stop()
		return -1, ctx.Err()
	}
}

// SelectVisual implements runfx.Visual for interactive selection with arrow keys.
type SelectVisual struct {
	config    SelectConfig
	selected  int
	scrollTop int
	done      chan int
	canceled  chan bool
}

// Render implements runfx.Visual - displays the selection menu.
func (sv *SelectVisual) Render(w share.Writer) {
	var output strings.Builder

	output.WriteString(sv.config.Label)
	output.WriteString("\n\n")

	// Calculate visible range based on page size
	start := sv.scrollTop
	end := start + sv.config.PageSize
	if end > len(sv.config.Options) {
		end = len(sv.config.Options)
	}

	// Render visible options
	for i := start; i < end; i++ {
		prefix := "  "
		if i == sv.selected {
			prefix = "▶ " // Arrow indicator for selected item
		}
		output.WriteString(fmt.Sprintf("%s%s\n", prefix, sv.config.Options[i]))
	}

	// Show navigation hints
	output.WriteString("\n")
	output.WriteString("Use ↑↓ arrows or WASD to navigate, Enter to select, Esc/q to cancel")

	// Write as a share.Entry
	entry := &share.Entry{
		Message: output.String(),
	}
	w.Write(entry)
}

// OnKey implements runfx.Interactive - handles keyboard input.
func (sv *SelectVisual) OnKey(key runfx.Key) bool {
	switch key {
	case runfx.KeyArrowUp, runfx.KeyW:
		sv.moveUp()
		return true
	case runfx.KeyArrowDown, runfx.KeyS:
		sv.moveDown()
		return true
	case runfx.KeyEnter:
		sv.done <- sv.selected
		return true
	case runfx.KeyEscape, runfx.KeyQ:
		sv.canceled <- true
		return true
	}
	return false // Key not handled
}

// OnResize implements runfx.Visual - handles terminal resize.
func (sv *SelectVisual) OnResize(cols, rows int) {
	// Adjust page size based on available terminal height
	if rows > 10 {
		sv.config.PageSize = rows - 6 // Leave space for prompt and instructions
	}
}

// Tick implements runfx.Visual - called on each render cycle.
func (sv *SelectVisual) Tick(now time.Time) {
	// SelectVisual doesn't need tick-based updates
}

// moveUp moves the selection cursor up.
func (sv *SelectVisual) moveUp() {
	if sv.selected > 0 {
		sv.selected--

		// Adjust scroll if needed
		if sv.selected < sv.scrollTop {
			sv.scrollTop = sv.selected
		}
	}
}

// moveDown moves the selection cursor down.
func (sv *SelectVisual) moveDown() {
	if sv.selected < len(sv.config.Options)-1 {
		sv.selected++

		// Adjust scroll if needed
		if sv.selected >= sv.scrollTop+sv.config.PageSize {
			sv.scrollTop = sv.selected - sv.config.PageSize + 1
		}
	}
}

// --- DSL BUILDER ---

// SelectBuilder provides a fluent API for building selection prompts.
type SelectBuilder struct {
	config SelectConfig
}

// Label sets the prompt label.
func (sb *SelectBuilder) Label(label string) *SelectBuilder {
	sb.config.Label = label
	return sb
}

// Options sets the available options.
func (sb *SelectBuilder) Options(options []string) *SelectBuilder {
	sb.config.Options = options
	return sb
}

// Default sets the default selected index.
func (sb *SelectBuilder) Default(index int) *SelectBuilder {
	sb.config.Default = index
	return sb
}

// Writer sets the output writer.
func (sb *SelectBuilder) Writer(writer io.Writer) *SelectBuilder {
	sb.config.Writer = writer
	return sb
}

// Reader sets the input reader.
func (sb *SelectBuilder) Reader(reader Reader) *SelectBuilder {
	sb.config.Reader = reader
	return sb
}

// PageSize sets the number of options to display at once.
func (sb *SelectBuilder) PageSize(size int) *SelectBuilder {
	sb.config.PageSize = size
	return sb
}

// Interactive enables or disables interactive mode.
func (sb *SelectBuilder) Interactive(enabled bool) *SelectBuilder {
	sb.config.Interactive = enabled
	return sb
}

// Build creates a function that shows the selection prompt.
func (sb *SelectBuilder) Build() func() (int, error) {
	config := sb.config
	return func() (int, error) {
		return SelectWithConfig(config)
	}
}

// Show displays the selection prompt and returns the result.
func (sb *SelectBuilder) Show() (int, error) {
	return SelectWithConfig(sb.config)
}
