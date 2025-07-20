package color

import (
	"strings"
	"testing"
)

func TestColorCreation(t *testing.T) {
	// Test EXPRESS API
	c1 := RGB(255, 0, 0)
	if c1.R != 255 || c1.G != 0 || c1.B != 0 {
		t.Errorf("RGB(255, 0, 0) failed, got RGB(%d, %d, %d)", c1.R, c1.G, c1.B)
	}

	c2 := Hex("#FF0000")
	if c2.Hex != "#FF0000" {
		t.Errorf("Hex(\"#FF0000\") failed, got %s", c2.Hex)
	}

	c3 := ANSI(1)
	if c3.ANSI != 1 {
		t.Errorf("ANSI(1) failed, got %d", c3.ANSI)
	}

	c4 := Color256(196)
	if c4.Color256 != 196 {
		t.Errorf("Color256(196) failed, got %d", c4.Color256)
	}
}

func TestColorConfigurationAPI(t *testing.T) {
	// Test INSTANTIATED API
	cfg := DefaultColorConfig()
	cfg.R = 255
	cfg.G = 128
	cfg.B = 64
	cfg.Name = "test_color"

	c := NewColor(cfg)
	if c.R != 255 || c.G != 128 || c.B != 64 || c.Name != "test_color" {
		t.Errorf("NewColor failed, got RGB(%d, %d, %d) name=%s", c.R, c.G, c.B, c.Name)
	}
}

func TestColorFluentAPI(t *testing.T) {
	// Test FLUENT API
	c := NewColorWith(
		WithRGB(255, 255, 0),
		WithName("yellow"),
	)

	if c.R != 255 || c.G != 255 || c.B != 0 || c.Name != "yellow" {
		t.Errorf("NewColorWith failed, got RGB(%d, %d, %d) name=%s", c.R, c.G, c.B, c.Name)
	}
}

func TestHexParsing(t *testing.T) {
	tests := []struct {
		input string
		r, g, b uint8
	}{
		{"#FF0000", 255, 0, 0},
		{"FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
		{"#F0F", 255, 0, 255}, // Short hex
	}

	for _, tt := range tests {
		c := Hex(tt.input)
		if c.R != tt.r || c.G != tt.g || c.B != tt.b {
			t.Errorf("Hex(%s) = RGB(%d, %d, %d), want RGB(%d, %d, %d)", 
				tt.input, c.R, c.G, c.B, tt.r, tt.g, tt.b)
		}
	}
}

func TestColorMethods(t *testing.T) {
	c := RGB(255, 128, 64).WithName("orange")

	// Test WithName
	if c.Name != "orange" {
		t.Errorf("WithName failed, got %s", c.Name)
	}

	// Test Bg
	bg := c.Bg()
	if !bg.IsBg {
		t.Error("Bg() should set IsBg to true")
	}
	if !strings.Contains(bg.Name, "_bg") {
		t.Error("Bg() should append _bg to name")
	}

	// Test String
	if c.String() != "orange" {
		t.Errorf("String() with name failed, got %s", c.String())
	}

	unnamed := RGB(255, 255, 255)
	if unnamed.String() != "#FFFFFF" {
		t.Errorf("String() without name failed, got %s", unnamed.String())
	}
}

func TestColorRendering(t *testing.T) {
	c := RGB(255, 0, 0)

	// Test different modes
	noColor := c.Render(ModeNoColor)
	if noColor != "" {
		t.Error("ModeNoColor should return empty string")
	}

	ansi := c.Render(ModeANSI)
	if ansi == "" {
		t.Error("ModeANSI should return non-empty string")
	}

	color256 := c.Render(Mode256Color)
	if !strings.Contains(color256, "38;5;") {
		t.Error("Mode256Color should contain 256-color sequence")
	}

	trueColor := c.Render(ModeTrueColor)
	if !strings.Contains(trueColor, "38;2;255;0;0") {
		t.Error("ModeTrueColor should contain RGB sequence")
	}
}

func TestBackgroundRendering(t *testing.T) {
	c := RGB(255, 0, 0)

	// Test background versions
	bg := c.Background(ModeTrueColor)
	if !strings.Contains(bg, "48;2;255;0;0") {
		t.Error("Background should contain background RGB sequence")
	}

	bg256 := c.Background(Mode256Color)
	if !strings.Contains(bg256, "48;5;") {
		t.Error("Background should contain 256-color background sequence")
	}
}

func TestApplyMethods(t *testing.T) {
	c := RGB(255, 0, 0)
	text := "test"

	// Test Apply
	result := c.Apply(text)
	if !strings.Contains(result, text) {
		t.Error("Apply should contain original text")
	}
	if !strings.Contains(result, "\033[0m") {
		t.Error("Apply should contain reset sequence")
	}

	// Test ApplyMode
	resultMode := c.ApplyMode(text, ModeNoColor)
	if resultMode != text {
		t.Error("ApplyMode with ModeNoColor should return unchanged text")
	}
}

func TestModeString(t *testing.T) {
	tests := []struct {
		mode Mode
		expected string
	}{
		{ModeNoColor, "NoColor"},
		{ModeANSI, "ANSI"},
		{Mode256Color, "256Color"},
		{ModeTrueColor, "TrueColor"},
	}

	for _, tt := range tests {
		if tt.mode.String() != tt.expected {
			t.Errorf("Mode.String() = %s, want %s", tt.mode.String(), tt.expected)
		}
	}
}

func TestRGBToANSI(t *testing.T) {
	tests := []struct {
		r, g, b uint8
		expected int
	}{
		{0, 0, 0, 0},     // Black
		{255, 0, 0, 1},   // Red (not bright by default in algorithm)
		{0, 255, 0, 2},   // Green 
		{0, 0, 255, 4},   // Blue
		{255, 255, 255, 9}, // Bright (based on algorithm)
	}

	for _, tt := range tests {
		result := rgbToANSI(tt.r, tt.g, tt.b)
		if result != tt.expected {
			t.Errorf("rgbToANSI(%d, %d, %d) = %d, want %d", 
				tt.r, tt.g, tt.b, result, tt.expected)
		}
	}
}

func TestRGBTo256(t *testing.T) {
	// Test grayscale
	gray := rgbTo256(128, 128, 128)
	if gray < 232 || gray > 255 {
		t.Errorf("rgbTo256(128, 128, 128) should return grayscale color, got %d", gray)
	}

	// Test color cube
	red := rgbTo256(255, 0, 0)
	if red < 16 || red > 231 {
		t.Errorf("rgbTo256(255, 0, 0) should return color cube color, got %d", red)
	}
}

func TestANSIToRGB(t *testing.T) {
	// Test basic colors
	r, g, b := ansiToRGB(1) // Red
	if r != 128 || g != 0 || b != 0 {
		t.Errorf("ansiToRGB(1) = (%d, %d, %d), want (128, 0, 0)", r, g, b)
	}

	// Test bright colors
	r, g, b = ansiToRGB(9) // Bright Red
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("ansiToRGB(9) = (%d, %d, %d), want (255, 0, 0)", r, g, b)
	}

	// Test invalid color
	r, g, b = ansiToRGB(99)
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("ansiToRGB(99) should return (0, 0, 0) for invalid input, got (%d, %d, %d)", r, g, b)
	}
}

