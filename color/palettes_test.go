package color

import (
	"testing"
)

func TestPaletteCreation(t *testing.T) {
	// Test EXPRESS API
	p1 := BasicPalette()
	if len(p1) == 0 {
		t.Error("BasicPalette should return non-empty palette")
	}

	// Test INSTANTIATED API
	cfg := DefaultPaletteConfig()
	cfg.Name = "test"
	cfg.Colors = map[string]Color{
		"red": RGB(255, 0, 0),
		"green": RGB(0, 255, 0),
	}
	
	p2 := NewPalette(cfg)
	if len(p2) != 2 {
		t.Errorf("NewPalette should create palette with 2 colors, got %d", len(p2))
	}

	// Test FLUENT API
	p3 := NewPaletteWith(
		WithPaletteName("custom"),
		WithColor("blue", RGB(0, 0, 255)),
		WithColor("yellow", RGB(255, 255, 0)),
	)
	if len(p3) != 2 {
		t.Errorf("NewPaletteWith should create palette with 2 colors, got %d", len(p3))
	}
}

func TestPaletteMethods(t *testing.T) {
	p := NewPaletteWith(
		WithColor("red", RGB(255, 0, 0)),
		WithColor("green", RGB(0, 255, 0)),
	)

	// Test Get
	red, exists := p.Get("red")
	if !exists {
		t.Error("Get should find existing color")
	}
	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Error("Get should return correct color")
	}

	_, notExists := p.Get("nonexistent")
	if notExists {
		t.Error("Get should not find non-existent color")
	}

	// Test Set
	p.Set("blue", RGB(0, 0, 255))
	blue, exists := p.Get("blue")
	if !exists || blue.B != 255 {
		t.Error("Set should add new color to palette")
	}

	// Test Names
	names := p.Names()
	if len(names) != 3 {
		t.Errorf("Names should return 3 names, got %d", len(names))
	}

	// Test Merge
	p2 := NewPaletteWith(
		WithColor("yellow", RGB(255, 255, 0)),
		WithColor("cyan", RGB(0, 255, 255)),
	)
	
	merged := p.Merge(p2)
	if len(merged) != 5 {
		t.Errorf("Merge should create palette with 5 colors, got %d", len(merged))
	}
}

func TestPredefinedColors(t *testing.T) {
	// Test basic colors
	if ColorRed.ANSI != 1 {
		t.Error("ColorRed should have ANSI code 1")
	}
	if ColorGreen.ANSI != 2 {
		t.Error("ColorGreen should have ANSI code 2")
	}

	// Test semantic colors
	if ColorSuccess.Name != "success" {
		t.Error("ColorSuccess should have name 'success'")
	}
	if ColorError.Name != "error" {
		t.Error("ColorError should have name 'error'")
	}
}

func TestMaterialColors(t *testing.T) {
	if MaterialRed.Hex != "#F44336" {
		t.Errorf("MaterialRed hex should be #F44336, got %s", MaterialRed.Hex)
	}
	if MaterialBlue.Name != "material_blue" {
		t.Errorf("MaterialBlue name should be material_blue, got %s", MaterialBlue.Name)
	}
}

func TestTailwindColors(t *testing.T) {
	if TailwindRed.Hex != "#EF4444" {
		t.Errorf("TailwindRed hex should be #EF4444, got %s", TailwindRed.Hex)
	}
	if TailwindBlue.Name != "tailwind_blue" {
		t.Errorf("TailwindBlue name should be tailwind_blue, got %s", TailwindBlue.Name)
	}
}

func TestDraculaColors(t *testing.T) {
	if DraculaPurple.Hex != "#BD93F9" {
		t.Errorf("DraculaPurple hex should be #BD93F9, got %s", DraculaPurple.Hex)
	}
	if DraculaGreen.Name != "dracula_green" {
		t.Errorf("DraculaGreen name should be dracula_green, got %s", DraculaGreen.Name)
	}
}

func TestNordColors(t *testing.T) {
	if NordBlue.Hex != "#5E81AC" {
		t.Errorf("NordBlue hex should be #5E81AC, got %s", NordBlue.Hex)
	}
	if NordGreen.Name != "nord_green" {
		t.Errorf("NordGreen name should be nord_green, got %s", NordGreen.Name)
	}
}

func TestGitHubColors(t *testing.T) {
	if GithubGreenLight.Hex != "#28A745" {
		t.Errorf("GithubGreenLight hex should be #28A745, got %s", GithubGreenLight.Hex)
	}
	if GithubBlueLight.Name != "github_blue_light" {
		t.Errorf("GithubBlueLight name should be github_blue_light, got %s", GithubBlueLight.Name)
	}
}

