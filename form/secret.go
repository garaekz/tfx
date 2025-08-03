package form

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/garaekz/tfx/terminal"
)

// ReadPasswordFunc is a function that reads a password from the terminal.
type ReadPasswordFunc func(fd int) ([]byte, error)

// MakeRawFunc is a function that puts the terminal into raw mode.
type MakeRawFunc func(fd uintptr) (*term.State, error)

// RestoreTerminalFunc is a function that restores the terminal to its original mode.
type RestoreTerminalFunc func(fd uintptr, state *term.State) error

// SecretConfig provides configuration for the Secret prompt.
type SecretConfig struct {
	// Label is the prompt message displayed to the user.
	Label string
	// Writer is the output writer for the prompt.
	Writer io.Writer
	// Reader is the input reader for the prompt (for non-TTY).
	Reader Reader
	// Confirm, if true, requires the user to enter the secret twice.
	Confirm bool
	// ReadPassword is the function used to read the password (for mocking).
	ReadPassword ReadPasswordFunc
	// IsTerminal is the function used to check if stdin is a terminal (for mocking).
	IsTerminal IsTTYFunc
	// MakeRaw is the function used to put the terminal into raw mode (for mocking).
	MakeRaw MakeRawFunc
	// RestoreTerminal is the function used to restore the terminal (for mocking).
	RestoreTerminal RestoreTerminalFunc
}

// DefaultSecretConfig returns the default configuration for Secret.
func DefaultSecretConfig() SecretConfig {
	return SecretConfig{
		Label:           "Enter secret:",
		Writer:          os.Stdout,
		Reader:          NewStdinReader(os.Stdin),
		Confirm:         false,
		ReadPassword:    term.ReadPassword,
		IsTerminal:      IsTTY,
		MakeRaw:         terminal.MakeRaw,
		RestoreTerminal: terminal.RestoreTerminal,
	}
}

// Secret prompts the user for a secret input (with echo off).
// It supports the Express, Fluent, and Instantiated API patterns.
//
// Express API:
//
//	result, err := Secret("Your password?")
//
// Fluent API:
//
//	result, err := SecretWith(WithLabelSecret("Your password?"), WithConfirmSecret(true))
//
// Instantiated API:
//
//	cfg := DefaultSecretConfig()
//	cfg.Label = "Your password?"
//	result, err := SecretWithConfig(cfg)
func Secret(label string) (string, error) {
	return SecretWith(WithLabelSecret(label))
}

// SecretWith prompts the user for a secret input with functional options.
func SecretWith(opts ...Option[SecretConfig]) (string, error) {
	cfg := DefaultSecretConfig()
	ApplyOptions(&cfg, opts...)
	return SecretWithConfig(cfg)
}

// SecretWithConfig prompts the user for a secret input with an explicit config.
func SecretWithConfig(cfg SecretConfig) (string, error) {
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	// Fallback non-TTY: simple prompt o no-interactive
	if !cfg.IsTerminal(os.Stdin) {
		if os.Getenv("FORM_NONINTERACTIVE") == "1" {
			return "", ErrCanceled
		}
		fmt.Fprint(cfg.Writer, cfg.Label+": ")
		input, err := cfg.Reader.ReadLine(context.Background())
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		return strings.TrimSpace(input), nil
	}

	// Save old state and set raw mode
	oldState, err := cfg.MakeRaw(uintptr(fd))
	if err != nil {
		return "", fmt.Errorf("failed to set raw terminal mode: %w", err)
	}
	defer cfg.RestoreTerminal(uintptr(fd), oldState)

	for {
		fmt.Fprint(cfg.Writer, cfg.Label+": ")
		byteSecret, err := cfg.ReadPassword(fd)
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		fmt.Fprintln(cfg.Writer) // Newline after input
		secret := string(byteSecret)

		if cfg.Confirm {
			fmt.Fprint(cfg.Writer, "Confirm secret: ")
			byteConfirm, err := cfg.ReadPassword(fd)
			if err != nil {
				if err == io.EOF {
					return "", io.EOF
				}
				return "", err
			}
			fmt.Fprintln(cfg.Writer) // Newline after input
			confirm := string(byteConfirm)

			if secret != confirm {
				fmt.Fprintln(cfg.Writer, "Secrets do not match. Please try again.")
				continue
			}
		}
		return secret, nil
	}
}

// WithLabelSecret sets the label for the Secret prompt.
func WithLabelSecret(label string) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.Label = label
	}
}

// WithConfirmSecret sets whether the Secret prompt requires confirmation.
func WithConfirmSecret(confirm bool) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.Confirm = confirm
	}
}

// WithSecretWriter sets the output writer for the Secret prompt.
func WithSecretWriter(w io.Writer) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.Writer = w
	}
}

// WithSecretReader sets the input reader for the Secret prompt.
func WithSecretReader(r Reader) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.Reader = r
	}
}

// WithReadPasswordFunc sets the function used to read the password.
func WithReadPasswordFunc(f ReadPasswordFunc) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.ReadPassword = f
	}
}

// WithIsTerminalSecret sets the function used to check if stdin is a terminal.
func WithIsTerminalSecret(f IsTTYFunc) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.IsTerminal = f
	}
}

// WithMakeRawFunc sets the function used to put the terminal into raw mode.
func WithMakeRawFunc(f MakeRawFunc) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.MakeRaw = f
	}
}

// WithRestoreTerminalFunc sets the function used to restore the terminal.
func WithRestoreTerminalFunc(f RestoreTerminalFunc) Option[SecretConfig] {
	return func(cfg *SecretConfig) {
		cfg.RestoreTerminal = f
	}
}
