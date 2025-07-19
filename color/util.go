package color

import (
	"fmt"
	"strings"
)

// ColorFunc represents a function that applies color to text
type ColorFunc func(text string) string

// ApplyColor applies a color to text and resets it
func ApplyColor(text, color string) string {
	if color == "" {
		return text
	}
	return color + text + Reset
}

// Colorize is an alias for ApplyColor for backwards compatibility
func Colorize(text, color string) string {
	return ApplyColor(text, color)
}

// Sprint applies color and returns the colored string
func Sprint(color string, args ...interface{}) string {
	text := fmt.Sprint(args...)
	return ApplyColor(text, color)
}

// Sprintf applies color with formatting
func Sprintf(color, format string, args ...interface{}) string {
	text := fmt.Sprintf(format, args...)
	return ApplyColor(text, color)
}

// Combine multiple color codes (e.g., Bold + Red)
func Combine(colors ...string) string {
	return strings.Join(colors, "")
}

// Color factory functions for easy use
func MakeColorFunc(color string) ColorFunc {
	return func(text string) ColorFunc {
		return func(text string) string {
			return ApplyColor(text, color)
		}
	}(color)
}

// Predefined color functions for easy use
var (
	RedText     = MakeColorFunc(Red)
	GreenText   = MakeColorFunc(Green)
	YellowText  = MakeColorFunc(Yellow)
	BlueText    = MakeColorFunc(Blue)
	MagentaText = MakeColorFunc(Magenta)
	CyanText    = MakeColorFunc(Cyan)
	WhiteText   = MakeColorFunc(White)

	// Bright variants
	BrightRedText     = MakeColorFunc(BrightRed)
	BrightGreenText   = MakeColorFunc(BrightGreen)
	BrightYellowText  = MakeColorFunc(BrightYellow)
	BrightBlueText    = MakeColorFunc(BrightBlue)
	BrightMagentaText = MakeColorFunc(BrightMagenta)
	BrightCyanText    = MakeColorFunc(BrightCyan)

	// Semantic functions
	SuccessText = MakeColorFunc(Success)
	ErrorText   = MakeColorFunc(Error)
	WarningText = MakeColorFunc(Warning)
	InfoText    = MakeColorFunc(Info)
	DebugText   = MakeColorFunc(Debug)
)

// Gradient simulation using ANSI colors (creates visual effects)
func GradientText(text string, startColor, endColor string) string {
	// Simple gradient simulation by alternating colors
	// For true gradients, you'd need 256-color or truecolor
	if len(text) <= 1 {
		return ApplyColor(text, startColor)
	}

	result := ""
	mid := len(text) / 2

	for i, char := range text {
		if i < mid {
			result += ApplyColor(string(char), startColor)
		} else {
			result += ApplyColor(string(char), endColor)
		}
	}
	return result
}

// Rainbow text using ANSI colors
func RainbowText(text string) string {
	colors := []string{Red, Yellow, Green, Cyan, Blue, Magenta}
	result := ""

	for i, char := range text {
		color := colors[i%len(colors)]
		result += ApplyColor(string(char), color)
	}
	return result
}

// Strip ANSI color codes from text
func StripColor(text string) string {
	// Simple regex would be better, but avoiding external deps
	result := ""
	inEscape := false

	for i := 0; i < len(text); i++ {
		if text[i] == '\033' && i+1 < len(text) && text[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape && text[i] == 'm' {
			inEscape = false
			continue
		}
		if !inEscape {
			result += string(text[i])
		}
	}
	return result
}

// GetTheme returns a badge theme by name
func GetTheme(name string) (BadgeTheme, bool) {
	theme, exists := BadgeThemes[name]
	return theme, exists
}

// ListThemes returns all available theme names
func ListThemes() []string {
	themes := make([]string, 0, len(BadgeThemes))
	for name := range BadgeThemes {
		themes = append(themes, name)
	}
	return themes
}
