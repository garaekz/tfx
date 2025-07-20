// Package color provides a multipath API for terminal colors with ANSI, 256-color, and TrueColor support.
package color

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/garaekz/tfx/internal/share"
)

// Color represents a color that can be rendered in different terminal modes
type Color struct {
	R, G, B  uint8  // RGB values (0-255) for truecolor mode
	ANSI     int    // ANSI code for basic mode (0-15)
	Color256 int    // 256-color code (0-255)
	Hex      string // Hex representation (#RRGGBB)
	Name     string // Name for debugging/display
	IsBg     bool   // Whether this is a background color
}

// ColorConfig provides structured configuration for Color creation
type ColorConfig struct {
	R, G, B  uint8
	ANSI     int
	Color256 int
	Hex      string
	Name     string
	IsBg     bool
	Mode     Mode
}

// Mode represents different color rendering modes
type Mode int

const (
	ModeNoColor   Mode = iota // No colors (0)
	ModeANSI                  // Basic 16 colors (1)
	Mode256Color              // 256 colors (2)
	ModeTrueColor             // 24-bit RGB (3)
)

// String returns the string representation of a Mode
func (m Mode) String() string {
	switch m {
	case ModeNoColor:
		return "NoColor"
	case ModeANSI:
		return "ANSI"
	case Mode256Color:
		return "256Color"
	case ModeTrueColor:
		return "TrueColor"
	default:
		return "Unknown"
	}
}

// DefaultColorConfig returns the default configuration for Color
func DefaultColorConfig() ColorConfig {
	return ColorConfig{
		R:    128,
		G:    128,
		B:    128,
		Mode: ModeTrueColor,
		IsBg: false,
	}
}

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick default constructors (renamed to avoid conflicts)
func NewRGB(r, g, b uint8) Color {
	return NewColor(ColorConfig{R: r, G: g, B: b})
}

func NewHex(hex string) Color {
	return NewColor(ColorConfig{Hex: hex})
}

func NewANSI(code int) Color {
	return NewColor(ColorConfig{ANSI: code})
}

func NewColor256(code int) Color {
	return NewColor(ColorConfig{Color256: code})
}

// Backward compatibility - these will be deprecated
func RGB(r, g, b uint8) Color     { return NewRGB(r, g, b) }   // Deprecated: use NewRGB
func Hex(hex string) Color        { return NewHex(hex) }       // Deprecated: use NewHex
func ANSIFunc(code int) Color     { return NewANSI(code) }     // Deprecated: use NewANSI (renamed to avoid conflict)
func Color256Func(code int) Color { return NewColor256(code) } // Deprecated: use NewColor256 (renamed to avoid conflict)

// 2. INSTANTIATED: Config struct
func NewColor(cfg ColorConfig, opts ...share.Option[ColorConfig]) Color {
	// Apply functional options to config
	share.ApplyOptions(&cfg, opts...)

	// If hex is provided, parse it
	if cfg.Hex != "" {
		parseHexIntoConfig(&cfg)
	}

	// Calculate missing values
	if cfg.Hex == "" {
		cfg.Hex = fmt.Sprintf("#%02X%02X%02X", cfg.R, cfg.G, cfg.B)
	}
	if cfg.ANSI == 0 && (cfg.R != 0 || cfg.G != 0 || cfg.B != 0) {
		cfg.ANSI = rgbToANSI(cfg.R, cfg.G, cfg.B)
	}
	if cfg.Color256 == 0 && (cfg.R != 0 || cfg.G != 0 || cfg.B != 0) {
		cfg.Color256 = rgbTo256(cfg.R, cfg.G, cfg.B)
	}

	return Color{
		R:        cfg.R,
		G:        cfg.G,
		B:        cfg.B,
		ANSI:     cfg.ANSI,
		Color256: cfg.Color256,
		Hex:      cfg.Hex,
		Name:     cfg.Name,
		IsBg:     cfg.IsBg,
	}
}

// 3. FLUENT: Functional options + DSL chaining support
func NewColorWith(opts ...share.Option[ColorConfig]) Color {
	cfg := DefaultColorConfig()
	return NewColor(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

func WithRGB(r, g, b uint8) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.R = r
		cfg.G = g
		cfg.B = b
	}
}

func WithHex(hex string) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.Hex = hex
	}
}

func WithANSI(code int) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.ANSI = code
	}
}

func WithColor256(code int) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.Color256 = code
	}
}

func WithName(name string) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.Name = name
	}
}

func WithBackground() share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.IsBg = true
	}
}

func WithMode(mode Mode) share.Option[ColorConfig] {
	return func(cfg *ColorConfig) {
		cfg.Mode = mode
	}
}

// --- COLOR METHODS ---

// Render returns the ANSI escape sequence for the given mode
func (c Color) Render(mode Mode) string {
	if c.IsBg {
		return c.Background(mode)
	}

	switch mode {
	case ModeNoColor:
		return ""
	case ModeANSI:
		if c.ANSI >= 0 && c.ANSI <= 7 {
			return fmt.Sprintf("\033[3%dm", c.ANSI)
		} else if c.ANSI >= 8 && c.ANSI <= 15 {
			return fmt.Sprintf("\033[9%dm", c.ANSI-8)
		}
		return ""
	case Mode256Color:
		return fmt.Sprintf("\033[38;5;%dm", c.Color256)
	case ModeTrueColor:
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
	default:
		return ""
	}
}

