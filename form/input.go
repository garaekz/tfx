package form

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
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
}

// DefaultInputConfig returns the default configuration for Input.
func DefaultInputConfig() InputConfig {
	return InputConfig{
		Label:    "Enter value:",
		Default:  "",
		Writer:   os.Stdout,
		Reader:   NewStdinReader(os.Stdin),
		Validate: func(s string) error { return nil },
	}
}

// Input prompts the user for a string input.
// It supports the Express, Fluent, and Instantiated API patterns.
//
// Express API:
//
//	result, err := Input("Your name?")
//
// Fluent API:
//
//	result, err := InputWith(WithLabelInput("Your name?"), WithDefaultInput("John Doe"))
//
// Instantiated API:
//
//	cfg := DefaultInputConfig()
//	cfg.Label = "Your name?"
//	result, err := InputWithConfig(cfg)
func Input(label string) (string, error) {
	return InputWith(WithLabelInput(label))
}

// InputWith prompts the user for a string input with functional options.
func InputWith(opts ...Option[InputConfig]) (string, error) {
	cfg := DefaultInputConfig()
	ApplyOptions(&cfg, opts...)
	return InputWithConfig(cfg)
}

// InputWithConfig prompts the user for a string input with an explicit config.
func InputWithConfig(cfg InputConfig) (string, error) {
	// Bucle interactivo para todo tipo de Reader (maneja validación y EOF)
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
		if err := cfg.Validate(input); err != nil {
			fmt.Fprintf(cfg.Writer, "Validation error: %v\n", err)
			continue
		}
		return input, nil
	}
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
			if err == io.EOF {
				return "", ErrCanceled
			}
			return "", err
		}

		input = strings.TrimSpace(input)

		if input == "" {
			input = cfg.Default
		}

		if err := cfg.Validate(input); err != nil {
			fmt.Fprintf(cfg.Writer, "Validation error: %v\n", err)
			continue
		}

		return input, nil
	}
}

// WithLabelInput sets the label for the Input prompt.
func WithLabelInput(label string) Option[InputConfig] {
	return func(cfg *InputConfig) {
		cfg.Label = label
	}
}

// WithDefaultInput sets the default value for the Input prompt.
func WithDefaultInput(def string) Option[InputConfig] {
	return func(cfg *InputConfig) {
		cfg.Default = def
	}
}

// WithValidateInput sets the validation function for the Input prompt.
func WithValidateInput(validate func(string) error) Option[InputConfig] {
	return func(cfg *InputConfig) {
		cfg.Validate = validate
	}
}

// WithInputWriter sets the output writer for the Input prompt.
func WithInputWriter(w io.Writer) Option[InputConfig] {
	return func(cfg *InputConfig) {
		cfg.Writer = w
	}
}

// WithInputReader sets the input reader for the Input prompt.
func WithInputReader(r Reader) Option[InputConfig] {
	return func(cfg *InputConfig) {
		cfg.Reader = r
	}
}
