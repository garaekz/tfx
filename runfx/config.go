package runfx

import (
	"io"
	"os"
	"time"

	"github.com/garaekz/tfx/internal/share"
)

// Config provides structured configuration for RunFX Loop
type Config struct {
	TickInterval time.Duration
	Output       io.Writer
	TestMode     bool
}

// DefaultConfig returns default configuration for RunFX
func DefaultConfig() Config {
	return Config{
		TickInterval: 50 * time.Millisecond, // Default 50ms for smooth animation
		Output:       os.Stdout,
		TestMode:     false,
	}
}

// --- FUNCTIONAL OPTIONS ---

// WithTickInterval sets the tick interval for the event loop
func WithTickInterval(interval time.Duration) share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = interval
	}
}

// WithOutput sets the output writer
func WithOutput(output io.Writer) share.Option[Config] {
	return func(cfg *Config) {
		cfg.Output = output
	}
}

// WithTestMode enables test mode
func WithTestMode() share.Option[Config] {
	return func(cfg *Config) {
		cfg.TestMode = true
	}
}

// WithSmoothAnimation sets tick interval to 30ms for very smooth animations
func WithSmoothAnimation() share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = 30 * time.Millisecond
	}
}

// WithFastAnimation sets tick interval to 100ms for faster/less resource intensive animations
func WithFastAnimation() share.Option[Config] {
	return func(cfg *Config) {
		cfg.TickInterval = 100 * time.Millisecond
	}
}
