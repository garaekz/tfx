package color

import "fmt"

// ShowPalette displays all available colors
func ShowPalette() {
	fmt.Println("🎨 LFX Color Palette")

	// Basic Colors
	fmt.Println("📋 Basic ANSI Colors:")
	basicColors := map[string]string{
		"Black":   Black,
		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
	}

	for name, color := range basicColors {
		fmt.Printf("  %s %-8s %s\n",
			ApplyColor("●", color),
			name,
			ApplyColor("Sample text", color))
	}

	// Bright Colors
	fmt.Println("\n✨ Bright Colors:")
	brightColors := map[string]string{
		"BrightRed":     BrightRed,
		"BrightGreen":   BrightGreen,
		"BrightYellow":  BrightYellow,
		"BrightBlue":    BrightBlue,
		"BrightMagenta": BrightMagenta,
		"BrightCyan":    BrightCyan,
		"BrightWhite":   BrightWhite,
		"BrightBlack":   BrightBlack,
	}

	for name, color := range brightColors {
		fmt.Printf("  %s %-13s %s\n",
			ApplyColor("●", color),
			name,
			ApplyColor("Sample text", color))
	}

	// Semantic Colors
	fmt.Println("\n🏷️  Semantic Colors:")
	semanticColors := map[string]string{
		"Success": Success,
		"Error":   Error,
		"Warning": Warning,
		"Info":    Info,
		"Debug":   Debug,
	}

	for name, color := range semanticColors {
		fmt.Printf("  %s %-8s %s\n",
			ApplyColor("●", color),
			name,
			ApplyColor("Sample text", color))
	}

	// Text Attributes
	fmt.Println("\n🎭 Text Attributes:")
	attributes := map[string]string{
		"Bold":      Bold + "Bold text" + Reset,
		"Dim":       Dim + "Dim text" + Reset,
		"Italic":    Italic + "Italic text" + Reset,
		"Underline": Underline + "Underlined text" + Reset,
		"Strike":    Strike + "Strikethrough text" + Reset,
	}

	for name, styled := range attributes {
		fmt.Printf("  %-10s %s\n", name+":", styled)
	}

	// Combined Effects
	fmt.Println("\n🌈 Combined Effects:")
	fmt.Printf("  %s\n", ApplyColor("Bold Red", Combine(Bold, Red)))
	fmt.Printf("  %s\n", ApplyColor("Underlined Blue", Combine(Underline, Blue)))
	fmt.Printf("  %s\n", RainbowText("Rainbow Text"))
	fmt.Printf("  %s\n", GradientText("Gradient Effect", Red, Blue))
}

// ShowThemes displays all available badge themes
func ShowThemes() {
	fmt.Println("🎨 Badge Themes")

	for themeName, theme := range BadgeThemes {
		fmt.Printf("📦 %s theme:\n", themeName)
		fmt.Printf("  Success: %s\n", ApplyColor("[OK] Success message", theme.Success))
		fmt.Printf("  Error:   %s\n", ApplyColor("[ERR] Error message", theme.Error))
		fmt.Printf("  Warning: %s\n", ApplyColor("[WARN] Warning message", theme.Warning))
		fmt.Printf("  Info:    %s\n", ApplyColor("[INFO] Info message", theme.Info))
		fmt.Printf("  Debug:   %s\n", ApplyColor("[DBG] Debug message", theme.Debug))
		fmt.Println()
	}
}

// ShowEffects demonstrates special effects
func ShowEffects() {
	fmt.Println("✨ Special Effects")

	// Background combinations
	fmt.Println("🎯 Background Effects:")
	fmt.Printf("  %s\n", ApplyColor(" ERROR ", Combine(BrightWhite, BgRed)))
	fmt.Printf("  %s\n", ApplyColor(" SUCCESS ", Combine(BrightWhite, BgGreen)))
	fmt.Printf("  %s\n", ApplyColor(" WARNING ", Combine(Black, BgYellow)))
	fmt.Printf("  %s\n", ApplyColor(" INFO ", Combine(BrightWhite, BgBlue)))

	// Progress bar simulation
	fmt.Println("\n📊 Progress Effects:")
	completed := ApplyColor("████████", ProgressComplete)
	incomplete := ApplyColor("░░░░", ProgressIncomplete)
	fmt.Printf("  [%s%s] 67%%\n", completed, incomplete)

	// Border effects
	fmt.Println("\n📐 Border Effects:")
	fmt.Printf("  %s\n", ApplyColor("┌─────────────────────┐", Border))
	fmt.Printf("  %s %s %s\n",
		ApplyColor("│", Border),
		"Bordered content",
		ApplyColor("│", Border))
	fmt.Printf("  %s\n", ApplyColor("└─────────────────────┘", Border))
}

// DemoAll shows all color capabilities
func DemoAll() {
	ShowPalette()
	fmt.Println()
	ShowThemes()
	ShowEffects()
}
