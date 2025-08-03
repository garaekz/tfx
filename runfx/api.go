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
// BEGINNER path - supports two usage patterns:
//   - Start()         // Zero-config, uses defaults
//   - Start(config)   // Config struct
func Start(args ...any) Loop {
	cfg := share.Overload(args, DefaultConfig())
	return NewLoopWithConfig(cfg)
}

// StartInteractive creates and starts an interactive Loop with keyboard input.
// INTERACTIVE path - supports two usage patterns:
//   - StartInteractive()         // Interactive with defaults
//   - StartInteractive(config)   // Interactive with config
func StartInteractive(args ...any) InteractiveLoop {
	cfg := share.Overload(args, DefaultConfig())
	loop := NewLoopWithConfig(cfg)

	// Enable interactive mode
	if mainLoop, ok := loop.(*MainLoop); ok {
		mainLoop.EnableInteractive()
		return mainLoop
	}

	// Fallback - shouldn't happen with current implementation
	return loop.(*MainLoop)
}

// NewLoopWithConfig creates a new Loop with the given configuration
func NewLoopWithConfig(cfg Config) Loop {
	ttyInfo := DetectTTY()

	// Create terminal writer with configuration
	terminalOpts := writer.TerminalOptions{
		DoubleBuffer: true,
		ForceColor:   false,
		DisableColor: ttyInfo.NoColor,
	}
	termWriter := writer.NewTerminalWriter(cfg.Output, terminalOpts)

	return &MainLoop{
		mux:          NewMultiplexer(),
		terminal:     termWriter,
		cursor:       &CursorManager{},
		screen:       NewScreenManager(),
		signals:      *terminal.NewSignalHandler(),
		eventLoop:    NewEventLoop(cfg.TickInterval),
		ttyInfo:      &ttyInfo,
		keyReader:    NewKeyReader(os.Stdin),       // Initialize key reader with stdin
		interactives: make(map[Visual]Interactive), // Initialize interactives map
		inputEnabled: ttyInfo.IsTTY,                // Enable input in TTY environments
		stopCh:       make(chan struct{}),
		output:       cfg.Output,
		nextRegion:   0,
		testMode:     cfg.TestMode,
	}
}

// StartWith creates and starts a new Loop using functional options only
// EXPERIMENTAL path - not currently in active use, experimental feature
func StartWith(opts ...share.Option[Config]) Loop {
	cfg := DefaultConfig()
	share.ApplyOptions(&cfg, opts...)
	return NewLoopWithConfig(cfg)
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

// TestMode enables test mode
func (b *LoopBuilder) TestMode() *LoopBuilder {
	b.config.TestMode = true
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
	return NewLoopWithConfig(b.config)
}
