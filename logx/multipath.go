package logx

import (
	"io"
	"os"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/share"
)

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick defaults
// Log creates a logger with smart defaults based on arguments
func Log(args ...any) *Logger {
	return newLogger(args...)
}

// 2. INSTANTIATED: Config struct
// NewWithConfig creates a logger with explicit configuration
func NewWithConfig(cfg LogOptions) *Logger {
	return New(cfg)
}

// 3. FLUENT: Functional options + DSL chaining support
// NewWith creates a logger with functional options
func NewWith(opts ...LogOption) *Logger {
	cfg := DefaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}
	return New(cfg)
}

// newLogger is the internal implementation supporting multipath overload
func newLogger(args ...any) *Logger {
	// Separate functional options from other args
	var opts []share.Option[LogOptions]
	var cfgArgs []any

	for _, arg := range args {
		switch v := arg.(type) {
		case LogOption:
			// Convert LogOption to share.Option[LogOptions]
			opts = append(opts, share.Option[LogOptions](v))
		case share.Option[LogOptions]:
			opts = append(opts, v)
		default:
			cfgArgs = append(cfgArgs, v)
		}
	}

	// Use overload to handle different argument combinations
	cfg := share.OverloadWithOptions(cfgArgs, DefaultOptions(), opts...)
	return New(cfg)
}

// LogOption is a functional option for logger configuration
type LogOption func(*LogOptions)

// --- FUNCTIONAL OPTIONS ---

// WithOutput sets the output writer
func WithOutput(w io.Writer) LogOption {
	return func(cfg *LogOptions) {
		cfg.Output = w
	}
}

// WithLevel sets the minimum logging level
func WithLevel(level share.Level) LogOption {
	return func(cfg *LogOptions) {
		cfg.Level = level
	}
}

// WithFormat sets the output format
func WithFormat(format share.Format) LogOption {
	return func(cfg *LogOptions) {
		cfg.Format = format
	}
}

// WithTimestamp enables/disables timestamps
func WithTimestamp(enabled bool) LogOption {
	return func(cfg *LogOptions) {
		cfg.Timestamp = enabled
	}
}

// WithTimeFormat sets the timestamp format
func WithTimeFormat(format string) LogOption {
	return func(cfg *LogOptions) {
		cfg.TimeFormat = format
	}
}

// WithColorMode sets the color mode
func WithColorMode(mode color.Mode) LogOption {
	return func(cfg *LogOptions) {
		cfg.ColorMode = mode
	}
}

// WithTheme sets the color theme
func WithTheme(theme color.ColorTheme) LogOption {
	return func(cfg *LogOptions) {
		cfg.Theme = theme
	}
}

// WithForceColor forces color output
func WithForceColor(force bool) LogOption {
	return func(cfg *LogOptions) {
		cfg.ForceColor = force
	}
}

// WithDisableColor disables all color output
func WithDisableColor(disable bool) LogOption {
	return func(cfg *LogOptions) {
		cfg.DisableColor = disable
	}
}

// WithBadgeWidth sets the badge width
func WithBadgeWidth(width int) LogOption {
	return func(cfg *LogOptions) {
		cfg.BadgeWidth = width
	}
}

// WithLogBadgeStyle sets the badge style
func WithLogBadgeStyle(style share.BadgeStyle) LogOption {
	return func(cfg *LogOptions) {
		cfg.BadgeStyle = style
	}
}

// WithCaller enables/disables caller information
func WithCaller(enabled bool) LogOption {
	return func(cfg *LogOptions) {
		cfg.ShowCaller = enabled
	}
}

// WithCallerDepth sets the caller depth
func WithCallerDepth(depth int) LogOption {
	return func(cfg *LogOptions) {
		cfg.CallerDepth = depth
	}
}

// WithFileOutput enables file output
func WithFileOutput(filename string) LogOption {
	return func(cfg *LogOptions) {
		cfg.LogFile = filename
	}
}

// WithFileLevel sets the file logging level
func WithFileLevel(level share.Level) LogOption {
	return func(cfg *LogOptions) {
		cfg.FileLevel = level
	}
}

// WithFileRotation configures file rotation
func WithFileRotation(maxSize int64, maxBackups, maxAge int) LogOption {
	return func(cfg *LogOptions) {
		cfg.MaxFileSize = maxSize
		cfg.MaxBackups = maxBackups
		cfg.MaxAge = maxAge
	}
}

