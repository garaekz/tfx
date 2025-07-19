package lfx

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/garaekz/lfx/internal/color"
)

// Level represents logging levels
type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// String returns string representation of level
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// Format represents output format
type Format int

const (
	FormatBadge  Format = iota // [TAG] message (default)
	FormatJSON                 // {"level":"info","msg":"..."}
	FormatText                 // 2024-01-01 15:04:05 INFO message
	FormatCustom               // User-defined format
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

// BadgeStyle represents different badge appearances
type BadgeStyle int

const (
	BadgeStyleSquare BadgeStyle = iota // [TAG]
	BadgeStyleRound                    // (TAG)
	BadgeStyleArrow                    // >TAG<
	BadgeStyleDot                      // •TAG•
	BadgeStyleCustom                   // User-defined
)

// Context represents a logging context with fields
type Context struct {
	logger *Logger
	fields map[string]interface{}
	ctx    context.Context
}

// Fields represents key-value pairs for structured logging
type Fields map[string]interface{}

// Entry represents a single log entry
type Entry struct {
	Level     Level
	Message   string
	Fields    Fields
	Timestamp time.Time
	Caller    *CallerInfo
	Context   context.Context
}

// CallerInfo represents information about the calling function
type CallerInfo struct {
	File     string
	Function string
	Line     int
}

// Formatter defines the interface for custom formatters
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// Writer defines the interface for log writers
type Writer interface {
	Write(entry *Entry) error
	Close() error
}

// MultiWriter allows writing to multiple destinations
type MultiWriter struct {
	writers []Writer
}

// Hook represents a function that can modify log entries
type Hook func(entry *Entry) *Entry

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
