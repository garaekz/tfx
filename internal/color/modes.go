package color

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents a color that can be rendered in different modes
type Color struct {
	// RGB values (0-255) for truecolor mode
	R, G, B uint8

	// ANSI code for basic mode (0-15)
	ANSI int

	// 256-color code (0-255)
	Color256 int

	// Hex representation (#RRGGBB)
	Hex string

	// Name for debugging/display
	Name string
}

// Mode represents different color rendering modes
type Mode int

const (
	ModeNoColor   Mode = iota // No colors
	ModeANSI                  // Basic 16 colors
	Mode256Color              // 256 colors
	ModeTrueColor             // 24-bit RGB
)

// NewRGB creates a color from RGB values
func NewRGB(r, g, b uint8) Color {
	return Color{
		R:        r,
		G:        g,
		B:        b,
		ANSI:     rgbToANSI(r, g, b),
		Color256: rgbTo256(r, g, b),
		Hex:      fmt.Sprintf("#%02X%02X%02X", r, g, b),
	}
}

// NewHex creates a color from hex string (#RRGGBB or RRGGBB)
func NewHex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return Color{} // Invalid hex
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return NewRGB(uint8(r), uint8(g), uint8(b))
}

// NewANSI creates a color from ANSI code (0-15)
func NewANSI(code int) Color {
	if code < 0 || code > 15 {
		return Color{}
	}

	// Convert ANSI to approximate RGB
	r, g, b := ansiToRGB(code)
	return Color{
		R:        r,
		G:        g,
		B:        b,
		ANSI:     code,
		Color256: code, // First 16 colors in 256-color mode are same as ANSI
		Hex:      fmt.Sprintf("#%02X%02X%02X", r, g, b),
	}
}

// New256 creates a color from 256-color code (0-255)
func New256(code int) Color {
	if code < 0 || code > 255 {
		return Color{}
	}

	r, g, b := color256ToRGB(code)
	return Color{
		R:        r,
		G:        g,
		B:        b,
		ANSI:     rgbToANSI(r, g, b),
		Color256: code,
		Hex:      fmt.Sprintf("#%02X%02X%02X", r, g, b),
	}
}

// Render returns the ANSI escape sequence for the given mode
func (c Color) Render(mode Mode) string {
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

// String returns a string representation of the color
func (c Color) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.Hex
}

// Convert RGB to closest ANSI color (0-15)
func rgbToANSI(r, g, b uint8) int {
	// Simple mapping to 16 ANSI colors
	// This is a simplified version - could be more sophisticated

	// Convert to intensity
	brightness := (int(r) + int(g) + int(b)) / 3

	// Determine dominant color
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

// Convert RGB to 256-color code
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

// Convert ANSI color to RGB (approximate)
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

// Convert 256-color code to RGB
func color256ToRGB(code int) (uint8, uint8, uint8) {
	if code < 16 {
		// Standard colors
		return ansiToRGB(code)
	} else if code < 232 {
		// 6x6x6 color cube
		code -= 16
		r := code / 36
		g := (code % 36) / 6
		b := code % 6

		// Convert 0-5 to 0-255
		rVal := uint8(r * 255 / 5)
		gVal := uint8(g * 255 / 5)
		bVal := uint8(b * 255 / 5)

		return rVal, gVal, bVal
	} else {
		// Grayscale
		gray := uint8(8 + (code-232)*10)
		return gray, gray, gray
	}
}

// Helper function for max
func max(a, b, c uint8) uint8 {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}
