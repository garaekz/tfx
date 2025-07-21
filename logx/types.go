package logx

import (
	"context"
	"io"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/share"
)

// LogOptions represents logger configuration
type LogOptions struct {
	// Output settings
	Output     io.Writer
	Level      share.Level
	Format     share.Format
	Timestamp  bool
	TimeFormat string

	// Color settings
	ColorMode    color.Mode
	Theme        color.ColorTheme
	ForceColor   bool
	DisableColor bool

	// Badge settings
	BadgeWidth  int
	BadgeStyle  share.BadgeStyle
	ShowCaller  bool
	CallerDepth int

	// File output
	LogFile     string
	FileLevel   share.Level
	MaxFileSize int64
	MaxBackups  int
	MaxAge      int

	// Custom formatter
	CustomFormatter share.Formatter
}

// Hook represents a function that can modify log entries
type Hook func(entry *share.Entry) *share.Entry

// DefaultOptions returns sensible defaults
func DefaultOptions() LogOptions {
	return LogOptions{
		Level:       share.LevelInfo,
		Format:      share.FormatBadge,
		Timestamp:   false,
		TimeFormat:  "2006-01-02 15:04:05",
		ColorMode:   color.ModeANSI,
		Theme:       color.DefaultTheme,
		BadgeWidth:  5,
		BadgeStyle:  share.BadgeStyleSquare,
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
	fields map[string]any
	ctx    context.Context
}

// Logger represents the main logger instance
type Logger struct {
	options LogOptions
	writers []share.Writer
	hooks   []Hook
	mu      sync.RWMutex
	ctx     context.Context
	wg      sync.WaitGroup
}
