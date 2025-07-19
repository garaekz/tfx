package logx

import (
	"context"
	"io"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/shared"
)

// LogOptions represents logger configuration
type LogOptions struct {
	// Output settings
	Output     io.Writer
	Level      shared.Level
	Format     shared.Format
	Timestamp  bool
	TimeFormat string

	// Color settings
	ColorMode    color.Mode
	Theme        color.ColorTheme
	ForceColor   bool
	DisableColor bool

	// Badge settings
	BadgeWidth  int
	BadgeStyle  shared.BadgeStyle
	ShowCaller  bool
	CallerDepth int

	// File output
	LogFile     string
	FileLevel   shared.Level
	MaxFileSize int64
	MaxBackups  int
	MaxAge      int

	// Custom formatter
	CustomFormatter shared.Formatter
}

// Hook represents a function that can modify log entries
type Hook func(entry *shared.Entry) *shared.Entry

// DefaultOptions returns sensible defaults
func DefaultOptions() LogOptions {
	return LogOptions{
		Level:       shared.LevelInfo,
		Format:      shared.FormatBadge,
		Timestamp:   false,
		TimeFormat:  "2006-01-02 15:04:05",
		ColorMode:   color.ModeANSI,
		Theme:       color.DefaultTheme,
		BadgeWidth:  5,
		BadgeStyle:  shared.BadgeStyleSquare,
		ShowCaller:  false,
		CallerDepth: 3,
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		MaxBackups:  3,
		MaxAge:      30, // days
	}
}

// Context represents a logging context with fields
type Context struct {
	logger *Logger
	fields map[string]interface{}
	ctx    context.Context
}

// Logger represents the main logger instance
type Logger struct {
	options LogOptions
	writers []shared.Writer
	hooks   []Hook
	mu      sync.RWMutex
	ctx     context.Context
	wg      sync.WaitGroup
}
