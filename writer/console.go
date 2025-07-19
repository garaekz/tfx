package writers

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
)

// ConsoleWriter writes formatted logs to console/terminal
type ConsoleWriter struct {
	output     io.Writer
	options    Options
	detector   *terminal.Detector
	badgeWidth int
	mu         sync.Mutex
}

// Options for console writer
type Options struct {
	Level        Level
	Format       Format
	Timestamp    bool
	TimeFormat   string
	Theme        color.ColorTheme
	BadgeWidth   int
	BadgeStyle   BadgeStyle
	ShowCaller   bool
	ForceColor   bool
	DisableColor bool
}

// NewConsoleWriter creates a new console writer
func NewConsoleWriter(output io.Writer, opts Options) *ConsoleWriter {
	if opts.TimeFormat == "" {
		opts.TimeFormat = "15:04:05"
	}
	if opts.BadgeWidth == 0 {
		opts.BadgeWidth = 5
	}

	return &ConsoleWriter{
		output:     output,
		options:    opts,
		detector:   terminal.NewDetector(output),
		badgeWidth: opts.BadgeWidth,
	}
}

// Write writes a log entry to console
func (w *ConsoleWriter) Write(entry *Entry) error {
	if entry.Level < w.options.Level {
		return nil
	}

	var output string
	switch w.options.Format {
	case FormatBadge:
		output = w.formatBadge(entry)
	case FormatJSON:
		output = w.formatJSON(entry)
	case FormatText:
		output = w.formatText(entry)
	default:
		output = w.formatBadge(entry)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := fmt.Fprintln(w.output, output)
	return err
}

// formatBadge formats entry as a badge log
func (w *ConsoleWriter) formatBadge(entry *Entry) string {
	var parts []string

	// Timestamp
	if w.options.Timestamp {
		timestamp := entry.Timestamp.Format(w.options.TimeFormat)
		if w.supportsColor() {
			timestamp = color.ApplyColor(timestamp, color.BrightBlack)
		}
		parts = append(parts, fmt.Sprintf("[%s]", timestamp))
	}

	// Badge/Level
	badge := w.formatBadgeTag(entry)
	parts = append(parts, badge)

	// Caller info
	if w.options.ShowCaller && entry.Caller != nil {
		caller := fmt.Sprintf("%s:%d", w.shortFilename(entry.Caller.File), entry.Caller.Line)
		if w.supportsColor() {
			caller = color.ApplyColor(caller, color.BrightBlack)
		}
		parts = append(parts, fmt.Sprintf("[%s]", caller))
	}

	// Message
	message := entry.Message
	if w.supportsColor() {
		message = w.colorizeMessage(entry, message)
	}
	parts = append(parts, message)

	// Fields
	if len(entry.Fields) > 0 {
		fieldsStr := w.formatFields(entry.Fields)
		if fieldsStr != "" {
			if w.supportsColor() {
				fieldsStr = color.ApplyColor(fieldsStr, color.BrightBlack)
			}
			parts = append(parts, fieldsStr)
		}
	}

	return strings.Join(parts, " ")
}

// formatBadgeTag formats the badge/level part
func (w *ConsoleWriter) formatBadgeTag(entry *Entry) string {
	var tag string
	var tagColor color.Color

	// Check if it's a custom badge
	if badgeTag, ok := entry.Fields["badge"].(string); ok {
		tag = badgeTag
		if badgeColor, ok := entry.Fields["badge_color"].(color.Color); ok {
			tagColor = badgeColor
		}
	} else {
		// Use level-based tag
		tag = w.getLevelTag(entry.Level)
		tagColor = w.getLevelColor(entry.Level)
	}

	// Pad tag to consistent width
	paddedTag := w.padTag(tag)

	// Apply badge style
	styledTag := w.applyBadgeStyle(paddedTag)

	// Apply color if supported
	if w.supportsColor() {
		colorCode := tagColor.Render(w.getColorMode())
		if colorCode != "" {
			return colorCode + styledTag + color.Reset
		}
	}

	return styledTag
}

// padTag pads the tag to maintain consistent width
func (w *ConsoleWriter) padTag(tag string) string {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(tag) > w.badgeWidth {
		w.badgeWidth = len(tag)
	}
	return tag + strings.Repeat(" ", w.badgeWidth-len(tag))
}

// applyBadgeStyle applies the badge style
func (w *ConsoleWriter) applyBadgeStyle(tag string) string {
	switch w.options.BadgeStyle {
	case BadgeStyleSquare:
		return fmt.Sprintf("[%s]", tag)
	case BadgeStyleRound:
		return fmt.Sprintf("(%s)", tag)
	case BadgeStyleArrow:
		return fmt.Sprintf(">%s<", tag)
	case BadgeStyleDot:
		return fmt.Sprintf("•%s•", tag)
	default:
		return fmt.Sprintf("[%s]", tag)
	}
}

// getLevelTag returns the tag for a level
func (w *ConsoleWriter) getLevelTag(level Level) string {
	switch level {
	case LevelTrace:
		return "TRC"
	case LevelDebug:
		return "DBG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERR"
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	default:
		return "LOG"
	}
}

// getLevelColor returns the color for a level
func (w *ConsoleWriter) getLevelColor(level Level) color.Color {
	switch level {
	case LevelTrace:
		return color.NewANSI(8) // Dark gray
	case LevelDebug:
		return w.options.Theme.Debug
	case LevelInfo:
		return w.options.Theme.Info
	case LevelWarn:
		return w.options.Theme.Warning
	case LevelError:
		return w.options.Theme.Error
	case LevelFatal:
		return color.NewRGB(255, 255, 255) // White on red bg
	case LevelPanic:
		return color.NewRGB(255, 255, 255) // White on red bg
	default:
		return w.options.Theme.Info
	}
}

// colorizeMessage applies color to the message based on level
func (w *ConsoleWriter) colorizeMessage(entry *Entry, message string) string {
	// Check for success type
	if msgType, ok := entry.Fields["type"].(string); ok && msgType == "success" {
		return color.ApplyColor(message, w.options.Theme.Success.Render(w.getColorMode()))
	}

	// Apply subtle coloring based on level
	switch entry.Level {
	case LevelError, LevelFatal, LevelPanic:
		return color.ApplyColor(message, color.BrightRed)
	case LevelWarn:
		return message // Keep message neutral for warnings
	default:
		return message // Keep message neutral for info/debug
	}
}

// formatFields formats the fields for display
func (w *ConsoleWriter) formatFields(fields Fields) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	for key, value := range fields {
		// Skip internal fields
		if key == "badge" || key == "badge_color" || key == "type" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	if len(parts) == 0 {
		return ""
	}

	return fmt.Sprintf("(%s)", strings.Join(parts, " "))
}

// formatJSON formats entry as JSON (simple implementation)
func (w *ConsoleWriter) formatJSON(entry *Entry) string {
	// This is a simplified JSON formatter
	// In a real implementation, you'd use json.Marshal
	parts := []string{
		fmt.Sprintf(`"level":"%s"`, entry.Level.String()),
		fmt.Sprintf(`"msg":"%s"`, entry.Message),
		fmt.Sprintf(`"time":"%s"`, entry.Timestamp.Format("2006-01-02T15:04:05.000Z")),
	}

	// Add fields
	for key, value := range entry.Fields {
		if key == "badge" || key == "badge_color" {
			continue
		}
		parts = append(parts, fmt.Sprintf(`"%s":"%v"`, key, value))
	}

	// Add caller if available
	if entry.Caller != nil {
		parts = append(parts, fmt.Sprintf(`"caller":"%s:%d"`, entry.Caller.File, entry.Caller.Line))
	}

	return fmt.Sprintf("{%s}", strings.Join(parts, ","))
}

// formatText formats entry as plain text
func (w *ConsoleWriter) formatText(entry *Entry) string {
	var parts []string

	// Timestamp
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
	parts = append(parts, timestamp)

	// Level
	parts = append(parts, entry.Level.String())

	// Caller
	if w.options.ShowCaller && entry.Caller != nil {
		caller := fmt.Sprintf("%s:%d", w.shortFilename(entry.Caller.File), entry.Caller.Line)
		parts = append(parts, caller)
	}

	// Message
	parts = append(parts, entry.Message)

	// Fields
	if len(entry.Fields) > 0 {
		fieldsStr := w.formatFields(entry.Fields)
		if fieldsStr != "" {
			parts = append(parts, fieldsStr)
		}
	}

	return strings.Join(parts, " ")
}

// Helper methods
func (w *ConsoleWriter) supportsColor() bool {
	if w.options.ForceColor {
		return true
	}
	if w.options.DisableColor {
		return false
	}
	return w.detector.SupportsANSI()
}

func (w *ConsoleWriter) getColorMode() color.Mode {
	if !w.supportsColor() {
		return color.ModeNoColor
	}
	return w.detector.GetMode()
}

func (w *ConsoleWriter) shortFilename(filename string) string {
	parts := strings.Split(filename, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return filename
}

// Close closes the writer (no-op for console)
func (w *ConsoleWriter) Close() error {
	return nil
}
