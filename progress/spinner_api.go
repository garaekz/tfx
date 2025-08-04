package progress

import (
	"io"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// SpinnerConfig holds configuration for Spinner.
type SpinnerConfig struct {
	Label     string
	Frames    []string
	Theme     ProgressTheme
	Writer    io.Writer // Used only for TTY detection.
	DetectTTY func() runfx.TTYInfo
}

// DefaultSpinnerConfig provides sensible defaults.
func DefaultSpinnerConfig() SpinnerConfig {
	return SpinnerConfig{
		Label:     "Loading",
		Frames:    []string{"|", "/", "-", "\\"},
		Theme:     MaterialTheme,
		DetectTTY: runfx.DetectTTY,
	}
}

// StartSpinner is a convenience function that creates a Spinner.
func StartSpinner(opts ...any) *Spinner {
	cfg := share.OverloadWithOptions[SpinnerConfig](opts, DefaultSpinnerConfig())
	return newSpinner(cfg)
}

// SpinnerBuilder offers a fluent DSL for building a Spinner.
type SpinnerBuilder struct {
	config SpinnerConfig
}

// NewSpinnerBuilder creates a new builder with default configuration.
func NewSpinnerBuilder() *SpinnerBuilder {
	return &SpinnerBuilder{config: DefaultSpinnerConfig()}
}

// Label sets the spinner label.
func (b *SpinnerBuilder) Label(label string) *SpinnerBuilder {
	b.config.Label = label
	return b
}

// Frames sets custom spinner frames.
func (b *SpinnerBuilder) Frames(frames []string) *SpinnerBuilder {
	b.config.Frames = frames
	return b
}

// Theme sets the spinner theme.
func (b *SpinnerBuilder) Theme(theme ProgressTheme) *SpinnerBuilder {
	b.config.Theme = theme
	return b
}

// DetectTTY allows providing a custom TTY detection function.
func (b *SpinnerBuilder) DetectTTY(fn func() runfx.TTYInfo) *SpinnerBuilder {
	b.config.DetectTTY = fn
	return b
}

// Build constructs the Spinner with the configured options.
func (b *SpinnerBuilder) Build() *Spinner {
	return newSpinner(b.config)
}
