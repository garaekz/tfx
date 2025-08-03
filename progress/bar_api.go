package progress

import (
	"io"
	"os"

	"github.com/garaekz/tfx/internal/share"
)

// DefaultProgressConfig returns default configuration for Progress.
func DefaultProgressConfig() ProgressConfig {
	return ProgressConfig{
		Total:   100,
		Label:   "Progress",
		Width:   40,
		Theme:   MaterialTheme,
		Style:   ProgressStyleBar,
		Effect:  EffectNone,
		Writer:  os.Stdout,
		ShowETA: false,
	}
}

// ProgressBuilder is a DSL builder for creating and configuring Progress instances
type ProgressBuilder struct {
	config ProgressConfig
}

// --- MULTIPATH API FUNCTIONS ---

// Start creates and starts a progress bar with multipath configuration support.
// Supports two usage patterns:
//   - Start()                          // Zero-config, uses defaults
//   - Start(config)                    // Config struct
func Start(args ...any) *Progress {
	cfg := share.Overload(args, DefaultProgressConfig())
	p := newProgress(cfg)
	p.Start()
	return p
}

// New creates a new ProgressBuilder for DSL chaining (HARDCORE path)
func New() *ProgressBuilder {
	return &ProgressBuilder{config: DefaultProgressConfig()}
}

// --- DSL BUILDER ---

// Total sets the total value for progress tracking
func (b *ProgressBuilder) Total(total int) *ProgressBuilder {
	b.config.Total = total
	return b
}

// Label sets the progress label
func (b *ProgressBuilder) Label(label string) *ProgressBuilder {
	b.config.Label = label
	return b
}

// Width sets the width of the progress bar
func (b *ProgressBuilder) Width(width int) *ProgressBuilder {
	b.config.Width = width
	return b
}

// Theme applies a ProgressTheme
func (b *ProgressBuilder) Theme(theme ProgressTheme) *ProgressBuilder {
	b.config.Theme = theme
	return b
}

// Style applies a ProgressStyle
func (b *ProgressBuilder) Style(style ProgressStyle) *ProgressBuilder {
	b.config.Style = style
	return b
}

// Effect applies a ProgressEffect
func (b *ProgressBuilder) Effect(effect ProgressEffect) *ProgressBuilder {
	b.config.Effect = effect
	return b
}

// Writer sets a custom writer for progress output
func (b *ProgressBuilder) Writer(writer io.Writer) *ProgressBuilder {
	b.config.Writer = writer
	return b
}

// ShowETA enables/disables ETA display
func (b *ProgressBuilder) ShowETA() *ProgressBuilder {
	b.config.ShowETA = true
	return b
}

// Build creates a new Progress instance without starting it
func (b *ProgressBuilder) Build() *Progress {
	return newProgress(b.config)
}

// Start creates and starts the progress bar
func (b *ProgressBuilder) Start() *Progress {
	p := newProgress(b.config)
	p.Start()
	return p
}

// --- CONVENIENCE BUILDER METHODS ---

// MaterialTheme applies Material Design theme
func (b *ProgressBuilder) MaterialTheme() *ProgressBuilder {
	return b.Theme(MaterialTheme)
}

// DraculaTheme applies Dracula theme
func (b *ProgressBuilder) DraculaTheme() *ProgressBuilder {
	return b.Theme(DraculaTheme)
}

// NordTheme applies Nord theme
func (b *ProgressBuilder) NordTheme() *ProgressBuilder {
	return b.Theme(NordTheme)
}

// BarStyle sets bar style
func (b *ProgressBuilder) BarStyle() *ProgressBuilder {
	return b.Style(ProgressStyleBar)
}

// DotStyle sets dot style
func (b *ProgressBuilder) DotStyle() *ProgressBuilder {
	return b.Style(ProgressStyleDots)
}
