package progress

import (
	"io"
	"os"
	"time"

	"github.com/garaekz/tfx/internal/share"
)

// --- SPINNER MULTIPATH API ---

// DefaultSpinnerConfig returns default configuration for Spinner.
func DefaultSpinnerConfig() SpinnerConfig {
	return SpinnerConfig{
		Message:  "",
		Frames:   []string{"|", "/", "-", "\\"},
		Interval: 100 * time.Millisecond,
		Theme:    MaterialTheme,
		Effect:   SpinnerEffectNone,
		Writer:   os.Stdout,
	}
}

// SpinnerBuilder provides a fluent interface for configuring spinners
type SpinnerBuilder struct {
	config SpinnerConfig
}

// --- MULTIPATH API FUNCTIONS ---

// StartSpinner creates and starts a spinner with multipath configuration support.
// Supports multiple usage patterns:
//   - StartSpinner()                   // Zero-config, uses defaults
//   - StartSpinner("message")          // Express mode with message
//   - StartSpinner(config)             // Config struct
func StartSpinner(args ...any) *Spinner {
	var cfg SpinnerConfig

	if len(args) == 0 {
		cfg = DefaultSpinnerConfig()
	} else if len(args) == 1 {
		switch v := args[0].(type) {
		case string:
			// Express mode: string message
			cfg = DefaultSpinnerConfig()
			cfg.Message = v
		case SpinnerConfig:
			cfg = v
		default:
			cfg = share.Overload(args, DefaultSpinnerConfig())
		}
	} else {
		cfg = share.Overload(args, DefaultSpinnerConfig())
	}

	s := NewSpinnerWithConfig(cfg)
	s.Start()
	return s
}

// NewSpinner creates a new SpinnerBuilder for DSL chaining (HARDCORE path)
func NewSpinner() *SpinnerBuilder {
	return &SpinnerBuilder{config: DefaultSpinnerConfig()}
}

// --- DSL BUILDER (HARDCORE PATH) ---

// Message sets the spinner message
func (b *SpinnerBuilder) Message(message string) *SpinnerBuilder {
	b.config.Message = message
	return b
}

// Frames sets custom frames for the spinner animation
func (b *SpinnerBuilder) Frames(frames []string) *SpinnerBuilder {
	b.config.Frames = frames
	return b
}

// Interval sets the frame interval duration
func (b *SpinnerBuilder) Interval(interval time.Duration) *SpinnerBuilder {
	b.config.Interval = interval
	return b
}

// Theme applies a ProgressTheme to colorize the spinner
func (b *SpinnerBuilder) Theme(theme ProgressTheme) *SpinnerBuilder {
	b.config.Theme = theme
	return b
}

// Effect applies visual effects
func (b *SpinnerBuilder) Effect(effect SpinnerEffect) *SpinnerBuilder {
	b.config.Effect = effect
	return b
}

// Writer sets a custom writer for spinner output
func (b *SpinnerBuilder) Writer(writer io.Writer) *SpinnerBuilder {
	b.config.Writer = writer
	return b
}

// Build creates a new Spinner instance without starting it
func (b *SpinnerBuilder) Build() *Spinner {
	return newSpinner(b.config)
}

// Start creates and starts the spinner
func (b *SpinnerBuilder) Start() *Spinner {
	s := newSpinner(b.config)
	s.Start()
	return s
}

// --- CONVENIENCE BUILDER METHODS ---

// MaterialTheme applies Material Design theme
func (b *SpinnerBuilder) MaterialTheme() *SpinnerBuilder {
	return b.Theme(MaterialTheme)
}

// DraculaTheme applies Dracula theme
func (b *SpinnerBuilder) DraculaTheme() *SpinnerBuilder {
	return b.Theme(DraculaTheme)
}

// NordTheme applies Nord theme
func (b *SpinnerBuilder) NordTheme() *SpinnerBuilder {
	return b.Theme(NordTheme)
}

// Rainbow enables rainbow effect
func (b *SpinnerBuilder) Rainbow() *SpinnerBuilder {
	return b.Effect(SpinnerEffectRainbow)
}

// --- PRESET FRAMES ---

// DotsFrames uses dots animation
func (b *SpinnerBuilder) DotsFrames() *SpinnerBuilder {
	return b.Frames([]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"})
}

// ArrowFrames uses arrow animation
func (b *SpinnerBuilder) ArrowFrames() *SpinnerBuilder {
	return b.Frames([]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"})
}

// BounceFrames uses bouncing animation
func (b *SpinnerBuilder) BounceFrames() *SpinnerBuilder {
	return b.Frames([]string{"⠁", "⠂", "⠄", "⠂"})
}

// --- INTERNAL IMPLEMENTATION ---

// newSpinner is the internal implementation
func newSpinner(cfg SpinnerConfig) *Spinner {
	return NewSpinnerWithConfig(cfg)
}
