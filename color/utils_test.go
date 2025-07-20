package color

import (
	"strings"
	"testing"
)

func TestStyleCreation(t *testing.T) {
	// Test EXPRESS API
	result1 := Style("test", ColorRed)
	if !strings.Contains(result1, "test") {
		t.Error("Style should contain the text")
	}

	result2 := StyleBg("test", ColorRed, ColorBlue.Bg())
	if !strings.Contains(result2, "test") {
		t.Error("StyleBg should contain the text")
	}

	// Test INSTANTIATED API
	cfg := DefaultStyleConfig()
	cfg.Text = "hello"
	cfg.ForeGround = ColorGreen
	cfg.Bold = true

	result3 := NewStyle(cfg)
	if !strings.Contains(result3, "hello") {
		t.Error("NewStyle should contain the text")
	}
	if !strings.Contains(result3, Bold) {
		t.Error("NewStyle should apply bold formatting")
	}

	// Test FLUENT API
	result4 := NewStyleWith(
		WithText("world"),
		WithForeground(ColorBlue),
		WithBold(),
		WithUnderline(),
	)
	if !strings.Contains(result4, "world") {
		t.Error("NewStyleWith should contain the text")
	}
}

func TestStyleOptions(t *testing.T) {
	result := NewStyleWith(
		WithText("styled"),
		WithForeground(ColorRed),
		WithBg(ColorBlue),
		WithBold(),
		WithDim(),
		WithItalic(),
		WithUnderline(),
		WithBlink(),
		WithReverse(),
		WithStrike(),
	)

	// Check that all formatting codes are present
	if !strings.Contains(result, Bold) {
		t.Error("Should contain bold")
	}
	if !strings.Contains(result, Italic) {
		t.Error("Should contain italic")
	}
	if !strings.Contains(result, Underline) {
		t.Error("Should contain underline")
	}
	if !strings.Contains(result, Reset) {
		t.Error("Should contain reset")
	}
}

func TestApplyFunctions(t *testing.T) {
	text := "test"
	color := ColorRed

	// Test Apply
	result := Apply(text, color, ModeTrueColor)
	if !strings.Contains(result, text) {
		t.Error("Apply should contain original text")
	}
	if !strings.Contains(result, Reset) {
		t.Error("Apply should contain reset")
	}

	// Test Apply with NoColor mode
	resultNoColor := Apply(text, color, ModeNoColor)
	if resultNoColor != text {
		t.Error("Apply with ModeNoColor should return unchanged text")
	}

	// Test ApplyBg
	resultBg := ApplyBg(text, ColorRed, ColorBlue, ModeTrueColor)
	if !strings.Contains(resultBg, text) {
		t.Error("ApplyBg should contain original text")
	}
	if !strings.Contains(resultBg, Reset) {
		t.Error("ApplyBg should contain reset")
	}
}

func TestSprintFunctions(t *testing.T) {
	// Test Sprint
	result1 := Sprint(ColorRed, "hello", " ", "world")
	if !strings.Contains(result1, "hello world") {
		t.Error("Sprint should format multiple arguments")
	}

	// Test Sprintf
	result2 := Sprintf(ColorBlue, "Hello %s!", "World")
	if !strings.Contains(result2, "Hello World!") {
		t.Error("Sprintf should format with placeholders")
	}
}

func TestCombine(t *testing.T) {
	result := Combine(Bold, Underline, Red)
	expected := Bold + Underline + Red
	if result != expected {
		t.Error("Combine should join all codes")
	}
}

func TestGradientText(t *testing.T) {
	text := "gradient"
	colors := []Color{ColorRed, ColorGreen, ColorBlue}

	// Test normal gradient
	result := GradientText(text, colors, ModeTrueColor)
	stripped := StripANSI(result)
	if !strings.Contains(stripped, "gradient") {
		t.Error("GradientText should contain original text")
	}

	// Test empty colors
	resultEmpty := GradientText(text, []Color{}, ModeTrueColor)
	if resultEmpty != text {
		t.Error("GradientText with empty colors should return unchanged text")
	}

	// Test single color
	resultSingle := GradientText(text, []Color{ColorRed}, ModeTrueColor)
	if !strings.Contains(resultSingle, "gradient") {
		t.Error("GradientText with single color should work")
	}

	// Test empty text
	resultEmptyText := GradientText("", colors, ModeTrueColor)
	if resultEmptyText != "" {
		t.Error("GradientText with empty text should return empty string")
	}
}

