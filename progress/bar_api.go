package progress

import (
	"io"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// ProgressConfig defines options for a progress bar.
type ProgressConfig struct {
	Total     int
	Label     string
	Width     int
	Theme     ProgressTheme
	Style     ProgressStyle
	Effect    ProgressEffect
	Writer    io.Writer // Used only for TTY detection, not direct writes.
	ShowETA   bool
	DetectTTY func() runfx.TTYInfo
}

// DefaultProgressConfig returns sensible defaults.
func DefaultProgressConfig() ProgressConfig {
	return ProgressConfig{
		Total:     100,
		Label:     "Progress",
		Width:     40,
		DetectTTY: runfx.DetectTTY,
	}
}

// Start creates a Progress component using the provided options.
func Start(opts ...any) *Progress {
	cfg := share.OverloadWithOptions[ProgressConfig](opts, DefaultProgressConfig())
	return newProgress(cfg)
}

// ProgressBuilder provides a fluent builder API.
type ProgressBuilder struct {
	config ProgressConfig
}

// NewProgressBuilder returns a builder with default configuration.
func NewProgressBuilder() *ProgressBuilder {
	return &ProgressBuilder{config: DefaultProgressConfig()}
}

// Total sets the total amount of work.
func (b *ProgressBuilder) Total(total int) *ProgressBuilder {
	b.config.Total = total
	return b
}

// Label sets the progress label.
func (b *ProgressBuilder) Label(label string) *ProgressBuilder {
	b.config.Label = label
	return b
}

// Width sets the bar width.
func (b *ProgressBuilder) Width(width int) *ProgressBuilder {
	b.config.Width = width
	return b
}

// DetectTTY allows providing a custom TTY detection function.
func (b *ProgressBuilder) DetectTTY(fn func() runfx.TTYInfo) *ProgressBuilder {
	b.config.DetectTTY = fn
	return b
}

// Build constructs the Progress component.
func (b *ProgressBuilder) Build() *Progress {
	return newProgress(b.config)
}
