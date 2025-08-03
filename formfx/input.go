package formfx

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// InputConfig provides configuration for the Input prompt.
type InputConfig struct {
	// Label is the prompt message displayed to the user.
	Label string
	// Default is the default value if the user just presses Enter.
	Default string
	// Writer is the output writer for the prompt.
	Writer io.Writer
	// Reader is the input reader for the prompt.
	Reader Reader
	// Validate is an optional validation function.
	Validate func(string) error
	// Interactive enables RunFX-powered interactive mode with enhanced editing.
	Interactive bool
}

// DefaultInputConfig returns the default configuration for Input.
func DefaultInputConfig() InputConfig {
	// Detect if we're in an interactive environment
	ttyInfo := runfx.DetectTTY()

	return InputConfig{
		Label:       "Enter value:",
		Default:     "",
		Writer:      os.Stdout,
		Reader:      NewStdinReader(os.Stdin),
		Validate:    func(s string) error { return nil },
		Interactive: ttyInfo.IsTTY, // Enable interactive mode if TTY available
	}
}

// --- MULTIPATH API FUNCTIONS ---

// Input prompts the user for a string input with multipath configuration support.
// Supports multiple usage patterns:
//   - Input(label)                               // Express: simple label
//   - Input(config)                              // Instantiated: config struct
func Input(args ...any) (string, error) {
	// Handle different argument patterns
	if len(args) == 0 {
		// No args: use default config
		cfg := DefaultInputConfig()
		return InputWithConfig(cfg)
	}

	// Check if first arg is a string (Express API)
	if label, ok := args[0].(string); ok {
		cfg := DefaultInputConfig()
		cfg.Label = label
		return InputWithConfig(cfg)
	}

	// Otherwise use Overload for config struct
	cfg := share.Overload(args, DefaultInputConfig())
	return InputWithConfig(cfg)
}

// NewInput creates a new InputBuilder for DSL chaining.
func NewInput() *InputBuilder {
	return &InputBuilder{config: DefaultInputConfig()}
}

// InputWithConfig prompts the user for a string input with an explicit config.
func InputWithConfig(cfg InputConfig) (string, error) {
	// Check for non-interactive environment
	ttyInfo := runfx.DetectTTY()
	if !ttyInfo.IsTTY && os.Getenv("FORM_NONINTERACTIVE") == "1" {
		return cfg.Default, nil
	}

	// Use interactive mode if enabled and available
	if cfg.Interactive && ttyInfo.IsTTY {
		return inputInteractive(cfg)
	}

	// Fall back to simple text mode
	return inputSimple(cfg)
}

// inputSimple provides a simple text-based input prompt.
func inputSimple(cfg InputConfig) (string, error) {
	for {
		var prompt string
		if cfg.Default != "" {
			prompt = fmt.Sprintf("%s (default: %s): ", cfg.Label, cfg.Default)
		} else {
			prompt = fmt.Sprintf("%s: ", cfg.Label)
		}
		fmt.Fprint(cfg.Writer, prompt)
		input, err := cfg.Reader.ReadLine(context.Background())
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			input = cfg.Default
		}

		// Validate input
		if err := cfg.Validate(input); err != nil {
			fmt.Fprintf(cfg.Writer, "Invalid input: %v\n", err)
			continue
		}

		return input, nil
	}
}

// inputInteractive provides an interactive input prompt with enhanced editing.
func inputInteractive(cfg InputConfig) (string, error) {
	// Create an interactive input using RunFX
	// This could support advanced editing features like history, autocomplete, etc.

	// For now, fall back to simple mode
	// TODO: Implement full RunFX integration with advanced text editing
	return inputSimple(cfg)
}

// --- DSL BUILDER ---

// InputBuilder provides a fluent API for building input prompts.
type InputBuilder struct {
	config InputConfig
}

// Label sets the prompt label.
func (ib *InputBuilder) Label(label string) *InputBuilder {
	ib.config.Label = label
	return ib
}

// Default sets the default value.
func (ib *InputBuilder) Default(value string) *InputBuilder {
	ib.config.Default = value
	return ib
}

// Writer sets the output writer.
func (ib *InputBuilder) Writer(writer io.Writer) *InputBuilder {
	ib.config.Writer = writer
	return ib
}

// Reader sets the input reader.
func (ib *InputBuilder) Reader(reader Reader) *InputBuilder {
	ib.config.Reader = reader
	return ib
}

// Validate sets the validation function.
func (ib *InputBuilder) Validate(validate func(string) error) *InputBuilder {
	ib.config.Validate = validate
	return ib
}

// Interactive enables or disables interactive mode.
func (ib *InputBuilder) Interactive(enabled bool) *InputBuilder {
	ib.config.Interactive = enabled
	return ib
}

// Build creates a function that shows the input prompt.
func (ib *InputBuilder) Build() func() (string, error) {
	config := ib.config
	return func() (string, error) {
		return InputWithConfig(config)
	}
}

// Show displays the input prompt and returns the result.
func (ib *InputBuilder) Show() (string, error) {
	return InputWithConfig(ib.config)
}
