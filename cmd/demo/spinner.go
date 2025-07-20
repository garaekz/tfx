package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/progress"
)

func runSpinnerDemo() {
	fmt.Println("🌀 TFX Enhanced Spinners Demo")
	fmt.Println("=============================")

	// Demo 1: Basic Spinner
	fmt.Println("\n⚡ Basic Spinner:")
	basicSpinner := progress.StartSpinner("Loading basic data...")
	basicSpinner.Start()
	time.Sleep(2 * time.Second)
	basicSpinner.Success("Basic data loaded!")

	time.Sleep(300 * time.Millisecond)

	// Demo 2: Custom Configuration
	fmt.Println("\n🔧 Custom Configuration:")
	customSpinner := progress.NewSpinner(progress.SpinnerConfig{
		Message:  "Processing files...",
		Interval: 120 * time.Millisecond,
		Theme:    progress.DraculaTheme,
	})
	customSpinner.Start()
	time.Sleep(2 * time.Second)
	customSpinner.Success("Files processed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 3: Functional Options
	fmt.Println("\n🎯 Functional Options:")
	functionalSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Syncing data..."),
		progress.WithSpinnerNordTheme(),
		progress.WithDotsFrames(),
	)
	functionalSpinner.Start()
	time.Sleep(2 * time.Second)
	functionalSpinner.Success("Sync completed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 4: Different Frame Styles
	fmt.Println("\n🎬 Different Frame Styles:")

	fmt.Println("   Classic Dots:")
	dotsSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Classic loading..."),
		progress.WithDotsFrames(),
	)
	dotsSpinner.Start()
	time.Sleep(1500 * time.Millisecond)
	dotsSpinner.Success("Classic complete!")

	fmt.Println("   Modern Arrows:")
	arrowsSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Arrow loading..."),
		progress.WithArrowFrames(),
	)
	arrowsSpinner.Start()
	time.Sleep(1500 * time.Millisecond)
	arrowsSpinner.Success("Arrows complete!")

	fmt.Println("   Moon Phases:")
	moonSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Casting lunar spells..."),
		progress.WithSpinnerDraculaTheme(),
		progress.WithSpinnerFrames([]string{"🌘", "🌗", "🌖", "🌕", "🌔", "🌓", "🌒"}),
	)
	moonSpinner.Start()
	time.Sleep(2 * time.Second)
	moonSpinner.Success("Lunar magic complete! ✨")

	fmt.Println("   Growing Plant:")
	plantSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Growing garden..."),
		progress.WithSpinnerMaterialTheme(),
		progress.WithSpinnerFrames([]string{"🌱", "🌿", "🍃", "🌳", "🌲"}),
	)
	plantSpinner.Start()
	time.Sleep(2 * time.Second)
	plantSpinner.Success("Garden flourished! 🌳")

	fmt.Println("   Space Mission:")
	spaceSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Launching to space..."),
		progress.WithSpinnerDraculaTheme(),
		progress.WithSpinnerFrames([]string{"🚀", "✨", "🌟", "⭐", "🌠", "💫"}),
	)
	spaceSpinner.Start()
	time.Sleep(2 * time.Second)
	spaceSpinner.Success("Mission accomplished! 🛸")

	fmt.Println("   Snowflake Melting:")
	snowSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Winter magic..."),
		progress.WithSpinnerNordTheme(),
		progress.WithSpinnerFrames([]string{"❄️", "❅", "❆", "✳", "✲", "✱", "·", "."}),
	)
	snowSpinner.Start()
	time.Sleep(2 * time.Second)
	snowSpinner.Success("Winter spell complete! ❄️")

	fmt.Println("   Clock Ticking:")
	clockSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Time passing..."),
		progress.WithSpinnerMaterialTheme(),
		progress.WithSpinnerFrames([]string{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗"}),
	)
	clockSpinner.Start()
	time.Sleep(2 * time.Second)
	clockSpinner.Success("Time travel complete! ⏰")

	time.Sleep(300 * time.Millisecond)

	// Demo 5: Theme Showcase
	fmt.Println("\n🎨 Theme Showcase:")

	themes := []struct {
		name    string
		theme   progress.ProgressTheme
		message string
	}{
		{"Material", progress.MaterialTheme, "Material processing..."},
		{"Dracula", progress.DraculaTheme, "Dark magic..."},
		{"Nord", progress.NordTheme, "Arctic winds..."},
	}

	for _, t := range themes {
		fmt.Printf("   %s Theme: ", t.name)
		spinner := progress.NewSpinnerWith(
			progress.WithMessage(t.message),
			progress.WithSpinnerTheme(t.theme),
		)
		spinner.Start()
		time.Sleep(1200 * time.Millisecond)
		spinner.Success("Complete!")
		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(300 * time.Millisecond)

	// Demo 6: Error Handling
	fmt.Println("\n❌ Error Handling:")
	errorSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Risky operation..."),
		progress.WithSpinnerMaterialTheme(),
	)
	errorSpinner.Start()
	time.Sleep(1500 * time.Millisecond)
	errorSpinner.Fail("Operation failed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 7: Dynamic Message Updates
	fmt.Println("\n🔄 Dynamic Message Updates:")
	dynamicSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Starting process..."),
		progress.WithSpinnerDraculaTheme(),
	)
	dynamicSpinner.Start()

	messages := []string{
		"Loading configuration...",
		"Connecting to server...",
		"Downloading data...",
		"Processing results...",
		"Finalizing...",
	}

	for _, msg := range messages {
		dynamicSpinner.SetMessage(msg)
		time.Sleep(400 * time.Millisecond)
	}

	dynamicSpinner.Success("All operations completed!")

	fmt.Println("\n🎉 Spinner Demo Complete!")
	fmt.Println("TFX Spinner features:")
	fmt.Println("• 3 API styles: EXPRESS, CONFIG, FLUENT")
	fmt.Println("• Multiple frame styles and themes")
	fmt.Println("• Dynamic message updates")
	fmt.Println("• Error handling")
	fmt.Println("• Custom frame animations")
	fmt.Println("• Smart terminal detection")
}