func TestVSCodeColors(t *testing.T) {
	if VSCodeBlue.Hex != "#007ACC" {
		t.Errorf("VSCodeBlue hex should be #007ACC, got %s", VSCodeBlue.Hex)
	}
	if VSCodeGreen.Name != "vscode_green" {
		t.Errorf("VSCodeGreen name should be vscode_green, got %s", VSCodeGreen.Name)
	}
}

func TestPredefinedPalettes(t *testing.T) {
	// Test StatusPalette
	status := StatusPalette()
	if len(status) == 0 {
		t.Error("StatusPalette should not be empty")
	}
	if _, exists := status.Get("success"); !exists {
		t.Error("StatusPalette should contain 'success' color")
	}

	// Test MaterialPalette
	material := MaterialPalette()
	if len(material) != 16 {
		t.Errorf("MaterialPalette should have 16 colors, got %d", len(material))
	}

	// Test DraculaPalette
	dracula := DraculaPalette()
	if len(dracula) != 7 {
		t.Errorf("DraculaPalette should have 7 colors, got %d", len(dracula))
	}

	// Test NordPalette
	nord := NordPalette()
	if len(nord) != 8 {
		t.Errorf("NordPalette should have 8 colors, got %d", len(nord))
	}

	// Test TailwindPalette
	tailwind := TailwindPalette()
	if len(tailwind) != 17 {
		t.Errorf("TailwindPalette should have 17 colors, got %d", len(tailwind))
	}

	// Test GitHubPalette
	github := GitHubPalette()
	if len(github) != 10 {
		t.Errorf("GitHubPalette should have 10 colors, got %d", len(github))
	}

	// Test VSCodePalette
	vscode := VSCodePalette()
	if len(vscode) != 7 {
		t.Errorf("VSCodePalette should have 7 colors, got %d", len(vscode))
	}
}

func TestGetPalette(t *testing.T) {
	// Test existing palette
	material, exists := GetPalette("material")
	if !exists {
		t.Error("GetPalette should find 'material' palette")
	}
	if len(material) == 0 {
		t.Error("Retrieved material palette should not be empty")
	}

	// Test non-existent palette
	_, notExists := GetPalette("nonexistent")
	if notExists {
		t.Error("GetPalette should not find non-existent palette")
	}
}

func TestListPalettes(t *testing.T) {
	names := ListPalettes()
	if len(names) != len(AllPalettes) {
		t.Errorf("ListPalettes should return %d names, got %d", len(AllPalettes), len(names))
	}

	// Check that all expected palettes are present
	expectedPalettes := []string{"status", "material", "dracula", "nord", "tailwind", "github", "vscode"}
	for _, expected := range expectedPalettes {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListPalettes should include '%s'", expected)
		}
	}
}

func TestAllPalettes(t *testing.T) {
	if len(AllPalettes) == 0 {
		t.Error("AllPalettes should not be empty")
	}

	// Test that all palette functions work
	for name, paletteFunc := range AllPalettes {
		palette := paletteFunc()
		if len(palette) == 0 {
			t.Errorf("Palette '%s' should not be empty", name)
		}
	}
}

func TestPaletteColorConsistency(t *testing.T) {
	// Test that colors in palettes match their individual definitions
	material := MaterialPalette()
	
	red, exists := material.Get("red")
	if !exists {
		t.Error("Material palette should contain 'red'")
	}
	if red.Hex != MaterialRed.Hex {
		t.Error("Material palette red should match MaterialRed")
	}

	blue, exists := material.Get("blue")
	if !exists {
		t.Error("Material palette should contain 'blue'")
	}
	if blue.Hex != MaterialBlue.Hex {
		t.Error("Material palette blue should match MaterialBlue")
	}
}

func TestWithColorsOption(t *testing.T) {
	colors := map[string]Color{
		"custom1": RGB(100, 100, 100),
		"custom2": RGB(200, 200, 200),
	}
	
	p := NewPaletteWith(
		WithColors(colors),
		WithColor("custom3", RGB(50, 50, 50)),
	)
	
	if len(p) != 3 {
		t.Errorf("Palette should have 3 colors, got %d", len(p))
	}
	
	c1, exists1 := p.Get("custom1")
	c2, exists2 := p.Get("custom2")
	c3, exists3 := p.Get("custom3")
	
	if !exists1 || !exists2 || !exists3 {
		t.Error("All custom colors should exist in palette")
	}
	
	if c1.R != 100 || c2.R != 200 || c3.R != 50 {
		t.Error("Custom colors should have correct RGB values")
	}
}

func TestDefaultPaletteConfig(t *testing.T) {
	cfg := DefaultPaletteConfig()
	
	if cfg.Name != "default" {
		t.Error("Default palette config should have name 'default'")
	}
	if cfg.Colors == nil {
		t.Error("Default palette config should have initialized Colors map")
	}
	if len(cfg.Colors) != 0 {
		t.Error("Default palette config should start with empty Colors map")
	}
}