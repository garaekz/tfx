package logx

import (
	"context"
	"fmt"
	"maps"
	"os"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/share"
)

// WithField adds a single field to the context
func (c *Context) WithField(key string, value any) *Context {
	newFields := make(map[string]any)
	maps.Copy(newFields, c.fields)
	newFields[key] = value

	return &Context{
		logger: c.logger,
		fields: newFields,
		ctx:    c.ctx,
	}
}

// WithFields adds multiple fields to the context
func (c *Context) WithFields(fields share.Fields) *Context {
	newFields := make(map[string]any)
	maps.Copy(newFields, c.fields)
	maps.Copy(newFields, fields)

	return &Context{
		logger: c.logger,
		fields: newFields,
		ctx:    c.ctx,
	}
}

// WithContext adds or updates the context.Context
func (c *Context) WithContext(ctx context.Context) *Context {
	return &Context{
		logger: c.logger,
		fields: c.fields,
		ctx:    ctx,
	}
}

// WithError adds an error field to the context
func (c *Context) WithError(err error) *Context {
	return c.WithField("error", err.Error())
}

// WithUser adds user-related fields (common pattern)
func (c *Context) WithUser(userID any) *Context {
	return c.WithField("user_id", userID)
}

// WithRequestID adds a request ID field (common in web apps)
func (c *Context) WithRequestID(requestID string) *Context {
	return c.WithField("request_id", requestID)
}

// WithSession adds a session ID field
func (c *Context) WithSession(sessionID string) *Context {
	return c.WithField("session_id", sessionID)
}

// WithTraceID adds a trace ID for distributed tracing
func (c *Context) WithTraceID(traceID string) *Context {
	return c.WithField("trace_id", traceID)
}

// log is the internal method that creates entries with fields
func (c *Context) log(level share.Level, msg string) {
	if !c.logger.shouldLog(level) {
		return
	}

	// Merge context fields with any fields from context.Context
	allFields := make(share.Fields)

	// Add fields from the logging context
	maps.Copy(allFields, c.fields)

	// Extract fields from context.Context if any exist
	if c.ctx != nil {
		if ctxFields := extractContextFields(c.ctx); ctxFields != nil {
			maps.Copy(allFields, ctxFields)
		}
	}

	entry := c.logger.createEntry(level, msg, allFields)
	entry.Context = c.ctx

	c.logger.mu.RLock()
	writers := c.logger.writers
	c.logger.mu.RUnlock()

	for _, writer := range writers {
		go func(w share.Writer) {
			w.Write(entry)
		}(writer)
	}
}

// Logging methods for Context
func (c *Context) Trace(msg string, args ...any) {
	c.log(share.LevelTrace, fmt.Sprintf(msg, args...))
}

func (c *Context) Debug(msg string, args ...any) {
	c.log(share.LevelDebug, fmt.Sprintf(msg, args...))
}

func (c *Context) Info(msg string, args ...any) {
	c.log(share.LevelInfo, fmt.Sprintf(msg, args...))
}

func (c *Context) Warn(msg string, args ...any) {
	c.log(share.LevelWarn, fmt.Sprintf(msg, args...))
}

func (c *Context) Error(msg string, args ...any) {
	c.log(share.LevelError, fmt.Sprintf(msg, args...))
}

