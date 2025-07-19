package core

import (
	"context"
	"time"
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
	FormatBadge Format = iota
	FormatJSON
	FormatText
	FormatCustom
)

// BadgeStyle represents different badge appearances
type BadgeStyle int

const (
	BadgeStyleSquare BadgeStyle = iota
	BadgeStyleRound
	BadgeStyleArrow
	BadgeStyleDot
	BadgeStyleCustom
)

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

// Hook represents a function that can modify log entries
type Hook func(entry *Entry) *Entry