func TestRainbowText(t *testing.T) {
	text := "rainbow"
	result := RainbowText(text, ModeTrueColor)
	stripped := StripANSI(result)
	if !strings.Contains(stripped, "rainbow") {
		t.Error("RainbowText should contain original text")
	}
}

func TestPulseText(t *testing.T) {
	text := "pulse"
	color := ColorRed

	// Test with pulse
	resultPulse := PulseText(text, color, true)
	if !strings.Contains(resultPulse, text) {
		t.Error("PulseText should contain original text")
	}
	if !strings.Contains(resultPulse, Dim) {
		t.Error("PulseText with pulse=true should contain dim")
	}

	// Test without pulse
	resultNoPulse := PulseText(text, color, false)
	if !strings.Contains(resultNoPulse, text) {
		t.Error("PulseText should contain original text")
	}
}

func TestStripANSI(t *testing.T) {
	// Test with ANSI codes
	coloredText := Red + "hello" + Reset + " " + Bold + "world" + Reset
	stripped := StripANSI(coloredText)
	expected := "hello world"
	if stripped != expected {
		t.Errorf("StripANSI should remove all ANSI codes, got %q, want %q", stripped, expected)
	}

	// Test with plain text
	plainText := "plain text"
	strippedPlain := StripANSI(plainText)
	if strippedPlain != plainText {
		t.Error("StripANSI should not change plain text")
	}
}

func TestGetLength(t *testing.T) {
	// Test with ANSI codes
	coloredText := Red + "hello" + Reset
	length := GetLength(coloredText)
	if length != 5 {
		t.Errorf("GetLength should return 5 for 'hello', got %d", length)
	}

	// Test with plain text
	plainText := "world"
	lengthPlain := GetLength(plainText)
	if lengthPlain != 5 {
		t.Errorf("GetLength should return 5 for 'world', got %d", lengthPlain)
	}
}

func TestPadString(t *testing.T) {
	text := "test"
	padded := PadString(text, 10, ' ')
	if GetLength(padded) != 10 {
		t.Errorf("PadString should pad to width 10, got length %d", GetLength(padded))
	}
	if !strings.HasPrefix(padded, text) {
		t.Error("PadString should start with original text")
	}

	// Test with text longer than width
	longText := "very long text"
	paddedLong := PadString(longText, 5, ' ')
	if paddedLong != longText {
		t.Error("PadString should not truncate text longer than width")
	}
}

func TestCenterString(t *testing.T) {
	text := "test"
	centered := CenterString(text, 10)
	if GetLength(centered) != 10 {
		t.Errorf("CenterString should pad to width 10, got length %d", GetLength(centered))
	}

	// Test with text longer than width
	longText := "very long text"
	centeredLong := CenterString(longText, 5)
	if centeredLong != longText {
		t.Error("CenterString should not truncate text longer than width")
	}
}

func TestSemanticFunctions(t *testing.T) {
	text := "message"

	// Test all semantic functions
	successResult := Success(text)
	if !strings.Contains(successResult, text) {
		t.Error("Success should contain original text")
	}

	errorResult := Error(text)
	if !strings.Contains(errorResult, text) {
		t.Error("Error should contain original text")
	}

	warningResult := Warning(text)
	if !strings.Contains(warningResult, text) {
		t.Error("Warning should contain original text")
	}

	infoResult := Info(text)
	if !strings.Contains(infoResult, text) {
		t.Error("Info should contain original text")
	}

	debugResult := Debug(text)
	if !strings.Contains(debugResult, text) {
		t.Error("Debug should contain original text")
	}
}