// Background returns the background version of this color
func (c Color) Background(mode Mode) string {
	switch mode {
	case ModeNoColor:
		return ""
	case ModeANSI:
		if c.ANSI >= 0 && c.ANSI <= 7 {
			return fmt.Sprintf("\033[4%dm", c.ANSI)
		} else if c.ANSI >= 8 && c.ANSI <= 15 {
			return fmt.Sprintf("\033[10%dm", c.ANSI-8)
		}
		return ""
	case Mode256Color:
		return fmt.Sprintf("\033[48;5;%dm", c.Color256)
	case ModeTrueColor:
		return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
	default:
		return ""
	}
}

// Bg returns a background version of this color
func (c Color) Bg() Color {
	return Color{
		R:        c.R,
		G:        c.G,
		B:        c.B,
		ANSI:     c.ANSI,
		Color256: c.Color256,
		Hex:      c.Hex,
		Name:     c.Name + "_bg",
		IsBg:     true,
	}
}

// WithName returns a copy of the color with a new name
func (c Color) WithName(name string) Color {
	c.Name = name
	return c
}

// String returns a string representation of the color
func (c Color) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.Hex
}

// Apply applies the color to text and adds reset
func (c Color) Apply(text string) string {
	return c.Render(ModeTrueColor) + text + Reset
}

// ApplyMode applies the color to text in the specified mode
func (c Color) ApplyMode(text string, mode Mode) string {
	if mode == ModeNoColor {
		return text
	}
	return c.Render(mode) + text + Reset
}

// --- INTERNAL HELPER FUNCTIONS ---

func parseHexIntoConfig(cfg *ColorConfig) {
	hex := strings.TrimPrefix(cfg.Hex, "#")
	if len(hex) == 3 {
		// Short hex format: #RGB -> #RRGGBB
		hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
	}
	if len(hex) != 6 {
		return // Invalid hex
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	cfg.R = uint8(r)
	cfg.G = uint8(g)
	cfg.B = uint8(b)
	cfg.Hex = "#" + strings.ToUpper(hex)
}

func rgbToANSI(r, g, b uint8) int {
	brightness := (int(r) + int(g) + int(b)) / 3
	maxVal := max(r, g, b)

	if maxVal == 0 {
		return 0 // Black
	}

	bright := brightness > 127

	switch {
	case r == maxVal && g > b:
		if bright {
			return 11
		} else {
			return 3
		} // Yellow
	case r == maxVal:
		if bright {
			return 9
		} else {
			return 1
		} // Red
	case g == maxVal && b > r:
		if bright {
			return 14
		} else {
			return 6
		} // Cyan
	case g == maxVal:
		if bright {
			return 10
		} else {
			return 2
		} // Green
	case b == maxVal && r > g:
		if bright {
			return 13
		} else {
			return 5
		} // Magenta
	case b == maxVal:
		if bright {
			return 12
		} else {
			return 4
		} // Blue
	default:
		if bright {
			return 15
		} else {
			return 7
		} // White/Gray
	}
}

func rgbTo256(r, g, b uint8) int {
	// Check if it's a grayscale color
	if r == g && g == b {
		if r < 8 {
			return 16
		}
		if r > 248 {
			return 231
		}
		return 232 + int(r-8)/10
	}

	// Convert to 6x6x6 color cube
	r6 := int(r) * 5 / 255
	g6 := int(g) * 5 / 255
	b6 := int(b) * 5 / 255

	return 16 + 36*r6 + 6*g6 + b6
}

func ansiToRGB(ansi int) (uint8, uint8, uint8) {
	colors := [][]uint8{
		{0, 0, 0},       // 0: Black
		{128, 0, 0},     // 1: Red
		{0, 128, 0},     // 2: Green
		{128, 128, 0},   // 3: Yellow
		{0, 0, 128},     // 4: Blue
		{128, 0, 128},   // 5: Magenta
		{0, 128, 128},   // 6: Cyan
		{192, 192, 192}, // 7: White
		{128, 128, 128}, // 8: Bright Black (Gray)
		{255, 0, 0},     // 9: Bright Red
		{0, 255, 0},     // 10: Bright Green
		{255, 255, 0},   // 11: Bright Yellow
		{0, 0, 255},     // 12: Bright Blue
		{255, 0, 255},   // 13: Bright Magenta
		{0, 255, 255},   // 14: Bright Cyan
		{255, 255, 255}, // 15: Bright White
	}

	if ansi >= 0 && ansi < len(colors) {
		return colors[ansi][0], colors[ansi][1], colors[ansi][2]
	}
	return 0, 0, 0
}

func color256ToRGB(code int) (uint8, uint8, uint8) {
	if code < 16 {
		return ansiToRGB(code)
	} else if code < 232 {
		code -= 16
		r := code / 36
		g := (code % 36) / 6
		b := code % 6

		rVal := uint8(r * 255 / 5)
		gVal := uint8(g * 255 / 5)
		bVal := uint8(b * 255 / 5)

		return rVal, gVal, bVal
	} else {
		gray := uint8(8 + (code-232)*10)
		return gray, gray, gray
	}
}

func max(a, b, c uint8) uint8 {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}
