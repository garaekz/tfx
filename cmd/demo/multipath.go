package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/progress"
)

func runMultipathDemo() {
	fmt.Println("🛤️  TFX Multipath API Demonstrations")
	fmt.Println("====================================")

	fmt.Println("\nTFX follows a multipath API philosophy:")
	fmt.Println("1. EXPRESS - Quick defaults: Start()")
	fmt.Println("2. CONFIG - Structured: StartWith(cfg)")
	fmt.Println("3. FLUENT - Functional options: Start(WithX, WithY)")
	fmt.Println("4. OBJECT - Manual lifecycle: New(), NewWithConfig()")

	// 1. EXPRESS API - Zero-config usage
	fmt.Println("\n═══ 1. EXPRESS API ═══")
	fmt.Println("progress.Start() // Zero-config, auto-start")
	p1 := progress.Start()
	for i := 0; i <= 100; i += 10 {
		p1.Set(i)
		time.Sleep(50 * time.Millisecond)
	}
	p1.Complete("Express completed!")

	time.Sleep(800 * time.Millisecond)

	// 2. CONFIG API - Structured configuration
	fmt.Println("\n═══ 2. CONFIG API ═══")
	fmt.Println("cfg := ProgressConfig{...}")
	fmt.Println("progress.StartWith(cfg) // Typed config")
	cfg := progress.ProgressConfig{
		Total: 50,
		Label: "Config Process",
		Width: 30,
		Theme: progress.DraculaTheme,
	}
	p2 := progress.StartWith(cfg)
	for i := 0; i <= 50; i += 5 {
		p2.Set(i)
		time.Sleep(80 * time.Millisecond)
	}
	p2.Complete("Config completed!")

	time.Sleep(800 * time.Millisecond)

	// 3. FLUENT API - Functional options
	fmt.Println("\n═══ 3. FLUENT API ═══")
	fmt.Println("progress.Start(WithTotal(75), WithLabel(...), WithTheme(...))")
	p3 := progress.Start(
		progress.WithTotal(75),
		progress.WithLabel("Fluent Process"),
		progress.WithProgressWidth(25),
		progress.WithNordTheme(),
		progress.WithRainbowEffect(),
	)
	for i := 0; i <= 75; i += 3 {
		p3.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	p3.Complete("Fluent completed!")

	time.Sleep(800 * time.Millisecond)

	// 4. OBJECT LIFECYCLE API - Manual control
	fmt.Println("\n═══ 4. OBJECT LIFECYCLE API ═══")
	fmt.Println("bar := progress.New() // Create without starting")
	fmt.Println("bar.Start() // Manual control")

	bar := progress.New()
	fmt.Println("Created bar (not started yet)...")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Starting bar manually...")
	bar.Start()
	for i := 0; i <= 100; i += 5 {
		bar.Set(i)
		time.Sleep(40 * time.Millisecond)
	}
	bar.Complete("Object lifecycle completed!")

	time.Sleep(800 * time.Millisecond)

	// 5. OBJECT WITH CONFIG
	fmt.Println("\n═══ 5. OBJECT WITH CONFIG ═══")
	fmt.Println("bar := progress.NewWithConfig(cfg)")
	configuredBar := progress.NewWithConfig(progress.ProgressConfig{
		Total: 60,
		Label: "Configured Object",
		Theme: progress.MaterialTheme,
		Width: 35,
	})

	fmt.Println("Starting configured object...")
	configuredBar.Start()
	for i := 0; i <= 60; i += 4 {
		configuredBar.Set(i)
		time.Sleep(50 * time.Millisecond)
	}
	configuredBar.Complete("Configured object completed!")

	time.Sleep(800 * time.Millisecond)

	// 6. Same internal implementation
	fmt.Println("\n═══ 6. UNIFIED IMPLEMENTATION ═══")
	fmt.Println("All paths use the same internal 'newProgress()' function")
	fmt.Println("with 'share.OverloadWithOptions' for type safety")

	fmt.Println("\nDemonstrating equivalent calls:")

	// Three equivalent ways to create the same progress bar
	fmt.Println("1. Start(cfg)")
	quick1 := progress.Start(progress.ProgressConfig{Total: 20, Label: "Method 1"})
	for i := 0; i <= 20; i += 5 {
		quick1.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	quick1.Complete("Method 1 done!")

	time.Sleep(300 * time.Millisecond)

	fmt.Println("2. StartWith(cfg)")
	quick2 := progress.StartWith(progress.ProgressConfig{Total: 20, Label: "Method 2"})
	for i := 0; i <= 20; i += 5 {
		quick2.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	quick2.Complete("Method 2 done!")

	time.Sleep(300 * time.Millisecond)

	fmt.Println("3. Start(WithTotal, WithLabel)")
	quick3 := progress.Start(
		progress.WithTotal(20),
		progress.WithLabel("Method 3"),
	)
	for i := 0; i <= 20; i += 5 {
		quick3.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
	quick3.Complete("Method 3 done!")

	fmt.Println("\n✅ All methods produce the same result!")
	fmt.Println("✅ Choose the style that fits your use case!")
	fmt.Println("✅ Multipath API demonstration completed!")
}
