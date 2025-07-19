package main

import (
	"github.com/garaekz/tfx"
)

func main() {
	// Show terminal capabilities
	tfx.PrintCapabilities()

	// Basic usage
	tfx.Success("Server started successfully on port %d", 8080)
	tfx.Info("Processing %d requests", 42)
	tfx.Warn("Memory usage is at %d%%", 85)
	tfx.Error("Failed to connect to database: %s", "connection timeout")
	tfx.Debug("Debug info: user session %s", "abc123")

	// Custom badges
	tfx.Badge("API", "Request processed in %dms", tfx.AnsiGreen, 150)
	tfx.Badge("DB", "Query executed", tfx.AnsiBlue)
	tfx.Badge("CACHE", "Hit ratio: %.2f%%", tfx.AnsiYellow, 89.5)

	// Background colors
	tfx.BadgeBg("CRITICAL", "System overload detected", tfx.AnsiWhite, tfx.BgRed)
	// tfx.BadgeBg("TODO", "Implement user authentication", tfx.AnsiBlack, tfx.BgYellow)

	// Conditional coloring
	if tfx.SupportsColor() {
		tfx.Info("Your terminal supports colors! 🎨")
	} else {
		tfx.Info("No color support detected")
	}

	// Manual colorization
	text := tfx.Colorize("This text is colored", tfx.AnsiMagenta)
	tfx.Info("Message: %s", text)
}
