package main

import (
	"github.com/garaekz/lfx"
)

func main() {
	// Show terminal capabilities
	lfx.PrintCapabilities()

	// Basic usage
	lfx.Success("Server started successfully on port %d", 8080)
	lfx.Info("Processing %d requests", 42)
	lfx.Warn("Memory usage is at %d%%", 85)
	lfx.Error("Failed to connect to database: %s", "connection timeout")
	lfx.Debug("Debug info: user session %s", "abc123")

	// Custom badges
	lfx.Badge("API", "Request processed in %dms", lfx.AnsiGreen, 150)
	lfx.Badge("DB", "Query executed", lfx.AnsiBlue)
	lfx.Badge("CACHE", "Hit ratio: %.2f%%", lfx.AnsiYellow, 89.5)

	// Background colors
	lfx.BadgeBg("CRITICAL", "System overload detected", lfx.AnsiWhite, lfx.BgRed)
	// lfx.BadgeBg("TODO", "Implement user authentication", lfx.AnsiBlack, lfx.BgYellow)

	// Conditional coloring
	if lfx.SupportsColor() {
		lfx.Info("Your terminal supports colors! 🎨")
	} else {
		lfx.Info("No color support detected")
	}

	// Manual colorization
	text := lfx.Colorize("This text is colored", lfx.AnsiMagenta)
	lfx.Info("Message: %s", text)
}