// WithCustomFormatter sets a custom formatter
func WithCustomFormatter(formatter share.Formatter) LogOption {
	return func(cfg *LogOptions) {
		cfg.CustomFormatter = formatter
	}
}

// --- CONVENIENCE OPTIONS ---

// WithJSON enables JSON format
func WithJSON() LogOption {
	return WithFormat(share.FormatJSON)
}

// WithBadges enables badge format (default)
func WithBadges() LogOption {
	return WithFormat(share.FormatBadge)
}

// WithText enables plain text format
func WithText() LogOption {
	return WithFormat(share.FormatText)
}

// WithDebugLevel sets level to Debug
func WithDebugLevel() LogOption {
	return WithLevel(share.LevelDebug)
}

// WithInfoLevel sets level to Info
func WithInfoLevel() LogOption {
	return WithLevel(share.LevelInfo)
}

// WithWarnLevel sets level to Warn
func WithWarnLevel() LogOption {
	return WithLevel(share.LevelWarn)
}

// WithErrorLevel sets level to Error
func WithErrorLevel() LogOption {
	return WithLevel(share.LevelError)
}

// WithDevelopment configures logger for development
func WithDevelopment() LogOption {
	return func(cfg *LogOptions) {
		cfg.Level = share.LevelDebug
		cfg.Format = share.FormatBadge
		cfg.Timestamp = true
		cfg.ShowCaller = true
		cfg.ForceColor = true
	}
}

// WithProduction configures logger for production
func WithProduction() LogOption {
	return func(cfg *LogOptions) {
		cfg.Level = share.LevelInfo
		cfg.Format = share.FormatJSON
		cfg.Timestamp = true
		cfg.ShowCaller = false
		cfg.DisableColor = true
	}
}

// WithConsoleOnly ensures only console output
func WithConsoleOnly() LogOption {
	return func(cfg *LogOptions) {
		cfg.LogFile = ""
	}
}

// WithStdout sets output to stdout
func WithStdout() LogOption {
	return WithOutput(os.Stdout)
}

// WithStderr sets output to stderr
func WithStderr() LogOption {
	return WithOutput(os.Stderr)
}

// --- THEME OPTIONS ---

// WithMaterialTheme uses Material Design colors
func WithMaterialTheme() LogOption {
	return WithTheme(color.MaterialTheme)
}

// WithDraculaTheme uses Dracula colors
func WithDraculaTheme() LogOption {
	return WithTheme(color.DraculaTheme)
}

// WithNordTheme uses Nord colors
func WithNordTheme() LogOption {
	return WithTheme(color.DefaultTheme) // TODO: Add NordTheme when available
}

// WithGitHubTheme uses GitHub colors
func WithGitHubTheme() LogOption {
	return WithTheme(color.DefaultTheme) // TODO: Add GitHubTheme when available
}

// --- PRESET CONSTRUCTORS ---

// DevLogger creates a development-focused logger
func DevLogger() *Logger {
	return NewWith(WithDevelopment())
}

// ProdLogger creates a production-focused logger
func ProdLogger() *Logger {
	return NewWith(WithProduction())
}

// TestLogger creates a logger suitable for testing
func TestLogger() *Logger {
	return NewWith(
		WithLevel(share.LevelDebug),
		WithFormat(share.FormatText),
		WithTimestamp(false),
		WithDisableColor(true),
	)
}

// ConsoleLogger creates a console-only logger with colors
func ConsoleLogger() *Logger {
	return NewWith(
		WithLevel(share.LevelInfo),
		WithFormat(share.FormatBadge),
		WithTimestamp(false),
		WithForceColor(true),
		WithConsoleOnly(),
	)
}

// FileLogger creates a file-only logger
func FileLogger(filename string) *Logger {
	return NewWith(
		WithLevel(share.LevelInfo),
		WithFormat(share.FormatJSON),
		WithTimestamp(true),
		WithFileOutput(filename),
		WithDisableColor(true),
	)
}

// StructuredLogger creates a JSON logger for structured logging
func StructuredLogger() *Logger {
	return NewWith(
		WithJSON(),
		WithTimestamp(true),
		WithLevel(share.LevelInfo),
		WithDisableColor(true),
	)
}