func (c *Context) Fatal(msg string, args ...any) {
	c.log(share.LevelFatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (c *Context) Panic(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	c.log(share.LevelPanic, msg)
	panic(msg)
}

func (c *Context) Success(msg string, args ...any) {
	c.log(share.LevelSuccess, fmt.Sprintf(msg, args...))
}

// FatalIf logs a fatal message with context if err is not nil and exits
func (c *Context) FatalIf(err error, msg string, args ...any) {
	if err != nil {
		formattedMsg := fmt.Sprintf(msg, args...)
		errorFields := make(share.Fields)
		maps.Copy(errorFields, c.fields)
		errorFields["error"] = err.Error()

		c.logger.log(share.LevelFatal, fmt.Sprintf("%s: %v", formattedMsg, err), errorFields)
		os.Exit(1)
	}
}

// ErrorIf logs an error message with context if err is not nil and returns true if error occurred
func (c *Context) ErrorIf(err error, msg string, args ...any) bool {
	if err != nil {
		formattedMsg := fmt.Sprintf(msg, args...)
		errorFields := make(share.Fields)
		maps.Copy(errorFields, c.fields)
		errorFields["error"] = err.Error()

		c.logger.log(share.LevelError, fmt.Sprintf("%s: %v", formattedMsg, err), errorFields)
		return true
	}
	return false
}

// WarnIf logs a warning message with context if err is not nil and returns true if error occurred
func (c *Context) WarnIf(err error, msg string, args ...any) bool {
	if err != nil {
		formattedMsg := fmt.Sprintf(msg, args...)
		errorFields := make(share.Fields)
		maps.Copy(errorFields, c.fields)
		errorFields["error"] = err.Error()

		c.logger.log(share.LevelWarn, fmt.Sprintf("%s: %v", formattedMsg, err), errorFields)
		return true
	}
	return false
}

// InfoIf logs an info message with context if err is not nil and returns true if error occurred
func (c *Context) InfoIf(err error, msg string, args ...any) bool {
	if err != nil {
		formattedMsg := fmt.Sprintf(msg, args...)
		errorFields := make(share.Fields)
		maps.Copy(errorFields, c.fields)
		errorFields["error"] = err.Error()

		c.logger.log(share.LevelInfo, fmt.Sprintf("%s: %v", formattedMsg, err), errorFields)
		return true
	}
	return false
}

// DebugIf logs a debug message with context if err is not nil and returns true if error occurred
func (c *Context) DebugIf(err error, msg string, args ...any) bool {
	if err != nil {
		formattedMsg := fmt.Sprintf(msg, args...)
		errorFields := make(share.Fields)
		maps.Copy(errorFields, c.fields)
		errorFields["error"] = err.Error()

		c.logger.log(share.LevelDebug, fmt.Sprintf("%s: %v", formattedMsg, err), errorFields)
		return true
	}
	return false
}

func (c *Context) Badge(tag, msg string, color color.Color, args ...any) {
	badgeFields := make(share.Fields)
	maps.Copy(badgeFields, c.fields)
	badgeFields["badge"] = tag
	badgeFields["badge_color"] = color

	entry := c.logger.createEntry(share.LevelInfo, fmt.Sprintf(msg, args...), badgeFields)
	entry.Context = c.ctx

	c.logger.mu.RLock()
	writers := c.logger.writers
	c.logger.mu.RUnlock()

	for _, writer := range writers {
		go func(w share.Writer) {
			w.Write(entry)
		}(writer)
	}
}

// GetFields returns a copy of all fields in the context
func (c *Context) GetFields() share.Fields {
	fields := make(share.Fields)
	maps.Copy(fields, c.fields)
	return fields
}

// GetContext returns the context.Context
func (c *Context) GetContext() context.Context {
	return c.ctx
}

// Helper function to extract fields from context.Context
// This is a common pattern where you store logging fields in context
func extractContextFields(ctx context.Context) share.Fields {
	if ctx == nil {
		return nil
	}

	// Check for common context keys used for logging
	fields := make(share.Fields)

	// Request ID
	if reqID := ctx.Value("request_id"); reqID != nil {
		fields["request_id"] = reqID
	}

	// User ID
	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}

	// Session ID
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		fields["session_id"] = sessionID
	}

	// Trace ID
	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}

	// Correlation ID
	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		fields["correlation_id"] = correlationID
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

// Convenience functions for creating contexts from context.Context
func FromContext(ctx context.Context) *Context {
	return GetLogger().WithContext(ctx)
}

func FromContextWithFields(ctx context.Context, fields share.Fields) *Context {
	return GetLogger().WithContext(ctx).WithFields(fields)
}
