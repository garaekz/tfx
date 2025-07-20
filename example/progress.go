// example/progress.go - SHOWCASE OF NEW MULTIPATH INTEGRATION
package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/progress"
)

func main() {
	fmt.Println("🎨 TFX Progress + Color System Integration Demo")
	fmt.Println("===============================================")

	// Demo 1: EXPRESS API - Quick default
	fmt.Println("\n📋 EXPRESS API - Material Design Theme:")
	materialBar := progress.NewMaterialProgress(50, "Material Design")
	materialBar.Start()
	for i := 0; i <= 50; i += 5 {
		materialBar.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	materialBar.Complete("Material theme complete!")

	// Demo 2: FLUENT API - Dracula Theme with Rainbow Effect
	fmt.Println("\n🧛 FLUENT API - Dracula Theme with Rainbow Effect:")
	draculaBar := progress.NewProgressWith(
		progress.WithTotal(40),
		progress.WithLabel("Dracula Magic"),
		progress.WithDraculaTheme(),
		progress.WithRainbowEffect(),
		progress.WithProgressWidth(35),
	)
	draculaBar.Start()
	for i := 0; i <= 40; i += 4 {
		draculaBar.Set(i)
		time.Sleep(80 * time.Millisecond)
	}
	draculaBar.Complete("Dracula power unleashed!")

	// Demo 3: CONFIG API - Nord Theme
	fmt.Println("\n❄️ CONFIG API - Nord Theme:")
	cfg := progress.DefaultProgressConfig()
	cfg.Total = 30
	cfg.Label = "Nordic Cool"
	cfg.Theme = progress.NordTheme
	cfg.Width = 25
	nordBar := progress.NewProgress(cfg)
	nordBar.Start()
	for i := 0; i <= 30; i += 3 {
		nordBar.Set(i)
		time.Sleep(120 * time.Millisecond)
	}
	nordBar.Complete("Nord vibes achieved!")

	// Demo 4: Custom Theme with Gradient Effect
	fmt.Println("\n🌈 Custom Theme with Gradient Effect:")
	customTheme := progress.CreateCustomTheme(
		"sunset",
		color.MaterialOrange, // Complete color
		color.ANSI(8),        // Incomplete color (dark gray)
		color.MaterialPurple, // Label color
	)
	customBar := progress.NewProgressWith(
		progress.WithTotal(25),
		progress.WithLabel("Sunset Gradient"),
		progress.WithProgressTheme(customTheme),
		progress.WithGradientEffect(),
		progress.WithProgressWidth(30),
	)
	customBar.Start()
	for i := 0; i <= 25; i += 2 {
		customBar.Set(i)
		time.Sleep(150 * time.Millisecond)
	}
	customBar.Complete("Sunset gradient complete!")

	// Demo 5: GitHub Theme
	fmt.Println("\n🐙 GitHub Theme:")
	githubBar := progress.NewProgressWith(
		progress.WithTotal(60),
		progress.WithLabel("Git Operations"),
		progress.WithProgressTheme(progress.GitHubTheme),
		progress.WithProgressWidth(40),
	)
	githubBar.Start()
	for i := 0; i <= 60; i += 6 {
		githubBar.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	githubBar.Complete("Repository cloned successfully!")

	// Demo 6: VS Code Theme
	fmt.Println("\n💻 VS Code Theme:")
	vscodeBar := progress.NewProgressWith(
		progress.WithTotal(35),
		progress.WithLabel("Code Compilation"),
		progress.WithProgressTheme(progress.VSCodeTheme),
		progress.WithProgressWidth(30),
	)
	vscodeBar.Start()
	for i := 0; i <= 35; i += 5 {
		vscodeBar.Set(i)
		time.Sleep(90 * time.Millisecond)
	}
	vscodeBar.Complete("Build successful!")

	// Demo 7: Neon Theme with Effects
	fmt.Println("\n✨ Neon Theme:")
	neonBar := progress.NewProgressWith(
		progress.WithTotal(20),
		progress.WithLabel("Neon Glow"),
		progress.WithProgressTheme(progress.NeonTheme),
		progress.WithProgressEffect(progress.EffectGlow),
		progress.WithProgressWidth(25),
	)
	neonBar.Start()
	for i := 0; i <= 20; i += 2 {
		neonBar.Set(i)
		time.Sleep(200 * time.Millisecond)
	}
	neonBar.Complete("Neon effect complete!")

	// Demo 8: Spinner Showcase - EXPRESS API
	fmt.Println("\n🌀 Enhanced Spinners:")

	// Material Spinner - EXPRESS API
	materialSpinner := progress.NewMaterialSpinner("Loading Material...")
	materialSpinner.Start()
	time.Sleep(2 * time.Second)
	materialSpinner.Stop("Material loaded!")

	// Dracula Rainbow Spinner - FLUENT API
	draculaSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Casting spells..."),
		progress.WithSpinnerDraculaTheme(),
		progress.WithSpinnerRainbow(),
		progress.WithSpinnerFrames([]string{"🌙", "🌘", "🌗", "🌖", "🌕", "🌔", "🌓", "🌒"}),
	)
	draculaSpinner.Start()
	time.Sleep(2 * time.Second)
	draculaSpinner.Stop("Spells cast successfully!")

	// Nord Spinner - CONFIG API
	spinnerCfg := progress.DefaultSpinnerConfig()
	spinnerCfg.Message = "Processing Nordic data..."
	spinnerCfg.Theme = progress.NordTheme
	spinnerCfg.Frames = []string{"❄", "❅", "❆", "✱", "✲", "✳"}
	nordSpinner := progress.NewSpinner(spinnerCfg)
	nordSpinner.Start()
	time.Sleep(time.Duration(1.5 * float64(time.Second)))
	nordSpinner.Stop("Nordic processing complete!")

	// Demo 9: Dynamic Theme Switching
	fmt.Println("\n🔄 Dynamic Theme Switching:")
	dynamicBar := progress.NewProgressWith(
		progress.WithTotal(45),
		progress.WithLabel("Theme Switcher"),
		progress.WithProgressWidth(35),
	)
	dynamicBar.Start()

	themes := []progress.ProgressTheme{
		progress.MaterialTheme,
		progress.DraculaTheme,
		progress.NordTheme,
		progress.GitHubTheme,
		progress.VSCodeTheme,
	}

	for i := 0; i <= 45; i += 9 {
		themeIndex := (i / 9) % len(themes)
		dynamicBar.SetTheme(themes[themeIndex])
		dynamicBar.Set(i)
		time.Sleep(300 * time.Millisecond)
	}
	dynamicBar.Complete("All themes showcased!")

	// Demo 10: Error Handling
	fmt.Println("\n❌ Error Handling Demo:")
	errorBar := progress.NewProgressWith(
		progress.WithTotal(20),
		progress.WithLabel("Risky Operation"),
		progress.WithProgressTheme(progress.MaterialTheme),
		progress.WithProgressWidth(25),
	)
	errorBar.Start()
	for i := 0; i <= 15; i += 3 {
		errorBar.Set(i)
		time.Sleep(200 * time.Millisecond)
	}
	errorBar.Fail("Operation failed - this is how errors look!")

	// Demo 11: All Available Themes Showcase
	fmt.Println("\n🎨 All Available Themes:")
	for _, theme := range progress.AllProgressThemes {
		fmt.Printf("Testing %s theme: ", theme.Name)
		quickBar := progress.NewProgressWith(
			progress.WithTotal(10),
			progress.WithLabel(theme.Name),
			progress.WithProgressTheme(theme),
			progress.WithProgressWidth(15),
		)
		quickBar.Start()
		for i := 0; i <= 10; i += 2 {
			quickBar.Set(i)
			time.Sleep(50 * time.Millisecond)
		}
		quickBar.Complete("✓")
	}

	// Demo 12: Global API showcase (EXPRESS)
	fmt.Println("\n🌍 Global API Demo:")
	progress.StartGlobalProgress(30, "Global Operation")
	for i := 0; i <= 30; i += 3 {
		progress.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	progress.Complete("Global operation completed!")

	fmt.Println("\n🎉 Demo Complete!")
	fmt.Println("TFX now has full Color + Progress integration with:")
	fmt.Println("• 3 API styles: EXPRESS, CONFIG, FLUENT")
	fmt.Println("• 8 professional themes (Material, Dracula, Nord, GitHub, VS Code, Tailwind, Rainbow, Neon)")
	fmt.Println("• 4 visual effects (None, Gradient, Rainbow, Glow)")
	fmt.Println("• Smart terminal detection and color optimization")
	fmt.Println("• Custom theme creation")
	fmt.Println("• Enhanced spinners with color effects")
	fmt.Println("• Dynamic theme switching")
	fmt.Println("• Professional completion and error messages")
}

// Usage examples showing the three API patterns
func usageExamples() {
	// 1. EXPRESS API - Quick constructors
	bar1 := progress.NewMaterialProgress(100, "Material Design")
	bar2 := progress.NewDraculaProgress(100, "Dracula Magic")
	bar3 := progress.NewNordProgress(100, "Nordic Cool")
	bar4 := progress.NewRainbowProgress(100, "Rainbow Effects")

	// 2. CONFIG API - Structured configuration
	cfg := progress.DefaultProgressConfig()
	cfg.Total = 100
	cfg.Label = "Custom Setup"
	cfg.Theme = progress.GitHubTheme
	cfg.Effect = progress.EffectGradient
	cfg.Width = 50
	bar5 := progress.NewProgress(cfg)

	// 3. FLUENT API - Functional options
	bar6 := progress.NewProgressWith(
		progress.WithTotal(100),
		progress.WithLabel("Fluent Style"),
		progress.WithProgressTheme(progress.VSCodeTheme),
		progress.WithGradientEffect(),
		progress.WithProgressWidth(60),
	)

	// Custom theme creation
	myTheme := progress.CreateCustomTheme(
		"corporate",
		color.Hex("#007ACC"), // Complete color
		color.Hex("#E1E4E8"), // Incomplete color
		color.Hex("#24292E"), // Label color
	)

	bar7 := progress.NewProgressWith(
		progress.WithTotal(100),
		progress.WithLabel("Corporate Style"),
		progress.WithProgressTheme(myTheme),
	)

	// Spinner examples - all three APIs
	// EXPRESS
	spinner1 := progress.NewMaterialSpinner("Loading...")

	// CONFIG
	spinnerCfg := progress.DefaultSpinnerConfig()
	spinnerCfg.Message = "Processing..."
	spinnerCfg.Theme = progress.DraculaTheme
	spinnerCfg.Effect = progress.SpinnerEffectRainbow
	spinner2 := progress.NewSpinner(spinnerCfg)

	// FLUENT
	spinner3 := progress.NewSpinnerWith(
		progress.WithMessage("Working..."),
		progress.WithSpinnerNordTheme(),
		progress.WithDotsFrames(),
	)

	// Use the progress bars and spinners...
	_ = bar1
	_ = bar2
	_ = bar3
	_ = bar4
	_ = bar5
	_ = bar6
	_ = bar7
	_ = spinner1
	_ = spinner2
	_ = spinner3
}
