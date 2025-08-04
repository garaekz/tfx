package runfx

import (
	"io"
	"os"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/terminal"
	"github.com/garaekz/tfx/writer"
)

// --- MULTIPATH API FUNCTIONS ---

// Start creates and starts a new Loop with multipath configuration support.
// opts Type: any = Option[Config] | Config
func Start(opts ...any) Loop {
	return newLoopWithConfig(share.OverloadWithOptions(opts, DefaultConfig()))
}

// newLoopWithConfig creates a new Loop with the given configuration
func newLoopWithConfig(cfg Config) Loop {
	ttyInfo := DetectTTYForOutput(cfg.Output)
	tw := writer.NewTerminalWriter(cfg.Output, writer.TerminalOptions{
		DoubleBuffer: true,
		DisableColor: ttyInfo.NoColor,
	})

	return &MainLoop{
		mux:     NewMultiplexer(),
		writer:  tw,
		reader:  NewKeyReader(os.Stdin),
		signals: terminal.NewSignalHandler(),
		events:  make(chan any, 64), // buffered channel for events
		ticker:  time.NewTicker(cfg.TickInterval),
	}
}

// --- DSL BUILDER API (Hardcore Path) ---

// LoopBuilder provides a fluent DSL interface for configuring RunFX loops
type LoopBuilder struct {
	config Config
}

// New creates a new LoopBuilder with default configuration.
// HARDCORE path - provides fluent DSL interface:
//
//	runfx.New().TickInterval(100*time.Millisecond).TestMode(true).Start()
func New() *LoopBuilder {
	return &LoopBuilder{config: DefaultConfig()}
}

// TickInterval sets the tick interval for the event loop
func (b *LoopBuilder) TickInterval(interval time.Duration) *LoopBuilder {
	b.config.TickInterval = interval
	return b
}

// Output sets the output writer
func (b *LoopBuilder) Output(output io.Writer) *LoopBuilder {
	b.config.Output = output
	return b
}

// SmoothAnimation sets tick interval to 30ms for very smooth animations
func (b *LoopBuilder) SmoothAnimation() *LoopBuilder {
	b.config.TickInterval = 30 * time.Millisecond
	return b
}

// FastAnimation sets tick interval to 100ms for faster/less resource intensive animations
func (b *LoopBuilder) FastAnimation() *LoopBuilder {
	b.config.TickInterval = 100 * time.Millisecond
	return b
}

// Start creates and returns the configured Loop instance
func (b *LoopBuilder) Start() Loop {
	return newLoopWithConfig(b.config)
}

func (b *LoopBuilder) AutoTick() *LoopBuilder {
	tty := DetectTTY()
	b.config.TickInterval = determineTickInterval(tty)
	return b
}

// --- Functional Options Path ---
// This approach provides a composable way to configure a loop.
// As noted in the TFX philosophy, it can serve as an alternative to the DSL path.

// StartWith creates and returns a new Loop configured with the provided functional options.
// It is the primary entry point for the Functional Options path.
func StartWith(cfg Config) Loop {
	return newLoopWithConfig(cfg)
}

// WithTickInterval returns an Option to set a custom tick interval.
func WithTickInterval(interval time.Duration) share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = interval
	}
}

// WithOutput returns an Option to set a custom output writer.
func WithOutput(output io.Writer) share.Option[Config] {
	return func(cfg *Config) {
		cfg.Output = output
	}
}

// WithSmoothAnimation returns an Option to set a 30ms tick interval for smooth animations.
func WithSmoothAnimation() share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = 30 * time.Millisecond
	}
}

// WithFastAnimation returns an Option to set a 100ms tick interval for efficient animations.
func WithFastAnimation() share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = 100 * time.Millisecond
	}
}

// WithAutoTick returns an Option to set the tick interval based on detected TTY capabilities.
func WithAutoTick() share.Option[Config] {
	return func(cfg *Config) {
		// Note: This requires access to the output writer from the config.
		// If the output is changed by another option, the order matters.
		tty := DetectTTYForOutput(cfg.Output)
		cfg.TickInterval = determineTickInterval(tty)
	}
}