func TestColor256ToRGB(t *testing.T) {
	// Test standard colors (0-15)
	r, g, b := color256ToRGB(1)
	if r != 128 || g != 0 || b != 0 {
		t.Errorf("color256ToRGB(1) = (%d, %d, %d), want (128, 0, 0)", r, g, b)
	}

	// Test grayscale (232-255)
	r, g, b = color256ToRGB(244) // Mid gray
	if r != g || g != b {
		t.Errorf("color256ToRGB(244) should return grayscale, got (%d, %d, %d)", r, g, b)
	}
}

func TestMaxFunction(t *testing.T) {
	tests := []struct {
		a, b, c uint8
		expected uint8
	}{
		{1, 2, 3, 3},
		{3, 2, 1, 3},
		{2, 3, 1, 3},
		{1, 1, 1, 1},
		{255, 0, 128, 255},
	}

	for _, tt := range tests {
		result := max(tt.a, tt.b, tt.c)
		if result != tt.expected {
			t.Errorf("max(%d, %d, %d) = %d, want %d", 
				tt.a, tt.b, tt.c, result, tt.expected)
		}
	}
}

func TestBackgroundColorFlag(t *testing.T) {
	c := RGB(255, 0, 0)
	
	// Test that regular color is not background
	if c.IsBg {
		t.Error("Regular color should not be background")
	}

	// Test background color rendering
	bg := c.Bg()
	result := bg.Render(ModeTrueColor)
	if !strings.Contains(result, "48;2;") {
		t.Error("Background color should use background sequence when rendered")
	}
}

func TestColorWithFunctionalOptions(t *testing.T) {
	c := NewColorWith(
		WithHex("#FF8000"),
		WithName("orange"),
		WithBackground(),
	)

	if c.Hex != "#FF8000" {
		t.Errorf("WithHex failed, got %s", c.Hex)
	}
	if c.Name != "orange" {
		t.Errorf("WithName failed, got %s", c.Name)
	}
	if !c.IsBg {
		t.Error("WithBackground failed, IsBg should be true")
	}
}

func TestDefaultColorConfig(t *testing.T) {
	cfg := DefaultColorConfig()
	
	if cfg.Mode != ModeTrueColor {
		t.Error("Default mode should be ModeTrueColor")
	}
	if cfg.IsBg {
		t.Error("Default should not be background")
	}
	if cfg.R != 128 || cfg.G != 128 || cfg.B != 128 {
		t.Error("Default RGB should be (128, 128, 128)")
	}
}