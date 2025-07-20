package color

import (
	"fmt"
	"strings"

	"github.com/garaekz/tfx/internal/share"
)

// StyleConfig provides structured configuration for text styling
type StyleConfig struct {
	Text       string
	ForeGround Color
	Background Color
	Bold       bool
	Dim        bool
	Italic     bool
	Underline  bool
	Blink      bool
	Reverse    bool
	Strike     bool
	Mode       Mode
}

// DefaultStyleConfig returns default styling configuration
func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		Mode: ModeTrueColor,
	}
}

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick styling functions
func Style(text string, fg Color) string {
	return NewStyle(StyleConfig{Text: text, ForeGround: fg})
}

func StyleBg(text string, fg, bg Color) string {
	return NewStyle(StyleConfig{Text: text, ForeGround: fg, Background: bg})
}

// 2. INSTANTIATED: Config struct
func NewStyle(cfg StyleConfig, opts ...share.Option[StyleConfig]) string {
	share.ApplyOptions(&cfg, opts...)
	return renderStyledText(cfg)
}

// 3. FLUENT: Functional options + DSL chaining support
func NewStyleWith(opts ...share.Option[StyleConfig]) string {
	cfg := DefaultStyleConfig()
	return NewStyle(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

func WithText(text string) share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Text = text
	}
}

func WithForeground(color Color) share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.ForeGround = color
	}
}

func WithBg(color Color) share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Background = color
	}
}

func WithBold() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Bold = true
	}
}

func WithDim() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Dim = true
	}
}

func WithItalic() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Italic = true
	}
}

func WithUnderline() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Underline = true
	}
}

func WithBlink() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Blink = true
	}
}

func WithReverse() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Reverse = true
	}
}

func WithStrike() share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Strike = true
	}
}

func WithStyleMode(mode Mode) share.Option[StyleConfig] {
	return func(cfg *StyleConfig) {
		cfg.Mode = mode
	}
}

// --- UTILITY FUNCTIONS ---

// Apply applies a color to text with optional reset
func Apply(text string, color Color, mode Mode) string {
	if mode == ModeNoColor {
		return text
	}
	return color.Render(mode) + text + Reset
}

// ApplyBg applies foreground and background colors to text
func ApplyBg(text string, fg, bg Color, mode Mode) string {
	if mode == ModeNoColor {
		return text
	}
	return fg.Render(mode) + bg.Background(mode) + text + Reset
}

// Sprint applies color and returns formatted string
func Sprint(color Color, args ...interface{}) string {
	text := fmt.Sprint(args...)
	return Apply(text, color, ModeTrueColor)
}

// Sprintf applies color with formatting
func Sprintf(color Color, format string, args ...interface{}) string {
	text := fmt.Sprintf(format, args...)
	return Apply(text, color, ModeTrueColor)
}

// Combine multiple ANSI codes
func Combine(codes ...string) string {
	return strings.Join(codes, "")
}

// --- ADVANCED STYLING FUNCTIONS ---

// GradientText creates a gradient effect using multiple colors
func GradientText(text string, colors []Color, mode Mode) string {
	if len(colors) == 0 || len(text) == 0 {
		return text
	}
	if len(colors) == 1 {
		return Apply(text, colors[0], mode)
	}

	result := ""
	textLen := len(text)
	colorLen := len(colors)

	for i, char := range text {
		colorIndex := (i * colorLen) / textLen
		if colorIndex >= colorLen {
			colorIndex = colorLen - 1
		}
		result += Apply(string(char), colors[colorIndex], mode)
	}

	return result
}

// RainbowText applies rainbow colors to text
func RainbowText(text string, mode Mode) string {
	rainbowColors := []Color{
		ColorRed, ColorYellow, ColorGreen, ColorCyan, ColorBlue, ColorMagenta,
	}
	return GradientText(text, rainbowColors, mode)
}

// PulseText creates a pulsing effect (simulation with dim)
func PulseText(text string, color Color, pulse bool) string {
	if pulse {
		return Combine(Dim, color.Render(ModeTrueColor)) + text + Reset
	}
	return Apply(text, color, ModeTrueColor)
}

// StripANSI removes ANSI escape sequences from text
func StripANSI(text string) string {
	result := ""
	inEscape := false
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\033' && i+1 < len(runes) && runes[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape && runes[i] == 'm' {
			inEscape = false
			continue
		}
		if !inEscape {
			result += string(runes[i])
		}
	}

	return result
}

// GetLength returns the visual length of text (excluding ANSI codes)
func GetLength(text string) int {
	return len([]rune(StripANSI(text)))
}

