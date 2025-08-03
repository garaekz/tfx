package form

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
}

// DefaultConfirmConfig returns the default configuration for Confirm.
func DefaultConfirmConfig() ConfirmConfig {
	return ConfirmConfig{
		Label:   "Are you sure?",
		Default: false,
		Writer:  os.Stdout,
		Reader:  NewStdinReader(os.Stdin),
	}
}

// Confirm prompts the user for a yes/no confirmation.
// It supports the Express, Fluent, and Instantiated API patterns.
//
// Express API:
//
//	result, err := Confirm("Proceed?")
//
// Fluent API:
//
//	result, err := ConfirmWith(WithLabel("Proceed?"), WithDefault(true))
//
// Instantiated API:
//
//	cfg := DefaultConfirmConfig()
//	cfg.Label = "Proceed?"
//	result, err := ConfirmWithConfig(cfg)
func Confirm(label string) (bool, error) {
	return ConfirmWith(WithLabelConfirm(label))
}

// ConfirmWith prompts the user for a yes/no confirmation with functional options.
func ConfirmWith(opts ...Option[ConfirmConfig]) (bool, error) {
	cfg := DefaultConfirmConfig()
	ApplyOptions(&cfg, opts...)
	return ConfirmWithConfig(cfg)
}

// ConfirmWithConfig prompts the user for a yes/no confirmation with an explicit config.
func ConfirmWithConfig(cfg ConfirmConfig) (bool, error) {
	// Determine prompt suffix based on default value
	suffix := " (y/N) "
	if cfg.Default {
		suffix = " (Y/n) "
	}

	prompt := cfg.Label + suffix
	// Override non-interactive
	if !terminal.IsTerminal(os.Stdin) && os.Getenv("FORM_NONINTERACTIVE") == "1" {
		return cfg.Default, nil
	}
	// Loop hasta entrada válida
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

// WithLabelConfirm sets the label for the Confirm prompt.
func WithLabelConfirm(label string) Option[ConfirmConfig] {
	return func(cfg *ConfirmConfig) {
		cfg.Label = label
	}
}

// WithDefault sets the default value for the Confirm prompt.
func WithDefault(def bool) Option[ConfirmConfig] {
	return func(cfg *ConfirmConfig) {
		cfg.Default = def
	}
}

// WithConfirmWriter sets the output writer for the Confirm prompt.
func WithConfirmWriter(w io.Writer) Option[ConfirmConfig] {
	return func(cfg *ConfirmConfig) {
		cfg.Writer = w
	}
}

// WithConfirmReader sets the input reader for the Confirm prompt.
func WithConfirmReader(r Reader) Option[ConfirmConfig] {
	return func(cfg *ConfirmConfig) {
		cfg.Reader = r
	}
}
