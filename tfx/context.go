package tfx

import (
	"context"
	"fmt"
	"os"

	"github.com/garaekz/tfx/color"
)

// WithField adds a single field to the context
func (c *Context) WithField(key string, value interface{}) *Context {
	newFields := make(map[string]interface{})
	for k, v := range c.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &Context{
		logger: c.logger,
		fields: newFields,
		ctx:    c.ctx,
	}
}

// WithFields adds multiple fields to the context
func (c *Context) WithFields(fields Fields) *Context {
	newFields := make(map[string]interface{})
	for k, v := range c.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

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
func (c *Context) WithUser(userID interface{}) *Context {
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
func (c *Context) log(level Level, msg string) {
	if !c.logger.shouldLog(level) {
		return
	}

	// Merge context fields with any fields from context.Context
	allFields := make(Fields)

	// Add fields from the logging context
	for k, v := range c.fields {
		allFields[k] = v
	}

	// Extract fields from context.Context if any exist
	if c.ctx != nil {
		if ctxFields := extractContextFields(c.ctx); ctxFields != nil {
			for k, v := range ctxFields {
				allFields[k] = v
			}
		}
	}

	entry := c.logger.createEntry(level, msg, allFields)
	entry.Context = c.ctx

	c.logger.mu.RLock()
	writers := c.logger.writers
	c.logger.mu.RUnlock()

	for _, writer := range writers {
		go func(w Writer) {
			w.Write(entry)
		}(writer)
	}
}

// Logging methods for Context
func (c *Context) Trace(msg string, args ...interface{}) {
	c.log(LevelTrace, fmt.Sprintf(msg, args...))
}

func (c *Context) Debug(msg string, args ...interface{}) {
	c.log(LevelDebug, fmt.Sprintf(msg, args...))
}

func (c *Context) Info(msg string, args ...interface{}) {
	c.log(LevelInfo, fmt.Sprintf(msg, args...))
}

func (c *Context) Warn(msg string, args ...interface{}) {
	c.log(LevelWarn, fmt.Sprintf(msg, args...))
}

func (c *Context) Error(msg string, args ...interface{}) {
	c.log(LevelError, fmt.Sprintf(msg, args...))
}

func (c *Context) Fatal(msg string, args ...interface{}) {
	c.log(LevelFatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (c *Context) Panic(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	c.log(LevelPanic, msg)
	panic(msg)
}

func (c *Context) Success(msg string, args ...interface{}) {
	successFields := make(Fields)
	for k, v := range c.fields {
		successFields[k] = v
	}
	successFields["type"] = "success"

	entry := c.logger.createEntry(LevelInfo, fmt.Sprintf(msg, args...), successFields)
	entry.Context = c.ctx

	c.logger.mu.RLock()
	writers := c.logger.writers
	c.logger.mu.RUnlock()

	for _, writer := range writers {
		go func(w Writer) {
			w.Write(entry)
		}(writer)
	}
}

func (c *Context) Badge(tag, msg string, color color.Color, args ...interface{}) {
	badgeFields := make(Fields)
	for k, v := range c.fields {
		badgeFields[k] = v
	}
	badgeFields["badge"] = tag
	badgeFields["badge_color"] = color

	entry := c.logger.createEntry(LevelInfo, fmt.Sprintf(msg, args...), badgeFields)
	entry.Context = c.ctx

	c.logger.mu.RLock()
	writers := c.logger.writers
	c.logger.mu.RUnlock()

	for _, writer := range writers {
		go func(w Writer) {
			w.Write(entry)
		}(writer)
	}
}

// GetFields returns a copy of all fields in the context
func (c *Context) GetFields() Fields {
	fields := make(Fields)
	for k, v := range c.fields {
		fields[k] = v
	}
	return fields
}

// GetContext returns the context.Context
func (c *Context) GetContext() context.Context {
	return c.ctx
}

// Helper function to extract fields from context.Context
// This is a common pattern where you store logging fields in context
func extractContextFields(ctx context.Context) Fields {
	if ctx == nil {
		return nil
	}

	// Check for common context keys used for logging
	fields := make(Fields)

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

func FromContextWithFields(ctx context.Context, fields Fields) *Context {
	return GetLogger().WithContext(ctx).WithFields(fields)
}
