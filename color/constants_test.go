package color

import (
	"strings"
	"testing"
)

func TestResetConstant(t *testing.T) {
	if Reset != "\033[0m" {
		t.Errorf("Reset should be \\033[0m, got %q", Reset)
	}
}

func TestTextAttributes(t *testing.T) {
	attributes := map[string]string{
		"Bold":      Bold,
		"Dim":       Dim,
		"Italic":    Italic,
		"Underline": Underline,
		"Blink":     Blink,
		"Reverse":   Reverse,
		"Strike":    Strike,
	}

	for name, attr := range attributes {
		if !strings.HasPrefix(attr, "\033[") {
			t.Errorf("%s should start with \\033[, got %q", name, attr)
		}
		if !strings.HasSuffix(attr, "m") {
			t.Errorf("%s should end with m, got %q", name, attr)
		}
	}
}

func TestBasicForegroundColors(t *testing.T) {
	colors := map[string]string{
		"Black":   Black,
		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
	}

	for name, color := range colors {
		if !strings.HasPrefix(color, "\033[3") {
			t.Errorf("%s should start with \\033[3, got %q", name, color)
		}
		if !strings.HasSuffix(color, "m") {
			t.Errorf("%s should end with m, got %q", name, color)
		}
	}
}

func TestBrightForegroundColors(t *testing.T) {
	brightColors := map[string]string{
		"BrightBlack":   BrightBlack,
		"BrightRed":     BrightRed,
		"BrightGreen":   BrightGreen,
		"BrightYellow":  BrightYellow,
		"BrightBlue":    BrightBlue,
		"BrightMagenta": BrightMagenta,
		"BrightCyan":    BrightCyan,
		"BrightWhite":   BrightWhite,
	}

	for name, color := range brightColors {
		if !strings.HasPrefix(color, "\033[9") {
			t.Errorf("%s should start with \\033[9, got %q", name, color)
		}
		if !strings.HasSuffix(color, "m") {
			t.Errorf("%s should end with m, got %q", name, color)
		}
	}
}

func TestBasicBackgroundColors(t *testing.T) {
	bgColors := map[string]string{
		"BgBlack":   BgBlack,
		"BgRed":     BgRed,
		"BgGreen":   BgGreen,
		"BgYellow":  BgYellow,
		"BgBlue":    BgBlue,
		"BgMagenta": BgMagenta,
		"BgCyan":    BgCyan,
		"BgWhite":   BgWhite,
	}

	for name, color := range bgColors {
		if !strings.HasPrefix(color, "\033[4") {
			t.Errorf("%s should start with \\033[4, got %q", name, color)
		}
		if !strings.HasSuffix(color, "m") {
			t.Errorf("%s should end with m, got %q", name, color)
		}
	}
}

func TestBrightBackgroundColors(t *testing.T) {
	brightBgColors := map[string]string{
		"BgBrightBlack":   BgBrightBlack,
		"BgBrightRed":     BgBrightRed,
		"BgBrightGreen":   BgBrightGreen,
		"BgBrightYellow":  BgBrightYellow,
		"BgBrightBlue":    BgBrightBlue,
		"BgBrightMagenta": BgBrightMagenta,
		"BgBrightCyan":    BgBrightCyan,
		"BgBrightWhite":   BgBrightWhite,
	}

	for name, color := range brightBgColors {
		if !strings.HasPrefix(color, "\033[10") {
			t.Errorf("%s should start with \\033[10, got %q", name, color)
		}
		if !strings.HasSuffix(color, "m") {
			t.Errorf("%s should end with m, got %q", name, color)
		}
	}
}

func TestCursorControlConstants(t *testing.T) {
	cursorControls := map[string]string{
		"CursorUp":    CursorUp,
		"CursorDown":  CursorDown,
		"CursorRight": CursorRight,
		"CursorLeft":  CursorLeft,
		"CursorHome":  CursorHome,
	}

	for name, control := range cursorControls {
		if !strings.HasPrefix(control, "\033[") {
			t.Errorf("%s should start with \\033[, got %q", name, control)
		}
	}
}

