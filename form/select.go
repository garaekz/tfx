package form

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	// IsTerminal is the function used to check if stdin is a terminal (for mocking).
	IsTerminal IsTTYFunc
}

// DefaultSelectConfig returns the default configuration for Select.
func DefaultSelectConfig() SelectConfig {
	return SelectConfig{
		Label:      "Choose an option:",
		Options:    []string{},
		Default:    0,
		Writer:     os.Stdout,
		Reader:     NewStdinReader(os.Stdin),
		PageSize:   10,
		IsTerminal: IsTTY,
	}
}

// Select prompts the user to choose from a list of options.
// It supports the Express, Fluent, and Instantiated API patterns.
//
// Express API:
//
//	index, err := Select("Choose a color:", []string{"Red", "Green", "Blue"})
//
// Fluent API:
//
//	index, err := SelectWith(WithLabelSelect("Choose a color:"), WithOptions([]string{"Red", "Green", "Blue"}))
//
// Instantiated API:
//
//	cfg := DefaultSelectConfig()
//	cfg.Label = "Choose a color:"
//	cfg.Options = []string{"Red", "Green", "Blue"}
//	index, err := SelectWithConfig(cfg)
func Select(label string, options []string) (int, error) {
	return SelectWith(WithLabelSelect(label), WithOptions(options))
}

// SelectWith prompts the user to choose from a list of options with functional options.
func SelectWith(opts ...Option[SelectConfig]) (int, error) {
	cfg := DefaultSelectConfig()
	ApplyOptions(&cfg, opts...)
	return SelectWithConfig(cfg)
}

// SelectWithConfig prompts the user to choose from a list of options with an explicit config.
func SelectWithConfig(cfg SelectConfig) (int, error) {
	if len(cfg.Options) == 0 {
		return -1, errors.New("select: no options provided")
	}

	// Non-TTY fallback: sin TTY, un solo prompt
	if !cfg.IsTerminal(os.Stdin) {
		if os.Getenv("FORM_NONINTERACTIVE") == "1" {
			return cfg.Default, nil
		}
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
				fmt.Fprintf(cfg.Writer, "invalid choice: %s\n", input)
				continue
			}
			return choice - 1, nil
		}
	}

	// TTY interactive mode (simplified)
	cursor := cfg.Default
	for {
		// Clear screen and redraw options
		fmt.Fprint(cfg.Writer, "\033[H\033[2J") // Clear screen
		fmt.Fprintln(cfg.Writer, cfg.Label)

		start := 0
		if cursor >= cfg.PageSize {
			start = cursor - cfg.PageSize + 1
		}
		end := min(start+cfg.PageSize, len(cfg.Options))

		for i := start; i < end; i++ {
			prefix := "  "
			if i == cursor {
				prefix = "> "
			}
			fmt.Fprintf(cfg.Writer, "%s%s\n", prefix, cfg.Options[i])
		}

		// Read input (simplified: just waits for Enter)
		input, err := cfg.Reader.ReadLine(context.Background())
		if err != nil {
			if err == io.EOF {
				return cursor, nil
			}
			return -1, ErrCanceled
		}
		trimmed := strings.ToLower(strings.TrimSpace(input))
		if trimmed == "up" {
			if cursor > 0 {
				cursor--
			}
			continue
		}
		if trimmed == "down" {
			if cursor < len(cfg.Options)-1 {
				cursor++
			}
			continue
		}
		if trimmed == "" {
			return cursor, nil
		}
		if trimmed == "q" || trimmed == "quit" || trimmed == "exit" {
			return -1, ErrCanceled
		}
		// cualquier otro input, continuar
	}
}

// WithLabelSelect sets the label for the Select prompt.
func WithLabelSelect(label string) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.Label = label
	}
}

// WithOptions sets the options for the Select prompt.
func WithOptions(options []string) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.Options = options
	}
}

// WithDefaultSelect sets the default selected option for the Select prompt.
func WithDefaultSelect(def int) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.Default = def
	}
}

// WithPageSize sets the number of options to display at once for the Select prompt.
func WithPageSize(size int) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.PageSize = size
	}
}

// WithSelectWriter sets the output writer for the Select prompt.
func WithSelectWriter(w io.Writer) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.Writer = w
	}
}

// WithSelectReader sets the input reader for the Select prompt.
func WithSelectReader(r Reader) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.Reader = r
	}
}

// WithIsTerminalSelect sets the function used to check if stdin is a terminal.
func WithIsTerminalSelect(f IsTTYFunc) Option[SelectConfig] {
	return func(cfg *SelectConfig) {
		cfg.IsTerminal = f
	}
}
