package progress

import (
	"bytes"
	"testing"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
)

func TestProgressThemeDefaults(t *testing.T) {
	tests := []struct {
		name  string
		theme ProgressTheme
	}{
		{"material", MaterialTheme},
		{"dracula", DraculaTheme},
		{"nord", NordTheme},
		{"github", GitHubTheme},
		{"tailwind", TailwindTheme},
		{"vscode", VSCodeTheme},
		{"rainbow", RainbowTheme},
		{"neon", NeonTheme},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.theme.Name != tt.name {
				t.Errorf("expected theme name '%s', got '%s'", tt.name, tt.theme.Name)
			}

			// Check that all color fields are set (not zero values)
			if tt.theme.CompleteColor == (color.Color{}) {
				t.Error("CompleteColor should not be zero value")
			}
			if tt.theme.IncompleteColor == (color.Color{}) {
				t.Error("IncompleteColor should not be zero value")
			}
			if tt.theme.LabelColor == (color.Color{}) {
				t.Error("LabelColor should not be zero value")
			}
			if tt.theme.PercentColor == (color.Color{}) {
				t.Error("PercentColor should not be zero value")
			}
			if tt.theme.BorderColor == (color.Color{}) {
				t.Error("BorderColor should not be zero value")
			}
		})
	}
}

func TestThemeEffectSettings(t *testing.T) {
	// Themes that should have effects enabled
	effectEnabledThemes := []ProgressTheme{DraculaTheme, RainbowTheme, NeonTheme}

	for _, theme := range effectEnabledThemes {
		if !theme.EffectEnabled {
			t.Errorf("theme %s should have EffectEnabled=true", theme.Name)
		}
	}

	// Themes that should not have effects enabled
	effectDisabledThemes := []ProgressTheme{MaterialTheme, NordTheme, GitHubTheme, TailwindTheme, VSCodeTheme}

	for _, theme := range effectDisabledThemes {
		if theme.EffectEnabled {
			t.Errorf("theme %s should have EffectEnabled=false", theme.Name)
		}
	}
}

func TestAllProgressThemes(t *testing.T) {
	expectedCount := 8 // Material, Dracula, Nord, GitHub, Tailwind, VSCode, Rainbow, Neon

	if len(AllProgressThemes) != expectedCount {
		t.Errorf("expected %d themes in AllProgressThemes, got %d", expectedCount, len(AllProgressThemes))
	}

	// Check that all themes are present
	themeNames := make(map[string]bool)
	for _, theme := range AllProgressThemes {
		themeNames[theme.Name] = true
	}

	expectedNames := []string{"material", "dracula", "nord", "github", "tailwind", "vscode", "rainbow", "neon"}
	for _, name := range expectedNames {
		if !themeNames[name] {
			t.Errorf("expected theme '%s' in AllProgressThemes", name)
		}
	}
}

func TestThemeRenderColor(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	result := theme.RenderColor(color.MaterialRed, detector)

	// Should return a string (may be empty in test environment)
	_ = result // Color rendering depends on terminal capabilities
}

func TestThemeRenderColorNilDetector(t *testing.T) {
	theme := MaterialTheme

	result := theme.RenderColor(color.MaterialRed, nil)

	// Should return a string even with nil detector (fallback)
	if len(result) == 0 {
		t.Error("expected non-empty color rendering with nil detector")
	}
}

func TestProgressEffectTypes(t *testing.T) {
	// Test that effect constants are defined
	effects := []ProgressEffect{
		EffectNone,
		EffectGradient,
		EffectRainbow,
		EffectPulse,
		EffectGlow,
	}

	// Check that they have different values
	for i, effect1 := range effects {
		for j, effect2 := range effects {
			if i != j && effect1 == effect2 {
				t.Errorf("effects at index %d and %d have same value %d", i, j, effect1)
			}
		}
	}
}

func TestThemeRenderProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	tests := []struct {
		name    string
		percent float64
		width   int
		effect  ProgressEffect
	}{
		{"Solid None", 0.5, 10, EffectNone},
		{"Solid Rainbow", 0.5, 10, EffectRainbow},
		{"Solid Gradient", 0.5, 10, EffectGradient},
		{"Solid Glow", 0.5, 10, EffectGlow},
		{"Empty", 0.0, 10, EffectNone},
		{"Full", 1.0, 10, EffectNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := theme.RenderProgress(tt.percent, tt.width, tt.effect, detector)

			if len(result) == 0 {
				t.Error("expected non-empty progress rendering")
			}

			// Result should contain some characters (█ or ░)
			if !containsProgressChars(result) {
				t.Errorf("expected progress characters in result, got %q", result)
			}
		})
	}
}

func TestThemeRenderSolidProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	// Test different fill levels
	tests := []struct {
		filled int
		width  int
	}{
		{0, 10},  // Empty
		{5, 10},  // Half
		{10, 10}, // Full
	}

	for _, tt := range tests {
		result := theme.renderSolidProgress(tt.filled, tt.width, detector)

		if len(result) == 0 {
			t.Error("expected non-empty solid progress rendering")
		}

		// Count progress characters (█ and ░) in result
		progressCount := countProgressChars(result)
		if progressCount != tt.width {
			t.Errorf("expected %d progress characters, got %d", tt.width, progressCount)
		}
	}
}

func TestThemeRenderRainbowProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	result := theme.renderRainbowProgress(5, 10, detector)

	if len(result) == 0 {
		t.Error("expected non-empty rainbow progress rendering")
	}

	// Should contain progress characters
	if !containsProgressChars(result) {
		t.Errorf("expected progress characters in rainbow result, got %q", result)
	}
}

func TestThemeRenderGradientProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	result := theme.renderGradientProgress(5, 10, detector)

	if len(result) == 0 {
		t.Error("expected non-empty gradient progress rendering")
	}

	// Should contain progress characters
	if !containsProgressChars(result) {
		t.Errorf("expected progress characters in gradient result, got %q", result)
	}
}

func TestThemeRenderGlowProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	detector := terminal.NewDetector(buf)
	theme := MaterialTheme

	result := theme.renderGlowProgress(5, 10, detector)

	if len(result) == 0 {
		t.Error("expected non-empty glow progress rendering")
	}

	// Should contain progress characters
	if !containsProgressChars(result) {
		t.Errorf("expected progress characters in glow result, got %q", result)
	}
}

func TestGetThemeByName(t *testing.T) {
	tests := []struct {
		name     string
		expected ProgressTheme
		found    bool
	}{
		{"material", MaterialTheme, true},
		{"dracula", DraculaTheme, true},
		{"nord", NordTheme, true},
		{"github", GitHubTheme, true},
		{"nonexistent", MaterialTheme, false}, // Should return MaterialTheme as fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, found := GetThemeByName(tt.name)

			if found != tt.found {
				t.Errorf("expected found=%v, got %v", tt.found, found)
			}

			if theme.Name != tt.expected.Name {
				t.Errorf("expected theme name '%s', got '%s'", tt.expected.Name, theme.Name)
			}
		})
	}
}

func TestCreateCustomTheme(t *testing.T) {
	customTheme := CreateCustomTheme(
		"custom",
		color.MaterialGreen,
		color.MaterialRed,
		color.MaterialBlue,
	)

	if customTheme.Name != "custom" {
		t.Errorf("expected name 'custom', got '%s'", customTheme.Name)
	}

	if customTheme.CompleteColor != color.MaterialGreen {
		t.Error("expected CompleteColor to be MaterialGreen")
	}

	if customTheme.IncompleteColor != color.MaterialRed {
		t.Error("expected IncompleteColor to be MaterialRed")
	}

	if customTheme.LabelColor != color.MaterialBlue {
		t.Error("expected LabelColor to be MaterialBlue")
	}

	if customTheme.PercentColor != color.MaterialBlue {
		t.Error("expected PercentColor to be MaterialBlue (same as LabelColor)")
	}

	if customTheme.EffectEnabled {
		t.Error("expected EffectEnabled to be false for custom theme")
	}
}

func TestNewThemeFromPalette(t *testing.T) {
	// Create a simple palette for testing
	palette := color.Palette{
		"green": color.MaterialGreen,
		"gray":  color.NewANSI(8), // Gray
		"blue":  color.MaterialBlue,
	}

	theme := NewThemeFromPalette("test", palette)

	if theme.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", theme.Name)
	}

	if theme.CompleteColor != color.MaterialGreen {
		t.Error("expected CompleteColor to be MaterialGreen from palette")
	}

	if theme.IncompleteColor != color.NewANSI(8) {
		t.Error("expected IncompleteColor to be gray from palette")
	}

	if theme.LabelColor != color.MaterialBlue {
		t.Error("expected LabelColor to be MaterialBlue from palette")
	}
}

func TestNewThemeFromPaletteWithMissingColors(t *testing.T) {
	// Create a palette with missing colors
	palette := color.Palette{
		"red": color.MaterialRed, // No green, gray, or blue
	}

	theme := NewThemeFromPalette("incomplete", palette)

	if theme.Name != "incomplete" {
		t.Errorf("expected name 'incomplete', got '%s'", theme.Name)
	}

	// Should use fallback colors
	if theme.CompleteColor == (color.Color{}) {
		t.Error("expected fallback CompleteColor to be set")
	}

	if theme.LabelColor == (color.Color{}) {
		t.Error("expected fallback LabelColor to be set")
	}
}

func TestAbsFunction(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1, 1},
		{10, 10},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

// Helper functions for testing
func containsProgressChars(s string) bool {
	return len(s) > 0 && (containsChar(s, '█') || containsChar(s, '░'))
}

func containsChar(s string, char rune) bool {
	for _, c := range s {
		if c == char {
			return true
		}
	}
	return false
}

func countProgressChars(s string) int {
	// Count █ and ░ characters in string
	count := 0
	for _, c := range s {
		if c == '█' || c == '░' {
			count++
		}
	}
	return count
}
