package lfx

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/garaekz/lfx/internal/terminal"
)

// ANSI escape codes
const (
	AnsiReset   = "\033[0m"
	AnsiBold    = "\033[1m"
	AnsiRed     = "\033[1;31m"
	AnsiGreen   = "\033[1;32m"
	AnsiYellow  = "\033[1;33m"
	AnsiBlue    = "\033[1;34m"
	AnsiMagenta = "\033[1;35m"
	AnsiCyan    = "\033[1;36m"
	AnsiWhite   = "\033[1;37m"
	BgRed       = "\033[41m"
	BgGreen     = "\033[42m"
	BgYellow    = "\033[43m"
	BgBlue      = "\033[44m"
	BgMagenta   = "\033[45m"
	BgCyan      = "\033[46m"
	BgWhite     = "\033[47m"
)

// Badge width management
var (
	badgeWidth   = 5
	badgeWidthMu sync.Mutex
)

// Global terminal detector
var detector = terminal.NewDetector(os.Stdout)

// badgePad pads the tag to maintain consistent badge width
func badgePad(tag string) string {
	badgeWidthMu.Lock()
	if len(tag) > badgeWidth {
		badgeWidth = len(tag)
	}
	badgeWidthMu.Unlock()
	return tag + strings.Repeat(" ", badgeWidth-len(tag))
}

// BadgeBg prints a badge with both foreground and background colors
func BadgeBg(tag, msg, fg, bg string, args ...any) {
	formattedMsg := fmt.Sprintf(msg, args...)

	if detector.SupportsANSI() {
		fmt.Printf("%s%s[%s]%s %s\n", fg, bg, badgePad(tag), AnsiReset, formattedMsg)
	} else {
		fmt.Printf("[%s] %s\n", badgePad(tag), formattedMsg)
	}
}

// Colorize applies color to text only if terminal supports it
func Colorize(text, color string) string {
	if !detector.SupportsANSI() {
		return text
	}
	return color + text + AnsiReset
}

// Terminal capability functions
func GetTerminalMode() terminal.Mode { return detector.GetMode() }
func SupportsColor() bool            { return detector.SupportsColor() }
func SupportsANSI() bool             { return detector.SupportsANSI() }
func Supports256Color() bool         { return detector.Supports256Color() }
func SupportsTrueColor() bool        { return detector.SupportsTrueColor() }

// Configuration functions
func ForceMode(mode terminal.Mode) { detector.ForceMode(mode) }

// PrintCapabilities shows diagnostic information about terminal support
func PrintCapabilities() {
	mode := detector.GetMode()
	fmt.Printf("LFX Terminal Capabilities:\n")
	fmt.Printf("  Mode: %s\n", mode.String())
	fmt.Printf("  Color Support: %v\n", SupportsColor())
	fmt.Printf("  ANSI Support: %v\n", SupportsANSI())
	fmt.Printf("  256 Color: %v\n", Supports256Color())
	fmt.Printf("  True Color: %v\n", SupportsTrueColor())
	fmt.Printf("  OS: %s\n", detector.GetOS())
	fmt.Printf("  TERM: %s\n", os.Getenv("TERM"))
	fmt.Printf("  COLORTERM: %s\n", os.Getenv("COLORTERM"))

	if SupportsANSI() {
		fmt.Printf("\nColor Test:\n")
		Success("Success message")
		Error("Error message")
		Warn("Warning message")
		Info("Info message")
		Debug("Debug message")
	}
}