func TestBadgeFunctions(t *testing.T) {
	text := "TEST"

	// Test generic Badge
	badge := Badge(text, ColorWhite, ColorRed)
	if !strings.Contains(badge, text) {
		t.Error("Badge should contain original text")
	}
	if !strings.Contains(badge, " "+text+" ") {
		t.Error("Badge should add padding around text")
	}

	// Test semantic badges
	successBadge := SuccessBadge(text)
	if !strings.Contains(successBadge, text) {
		t.Error("SuccessBadge should contain original text")
	}

	errorBadge := ErrorBadge(text)
	if !strings.Contains(errorBadge, text) {
		t.Error("ErrorBadge should contain original text")
	}

	warningBadge := WarningBadge(text)
	if !strings.Contains(warningBadge, text) {
		t.Error("WarningBadge should contain original text")
	}

	infoBadge := InfoBadge(text)
	if !strings.Contains(infoBadge, text) {
		t.Error("InfoBadge should contain original text")
	}

	debugBadge := DebugBadge(text)
	if !strings.Contains(debugBadge, text) {
		t.Error("DebugBadge should contain original text")
	}
}

func TestProgressBar(t *testing.T) {
	// Test normal progress
	bar := ProgressBar(50, 100, 20, ColorGreen, ColorRed)
	stripped := StripANSI(bar)
	actualLen := GetLength(stripped)
	if actualLen != 20 {
		t.Errorf("ProgressBar should have correct width: expected 20, got %d (stripped: %q)", actualLen, stripped)
	}

	// Test zero total
	barZero := ProgressBar(10, 0, 20, ColorGreen, ColorRed)
	strippedZero := StripANSI(barZero)
	actualZeroLen := GetLength(strippedZero)
	if actualZeroLen != 20 {
		t.Errorf("ProgressBar with zero total should still work: expected 20, got %d (stripped: %q)", actualZeroLen, strippedZero)
	}

	// Test overflow
	barOverflow := ProgressBar(150, 100, 20, ColorGreen, ColorRed)
	strippedOverflow := StripANSI(barOverflow)
	actualOverflowLen := GetLength(strippedOverflow)
	if actualOverflowLen != 20 {
		t.Errorf("ProgressBar with overflow should be clamped: expected 20, got %d (stripped: %q)", actualOverflowLen, strippedOverflow)
	}
}

func TestBorder(t *testing.T) {
	text := "hello\nworld"
	bordered := Border(text, ColorBlue)

	lines := strings.Split(bordered, "\n")
	if len(lines) < 4 {
		t.Error("Border should create at least 4 lines (top, content lines, bottom)")
	}

	// Check for border characters
	if !strings.Contains(bordered, "┌") || !strings.Contains(bordered, "┐") {
		t.Error("Border should contain top border characters")
	}
	if !strings.Contains(bordered, "└") || !strings.Contains(bordered, "┘") {
		t.Error("Border should contain bottom border characters")
	}
	if !strings.Contains(bordered, "│") {
		t.Error("Border should contain side border characters")
	}
}

func TestDefaultStyleConfig(t *testing.T) {
	cfg := DefaultStyleConfig()

	if cfg.Mode != ModeTrueColor {
		t.Error("Default style mode should be ModeTrueColor")
	}
	if cfg.Text != "" {
		t.Error("Default style text should be empty")
	}
	if cfg.Bold || cfg.Italic || cfg.Underline {
		t.Error("Default style should have no formatting enabled")
	}
}

func TestStyleWithNoColor(t *testing.T) {
	cfg := DefaultStyleConfig()
	cfg.Text = "test"
	cfg.Mode = ModeNoColor
	cfg.ForeGround = ColorRed
	cfg.Bold = true

	result := NewStyle(cfg)
	if result != "test" {
		t.Error("Style with ModeNoColor should return plain text")
	}
}

func TestStyleWithEmptyText(t *testing.T) {
	cfg := DefaultStyleConfig()
	cfg.Text = ""
	cfg.ForeGround = ColorRed

	result := NewStyle(cfg)
	if result != "" {
		t.Error("Style with empty text should return empty string")
	}
}

func TestRenderStyledTextWithoutColors(t *testing.T) {
	cfg := StyleConfig{
		Text: "test",
		Bold: true,
		Mode: ModeTrueColor,
		// No colors set (zero values)
	}

	result := renderStyledText(cfg)
	if !strings.Contains(result, Bold) {
		t.Error("Should contain bold even without colors")
	}
	if !strings.Contains(result, "test") {
		t.Error("Should contain the text")
	}
}
