package progrefx

import (
    "io"

    "github.com/garaekz/tfx/internal/share"
    "github.com/garaekz/tfx/runfx"
)

// ProgressConfig defines options for a progress bar.
// This mirrors the configuration from the legacy progress package but lives
// under progrefx to make it explicit that it is part of the new progress
// experience. Progress components created via this config are compatible
// with runfx for interactive rendering.
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

// DefaultProgressConfig returns sensible defaults for a progress bar.
func DefaultProgressConfig() ProgressConfig {
    return ProgressConfig{
        Total:     100,
        Label:     "Progress",
        Width:     40,
        DetectTTY: runfx.DetectTTY,
    }
}

// Start creates a Progress component using the provided options.  It accepts
// either a ProgressConfig or a sequence of functional options.  See
// internal/share.OverloadWithOptions for details.
func Start(opts ...any) *Progress {
    cfg := share.OverloadWithOptions[ProgressConfig](opts, DefaultProgressConfig())
    return newProgress(cfg)
}

// ProgressBuilder provides a fluent builder API for constructing a progress bar.
// It is useful when you need to set only a few fields on the ProgressConfig.
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

// Build constructs the Progress component with the configured options.
func (b *ProgressBuilder) Build() *Progress {
    return newProgress(b.config)
}