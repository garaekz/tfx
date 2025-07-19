package logx

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/shared"
	writerpkg "github.com/garaekz/tfx/writer"
)

// Global logger instance
var globalLogger *Logger
var globalMu sync.RWMutex

func init() {
	globalLogger = New(DefaultOptions())
}

// New creates a new logger instance
func New(opts LogOptions) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	logger := &Logger{
		options: opts,
		writers: []shared.Writer{},
		hooks:   []Hook{},
		ctx:     context.Background(),
	}

	// Add default console writer
	cwOpts := writerpkg.ConsoleOptions{
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
func Configure(opts LogOptions) {
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
func (l *Logger) SetLevel(level shared.Level) {
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
			cwOpts := writerpkg.ConsoleOptions{
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
			cw.Close()
			l.writers[i] = writerpkg.NewConsoleWriter(w, cwOpts)
			break
		}
	}
}

func (l *Logger) Flush() {
	l.wg.Wait()
}

// SetFormat sets the output format
func (l *Logger) SetFormat(format shared.Format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Format = format
}

// EnableTimestamp enables timestamp in logs
func (l *Logger) EnableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Timestamp = true
	// Update console writers
	for i, wr := range l.writers {
		if cw, ok := wr.(*writerpkg.ConsoleWriter); ok {
			cwOpts := writerpkg.ConsoleOptions{
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
			cw.Close()
			l.writers[i] = writerpkg.NewConsoleWriter(l.options.Output, cwOpts)
		}
	}
}

// DisableTimestamp disables timestamp in logs
func (l *Logger) DisableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Timestamp = false
	// Update console writers
	for i, wr := range l.writers {
		if cw, ok := wr.(*writerpkg.ConsoleWriter); ok {
			cwOpts := writerpkg.ConsoleOptions{
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
			cw.Close()
			l.writers[i] = writerpkg.NewConsoleWriter(l.options.Output, cwOpts)
		}
	}
}

// SetTheme sets the color theme
func (l *Logger) SetTheme(theme color.ColorTheme) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Theme = theme
}

// AddWriter adds a new writer
func (l *Logger) AddWriter(writer shared.Writer) {
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
func (l *Logger) WithFields(fields shared.Fields) *Context {
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
func (l *Logger) shouldLog(level shared.Level) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return level >= l.options.Level
}

// createEntry creates a log entry
func (l *Logger) createEntry(level shared.Level, msg string, fields shared.Fields) *shared.Entry {
	entry := &shared.Entry{
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
func (l *Logger) getCaller() *shared.CallerInfo {
	pc, file, line, ok := runtime.Caller(l.options.CallerDepth)
	if !ok {
		return nil
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return nil
	}

	return &shared.CallerInfo{
		File:     file,
		Function: fn.Name(),
		Line:     line,
	}
}

// log writes a log entry
func (l *Logger) log(level shared.Level, msg string, fields shared.Fields) {
	if !l.shouldLog(level) {
		return
	}

	entry := l.createEntry(level, msg, fields)

	l.mu.RLock()
	writers := l.writers
	l.mu.RUnlock()

	for _, wr := range writers {
		l.wg.Add(1)
		go func(w shared.Writer) {
			defer l.wg.Done()
			w.Write(entry)
		}(wr)
	}
}

// Logging methods
func (l *Logger) Trace(msg string, args ...any) {
	l.log(shared.LevelTrace, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(shared.LevelDebug, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(shared.LevelInfo, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(shared.LevelWarn, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(shared.LevelError, fmt.Sprintf(msg, args...), nil)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.log(shared.LevelFatal, fmt.Sprintf(msg, args...), nil)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	l.log(shared.LevelPanic, msg, nil)
	panic(msg)
}

// Success is a convenience method for successful operations
func (l *Logger) Success(msg string, args ...any) {
	l.log(shared.LevelInfo, fmt.Sprintf(msg, args...), shared.Fields{"type": "success"})
}

// Badge creates a custom badge log
func (l *Logger) Badge(tag, msg string, color color.Color, args ...any) {
	l.log(shared.LevelInfo, fmt.Sprintf(msg, args...), shared.Fields{
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
func SetLevel(level shared.Level)              { GetLogger().SetLevel(level) }
func SetOutput(w io.Writer)                    { GetLogger().SetOutput(w) }
func SetFormat(format shared.Format)           { GetLogger().SetFormat(format) }
func EnableTimestamp()                         { GetLogger().EnableTimestamp() }
func DisableTimestamp()                        { GetLogger().DisableTimestamp() }
func SetTheme(theme color.ColorTheme)          { GetLogger().SetTheme(theme) }
func AddWriter(writer shared.Writer)           { GetLogger().AddWriter(writer) }
func AddHook(hook Hook)                        { GetLogger().AddHook(hook) }
func WithFields(fields shared.Fields) *Context { return GetLogger().WithFields(fields) }
func WithContext(ctx context.Context) *Context { return GetLogger().WithContext(ctx) }

func Trace(msg string, args ...any)   { GetLogger().Trace(msg, args...) }
func Debug(msg string, args ...any)   { GetLogger().Debug(msg, args...) }
func Info(msg string, args ...any)    { GetLogger().Info(msg, args...) }
func Warn(msg string, args ...any)    { GetLogger().Warn(msg, args...) }
func Error(msg string, args ...any)   { GetLogger().Error(msg, args...) }
func Fatal(msg string, args ...any)   { GetLogger().Fatal(msg, args...) }
func Panic(msg string, args ...any)   { GetLogger().Panic(msg, args...) }
func Success(msg string, args ...any) { GetLogger().Success(msg, args...) }

func Badge(tag, msg string, color color.Color, args ...any) {
	GetLogger().Badge(tag, msg, color, args...)
}

func Flush() { GetLogger().Flush() }
