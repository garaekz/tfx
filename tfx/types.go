package tfx

import (
	"context"
	"io"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/core"
)

// Re-export common types from core package
type Level = core.Level
type Format = core.Format
type BadgeStyle = core.BadgeStyle
type Fields = core.Fields
type Entry = core.Entry
type CallerInfo = core.CallerInfo
type Formatter = core.Formatter
type Writer = core.Writer
type Hook = core.Hook

const (
	LevelTrace = core.LevelTrace
	LevelDebug = core.LevelDebug
	LevelInfo  = core.LevelInfo
	LevelWarn  = core.LevelWarn
	LevelError = core.LevelError
	LevelFatal = core.LevelFatal
	LevelPanic = core.LevelPanic
)

const (
	FormatBadge  = core.FormatBadge
	FormatJSON   = core.FormatJSON
	FormatText   = core.FormatText
	FormatCustom = core.FormatCustom
)

const (
	BadgeStyleSquare = core.BadgeStyleSquare
	BadgeStyleRound  = core.BadgeStyleRound
	BadgeStyleArrow  = core.BadgeStyleArrow
	BadgeStyleDot    = core.BadgeStyleDot
	BadgeStyleCustom = core.BadgeStyleCustom
)

// Options represents logger configuration
type Options struct {
	// Output settings
	Output     io.Writer
	Level      Level
	Format     Format
	Timestamp  bool
	TimeFormat string

	// Color settings
	ColorMode    color.Mode
	Theme        color.ColorTheme
	ForceColor   bool
	DisableColor bool

	// Badge settings
	BadgeWidth  int
	BadgeStyle  BadgeStyle
	ShowCaller  bool
	CallerDepth int

	// File output
	LogFile     string
	FileLevel   Level
	MaxFileSize int64
	MaxBackups  int
	MaxAge      int

	// Custom formatter
	CustomFormatter Formatter
}

// DefaultOptions returns sensible defaults
func DefaultOptions() Options {
	return Options{
		Level:       LevelInfo,
		Format:      FormatBadge,
		Timestamp:   false,
		TimeFormat:  "15:04:05",
		ColorMode:   color.ModeANSI,
		Theme:       color.DefaultTheme,
		BadgeWidth:  5,
		BadgeStyle:  BadgeStyleSquare,
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
	options Options
	writers []Writer
	hooks   []Hook
	mu      sync.RWMutex
	ctx     context.Context
}

// Progress represents a progress bar state
type Progress struct {
	current int
	total   int
	message string
	width   int
	style   ProgressStyle
}

// ProgressStyle represents different progress bar styles
type ProgressStyle int

const (
	ProgressStyleBar    ProgressStyle = iota // [████████░░░░] 75%
	ProgressStyleDots                        // ••••••••···· 75%
	ProgressStyleArrows                      // >>>>>>>>>--- 75%
	ProgressStyleCustom                      // User-defined
)

// Spinner represents a loading spinner state
type Spinner struct {
	message string
	frames  []string
	index   int
	running bool
	done    chan bool
}