func TestScreenControlConstants(t *testing.T) {
	screenControls := map[string]string{
		"ClearScreen":     ClearScreen,
		"ClearLine":       ClearLine,
		"ClearToEnd":      ClearToEnd,
		"ClearToStart":    ClearToStart,
		"SaveCursor":      SaveCursor,
		"RestoreCursor":   RestoreCursor,
		"HideCursor":      HideCursor,
		"ShowCursor":      ShowCursor,
	}

	for name, control := range screenControls {
		if !strings.HasPrefix(control, "\033[") {
			t.Errorf("%s should start with \\033[, got %q", name, control)
		}
	}
}

func TestConstantsUniqueness(t *testing.T) {
	// Test that different constants have different values
	if Red == Blue {
		t.Error("Red and Blue should have different values")
	}
	if Bold == Italic {
		t.Error("Bold and Italic should have different values")
	}
	if BgRed == BgBlue {
		t.Error("BgRed and BgBlue should have different values")
	}
}

func TestConstantsAreNotEmpty(t *testing.T) {
	constants := []struct {
		name  string
		value string
	}{
		{"Reset", Reset},
		{"Bold", Bold},
		{"Red", Red},
		{"BgRed", BgRed},
		{"BrightRed", BrightRed},
		{"CursorUp", CursorUp},
		{"ClearScreen", ClearScreen},
	}

	for _, c := range constants {
		if c.value == "" {
			t.Errorf("Constant %s should not be empty", c.name)
		}
	}
}

func TestANSIColorCodes(t *testing.T) {
	// Test specific ANSI color codes
	expectedCodes := map[string]string{
		Black:   "\033[30m",
		Red:     "\033[31m",
		Green:   "\033[32m",
		Yellow:  "\033[33m",
		Blue:    "\033[34m",
		Magenta: "\033[35m",
		Cyan:    "\033[36m",
		White:   "\033[37m",
	}

	for color, expected := range expectedCodes {
		if color != expected {
			t.Errorf("Expected %q, got %q", expected, color)
		}
	}
}

func TestBrightANSIColorCodes(t *testing.T) {
	// Test specific bright ANSI color codes
	expectedBrightCodes := map[string]string{
		BrightBlack:   "\033[90m",
		BrightRed:     "\033[91m",
		BrightGreen:   "\033[92m",
		BrightYellow:  "\033[93m",
		BrightBlue:    "\033[94m",
		BrightMagenta: "\033[95m",
		BrightCyan:    "\033[96m",
		BrightWhite:   "\033[97m",
	}

	for color, expected := range expectedBrightCodes {
		if color != expected {
			t.Errorf("Expected %q, got %q", expected, color)
		}
	}
}

func TestBackgroundANSIColorCodes(t *testing.T) {
	// Test specific background ANSI color codes
	expectedBgCodes := map[string]string{
		BgBlack:   "\033[40m",
		BgRed:     "\033[41m",
		BgGreen:   "\033[42m",
		BgYellow:  "\033[43m",
		BgBlue:    "\033[44m",
		BgMagenta: "\033[45m",
		BgCyan:    "\033[46m",
		BgWhite:   "\033[47m",
	}

	for color, expected := range expectedBgCodes {
		if color != expected {
			t.Errorf("Expected %q, got %q", expected, color)
		}
	}
}

func TestTextAttributeCodes(t *testing.T) {
	// Test specific text attribute codes
	expectedAttrCodes := map[string]string{
		Bold:      "\033[1m",
		Dim:       "\033[2m",
		Italic:    "\033[3m",
		Underline: "\033[4m",
		Blink:     "\033[5m",
		Reverse:   "\033[7m",
		Strike:    "\033[9m",
	}

	for attr, expected := range expectedAttrCodes {
		if attr != expected {
			t.Errorf("Expected %q, got %q", expected, attr)
		}
	}
}