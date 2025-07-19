package tfx

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/garaekz/tfx/color"
	writerpkg "github.com/garaekz/tfx/writer"
)

// Global logger instance
var globalLogger *Logger
var globalMu sync.RWMutex

func init() {
	globalLogger = New(DefaultOptions())
}

// New creates a new logger instance
func New(opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	logger := &Logger{
		options: opts,
		writers: []Writer{},
		hooks:   []Hook{},
		ctx:     context.Background(),
	}

	// Add default console writer
	cwOpts := writerpkg.Options{
		Level:        opts.Level,
		Format:       opts.Format,
		Timestamp:    opts.Timestamp,
		TimeFormat:   opts.TimeFormat,
		Theme:        opts.Theme,
		BadgeWidth:   opts.BadgeWidth,
		BadgeStyle:   opts.BadgeStyle,
		ShowCaller:   opts.ShowCaller,
		ForceColor:   opts.ForceColor,
		DisableColor: opts.DisableColor,
	}
	consoleWriter := writerpkg.NewConsoleWriter(opts.Output, cwOpts)
	logger.AddWriter(consoleWriter)

	// Add file writer if specified
	if opts.LogFile != "" {
		fwOpts := writerpkg.DefaultFileOptions()
		fwOpts.Level = opts.FileLevel
		fwOpts.Format = opts.Format
		fwOpts.MaxSize = opts.MaxFileSize
		fwOpts.MaxBackups = opts.MaxBackups
		fwOpts.MaxAge = opts.MaxAge

		fileWriter, err := writerpkg.NewFileWriter(opts.LogFile, fwOpts)
		if err == nil {
			logger.AddWriter(fileWriter)
		}
	}

	return logger
}

// Configure updates the global logger options
func Configure(opts Options) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalLogger = New(opts)
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalLogger
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Output = w
	// Update console writer
	for i, wr := range l.writers {
		if cw, ok := wr.(*writerpkg.ConsoleWriter); ok {
			cwOpts := writerpkg.Options{
				Level:        l.options.Level,
				Format:       l.options.Format,
				Timestamp:    l.options.Timestamp,
				TimeFormat:   l.options.TimeFormat,
				Theme:        l.options.Theme,
				BadgeWidth:   l.options.BadgeWidth,
				BadgeStyle:   l.options.BadgeStyle,
				ShowCaller:   l.options.ShowCaller,
				ForceColor:   l.options.ForceColor,
				DisableColor: l.options.DisableColor,
			}
			l.writers[i] = writerpkg.NewConsoleWriter(w, cwOpts)
			cw.Close()
			break
		}
	}
}

// SetFormat sets the output format
func (l *Logger) SetFormat(format Format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Format = format
}

// EnableTimestamp enables timestamp in logs
func (l *Logger) EnableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Timestamp = true
}

// DisableTimestamp disables timestamp in logs
func (l *Logger) DisableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Timestamp = false
}

// SetTheme sets the color theme
func (l *Logger) SetTheme(theme color.ColorTheme) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Theme = theme
}

// AddWriter adds a new writer
func (l *Logger) AddWriter(writer Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, writer)
}

// AddHook adds a new hook
func (l *Logger) AddHook(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, hook)
}

// WithFields creates a new context with fields
func (l *Logger) WithFields(fields Fields) *Context {
	return &Context{
		logger: l,
		fields: map[string]interface{}(fields),
		ctx:    l.ctx,
	}
}

// WithContext creates a new context with the given context
func (l *Logger) WithContext(ctx context.Context) *Context {
	return &Context{
		logger: l,
		fields: make(map[string]interface{}),
		ctx:    ctx,
	}
}

// shouldLog checks if the level should be logged
func (l *Logger) shouldLog(level Level) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return level >= l.options.Level
}

// createEntry creates a log entry
func (l *Logger) createEntry(level Level, msg string, fields Fields) *Entry {
	entry := &Entry{
		Level:     level,
		Message:   msg,
		Fields:    fields,
		Timestamp: time.Now(),
		Context:   l.ctx,
	}

	// Add caller info if enabled
	if l.options.ShowCaller {
		entry.Caller = l.getCaller()
	}

	// Apply hooks
	l.mu.RLock()
	hooks := l.hooks
	l.mu.RUnlock()

	for _, hook := range hooks {
		entry = hook(entry)
	}

	return entry
}

// getCaller gets caller information
func (l *Logger) getCaller() *CallerInfo {
	pc, file, line, ok := runtime.Caller(l.options.CallerDepth)
	if !ok {
		return nil
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return nil
	}

	return &CallerInfo{
		File:     file,
		Function: fn.Name(),
		Line:     line,
	}
}

// log writes a log entry
func (l *Logger) log(level Level, msg string, fields Fields) {
	if !l.shouldLog(level) {
		return
	}

	entry := l.createEntry(level, msg, fields)

	l.mu.RLock()
	writers := l.writers
	l.mu.RUnlock()

	for _, wr := range writers {
		go func(w Writer) {
			w.Write(entry)
		}(wr)
	}
}

// Logging methods
func (l *Logger) Trace(msg string, args ...interface{}) {
	l.log(LevelTrace, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(LevelFatal, fmt.Sprintf(msg, args...), nil)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.log(LevelPanic, msg, nil)
	panic(msg)
}

// Success is a convenience method for successful operations
func (l *Logger) Success(msg string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(msg, args...), Fields{"type": "success"})
}

// Badge creates a custom badge log
func (l *Logger) Badge(tag, msg string, color color.Color, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(msg, args...), Fields{
		"badge":       tag,
		"badge_color": color,
	})
}

// Close closes all writers
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, writer := range l.writers {
		if err := writer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Global functions that use the global logger
func SetLevel(level Level)                     { GetLogger().SetLevel(level) }
func SetOutput(w io.Writer)                    { GetLogger().SetOutput(w) }
func SetFormat(format Format)                  { GetLogger().SetFormat(format) }
func EnableTimestamp()                         { GetLogger().EnableTimestamp() }
func DisableTimestamp()                        { GetLogger().DisableTimestamp() }
func SetTheme(theme color.ColorTheme)          { GetLogger().SetTheme(theme) }
func AddWriter(writer Writer)                  { GetLogger().AddWriter(writer) }
func AddHook(hook Hook)                        { GetLogger().AddHook(hook) }
func WithFields(fields Fields) *Context        { return GetLogger().WithFields(fields) }
func WithContext(ctx context.Context) *Context { return GetLogger().WithContext(ctx) }

func Trace(msg string, args ...interface{})   { GetLogger().Trace(msg, args...) }
func Debug(msg string, args ...interface{})   { GetLogger().Debug(msg, args...) }
func Info(msg string, args ...interface{})    { GetLogger().Info(msg, args...) }
func Warn(msg string, args ...interface{})    { GetLogger().Warn(msg, args...) }
func Error(msg string, args ...interface{})   { GetLogger().Error(msg, args...) }
func Fatal(msg string, args ...interface{})   { GetLogger().Fatal(msg, args...) }
func Panic(msg string, args ...interface{})   { GetLogger().Panic(msg, args...) }
func Success(msg string, args ...interface{}) { GetLogger().Success(msg, args...) }

func Badge(tag, msg string, color color.Color, args ...interface{}) {
	GetLogger().Badge(tag, msg, color, args...)
}
