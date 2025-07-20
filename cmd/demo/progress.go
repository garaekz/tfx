package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/progress"
)

func runProgressDemo() {
	fmt.Println("🎨 TFX Progress + Color System Integration Demo")
	fmt.Println("===============================================")

	// Demo 1: EXPRESS API - Material Design Theme
	fmt.Println("\n📋 EXPRESS API - Material Design Theme:")
	materialBar := progress.NewMaterialProgress(50, "Material Design")
	materialBar.Start()
	for i := 0; i <= 50; i += 5 {
		materialBar.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	materialBar.Complete("Material theme complete!")

	time.Sleep(300 * time.Millisecond)

	// Demo 2: FLUENT API - Dracula Theme with Rainbow Effect
	fmt.Println("\n🧛 FLUENT API - Dracula Theme with Rainbow Effect:")
	draculaBar := progress.Start(
		progress.WithTotal(40),
		progress.WithLabel("Dracula Magic"),
		progress.WithDraculaTheme(),
		progress.WithRainbowEffect(),
		progress.WithProgressWidth(35),
	)
	for i := 0; i <= 40; i += 4 {
		draculaBar.Set(i)
		time.Sleep(80 * time.Millisecond)
	}
	draculaBar.Complete("Dracula power unleashed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 3: CONFIG API - Nord Theme
	fmt.Println("\n❄️ CONFIG API - Nord Theme:")
	cfg := progress.DefaultProgressConfig()
	cfg.Total = 30
	cfg.Label = "Nordic Cool"
	cfg.Theme = progress.NordTheme
	cfg.Width = 25
	nordBar := progress.StartWith(cfg)
	for i := 0; i <= 30; i += 3 {
		nordBar.Set(i)
		time.Sleep(120 * time.Millisecond)
	}
	nordBar.Complete("Nord vibes achieved!")

	time.Sleep(300 * time.Millisecond)

	// Demo 4: OBJECT LIFECYCLE - Manual control
	fmt.Println("\n🔧 OBJECT LIFECYCLE - Manual control:")
	objectBar := progress.NewWithConfig(progress.ProgressConfig{
		Total: 40,
		Label: "Object Demo",
		Theme: progress.MaterialTheme,
	})

	fmt.Println("   Creating bar...")
	time.Sleep(300 * time.Millisecond)

	fmt.Println("   Starting progress...")
	objectBar.Start()
	for i := 0; i <= 40; i += 4 {
		objectBar.Set(i)
		time.Sleep(80 * time.Millisecond)
	}
	objectBar.Complete("Object demo completed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 5: THEMED PROGRESS BARS
	fmt.Println("\n🎨 THEMED PROGRESS BARS:")

	fmt.Println("   Material Design Theme:")
	materialQuick := progress.NewMaterialProgress(20, "Material")
	materialQuick.Start()
	for i := 0; i <= 20; i += 4 {
		materialQuick.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	materialQuick.Complete("Material completed!")

	time.Sleep(200 * time.Millisecond)

	fmt.Println("   Dracula Theme with Rainbow:")
	draculaQuick := progress.NewDraculaProgress(20, "Dracula")
	draculaQuick.Start()
	for i := 0; i <= 20; i += 4 {
		draculaQuick.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	draculaQuick.Complete("Dracula completed!")

	time.Sleep(200 * time.Millisecond)

	fmt.Println("   Nord Theme:")
	nordQuick := progress.NewNordProgress(20, "Nord")
	nordQuick.Start()
	for i := 0; i <= 20; i += 4 {
		nordQuick.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	nordQuick.Complete("Nord completed!")

	time.Sleep(200 * time.Millisecond)

	fmt.Println("   Rainbow Effect:")
	rainbowQuick := progress.NewRainbowProgress(20, "Rainbow")
	rainbowQuick.Start()
	for i := 0; i <= 20; i += 4 {
		rainbowQuick.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	rainbowQuick.Complete("Rainbow completed!")

	time.Sleep(300 * time.Millisecond)

	// Demo 6: Dynamic Theme Switching
	fmt.Println("\n🔄 Dynamic Theme Switching:")
	dynamicBar := progress.Start(
		progress.WithTotal(45),
		progress.WithLabel("Theme Switcher"),
		progress.WithProgressWidth(35),
	)

	themes := []progress.ProgressTheme{
		progress.MaterialTheme,
		progress.DraculaTheme,
		progress.NordTheme,
	}

	for i := 0; i <= 45; i += 9 {
		if len(themes) > 0 {
			themeIndex := (i / 9) % len(themes)
			dynamicBar.SetTheme(themes[themeIndex])
		}
		dynamicBar.Set(i)
		time.Sleep(200 * time.Millisecond)
	}
	dynamicBar.Complete("All themes showcased!")

	time.Sleep(300 * time.Millisecond)

	// Demo 7: Error Handling
	fmt.Println("\n❌ Error Handling Demo:")
	errorBar := progress.Start(
		progress.WithTotal(20),
		progress.WithLabel("Risky Operation"),
		progress.WithMaterialTheme(),
		progress.WithProgressWidth(25),
	)
	for i := 0; i <= 15; i += 3 {
		errorBar.Set(i)
		time.Sleep(150 * time.Millisecond)
	}
	errorBar.Fail("Operation failed - this is how errors look!")

	time.Sleep(300 * time.Millisecond)

	// Demo 8: Global API showcase
	fmt.Println("\n🌍 Global API Demo:")
	globalBar := progress.Start(
		progress.WithTotal(30),
		progress.WithLabel("Global Operation"),
	)
	for i := 0; i <= 30; i += 3 {
		globalBar.Set(i)
		time.Sleep(80 * time.Millisecond)
	}
	globalBar.Complete("Global operation completed!")

	fmt.Println("\n🎉 Progress Demo Complete!")
	fmt.Println("TFX Progress features:")
	fmt.Println("• 3 API styles: EXPRESS, CONFIG, FLUENT")
	fmt.Println("• Multiple professional themes")
	fmt.Println("• Visual effects and animations")
	fmt.Println("• Smart terminal detection")
	fmt.Println("• Dynamic theme switching")
	fmt.Println("• Error handling and completion messages")
}