// PadString pads a string to a specific visual width
func PadString(text string, width int, padChar rune) string {
	textLen := GetLength(text)
	if textLen >= width {
		return text
	}

	padding := strings.Repeat(string(padChar), width-textLen)
	return text + padding
}

// CenterString centers text within a given width
func CenterString(text string, width int) string {
	textLen := GetLength(text)
	if textLen >= width {
		return text
	}

	leftPad := (width - textLen) / 2
	rightPad := width - textLen - leftPad

	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// --- INTERNAL HELPER FUNCTIONS ---

func renderStyledText(cfg StyleConfig) string {
	if cfg.Text == "" {
		return ""
	}

	if cfg.Mode == ModeNoColor {
		return cfg.Text
	}

	var codes []string

	// Add attributes
	if cfg.Bold {
		codes = append(codes, Bold)
	}
	if cfg.Dim {
		codes = append(codes, Dim)
	}
	if cfg.Italic {
		codes = append(codes, Italic)
	}
	if cfg.Underline {
		codes = append(codes, Underline)
	}
	if cfg.Blink {
		codes = append(codes, Blink)
	}
	if cfg.Reverse {
		codes = append(codes, Reverse)
	}
	if cfg.Strike {
		codes = append(codes, Strike)
	}

	// Add colors
	if cfg.ForeGround.R != 0 || cfg.ForeGround.G != 0 || cfg.ForeGround.B != 0 {
		codes = append(codes, cfg.ForeGround.Render(cfg.Mode))
	}
	if cfg.Background.R != 0 || cfg.Background.G != 0 || cfg.Background.B != 0 {
		codes = append(codes, cfg.Background.Background(cfg.Mode))
	}

	if len(codes) == 0 {
		return cfg.Text
	}

	return Combine(codes...) + cfg.Text + Reset
}

// --- SEMANTIC STYLING FUNCTIONS ---

// Success applies success styling to text
func Success(text string) string {
	return Apply(text, ColorSuccess, ModeTrueColor)
}

// Error applies error styling to text
func Error(text string) string {
	return Apply(text, ColorError, ModeTrueColor)
}

// Warning applies warning styling to text
func Warning(text string) string {
	return Apply(text, ColorWarning, ModeTrueColor)
}

// Info applies info styling to text
func Info(text string) string {
	return Apply(text, ColorInfo, ModeTrueColor)
}

// Debug applies debug styling to text
func Debug(text string) string {
	return Apply(text, ColorDebug, ModeTrueColor)
}

// --- TFX-SPECIFIC FEATURES ---

// Badge creates a styled badge with background
func Badge(text string, fg, bg Color) string {
	return ApplyBg(" "+text+" ", fg, bg, ModeTrueColor)
}

// SuccessBadge creates a success badge
func SuccessBadge(text string) string {
	return Badge(text, ColorBlack, ColorSuccess)
}

// ErrorBadge creates an error badge
func ErrorBadge(text string) string {
	return Badge(text, ColorWhite, ColorError)
}

// WarningBadge creates a warning badge
func WarningBadge(text string) string {
	return Badge(text, ColorBlack, ColorWarning)
}

// InfoBadge creates an info badge
func InfoBadge(text string) string {
	return Badge(text, ColorWhite, ColorInfo)
}

// DebugBadge creates a debug badge
func DebugBadge(text string) string {
	return Badge(text, ColorWhite, ColorDebug)
}

// ProgressBar creates a simple progress bar
func ProgressBar(current, total, width int, filledColor, emptyColor Color) string {
	if total == 0 {
		total = 1
	}

	filled := (current * width) / total
	if filled > width {
		filled = width
	}

	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", width-filled)

	return Apply(filledStr, filledColor, ModeTrueColor) + Apply(emptyStr, emptyColor, ModeTrueColor)
}

// Border creates a bordered text block
func Border(text string, borderColor Color) string {
	lines := strings.Split(text, "\n")
	maxWidth := 0

	// Find the maximum width
	for _, line := range lines {
		if w := GetLength(line); w > maxWidth {
			maxWidth = w
		}
	}

	// Create bordered output
	result := Apply("┌"+strings.Repeat("─", maxWidth+2)+"┐", borderColor, ModeTrueColor) + "\n"

	for _, line := range lines {
		padding := strings.Repeat(" ", maxWidth-GetLength(line))
		result += Apply("│", borderColor, ModeTrueColor) + " " + line + padding + " " + Apply("│", borderColor, ModeTrueColor) + "\n"
	}

	result += Apply("└"+strings.Repeat("─", maxWidth+2)+"┘", borderColor, ModeTrueColor)

	return result
}
