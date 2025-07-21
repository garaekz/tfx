package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/progress"
)

func runShowcaseDemo() {
	fmt.Println("🎨 TFX Complete System Showcase")
	fmt.Println("================================")

	// Demo 1: Custom Theme Creation
	fmt.Println("\n🌅 Custom Theme Creation:")
	customTheme := progress.CreateCustomTheme(
		"sunset",
		color.MaterialOrange,     // Complete color
		color.MaterialDeepPurple, // Incomplete color
		color.MaterialPurple,     // Label color
	)
	customBar := progress.Start(
		progress.WithTotal(25),
		progress.WithLabel("Sunset Gradient"),
		progress.WithProgressTheme(customTheme),
		progress.WithGradientEffect(),
		progress.WithProgressWidth(30),
	)
	for i := 0; i <= 25; i += 2 {
		customBar.Set(i)
		time.Sleep(150 * time.Millisecond)
	}
	customBar.Complete("Sunset gradient complete!")

	time.Sleep(300 * time.Millisecond)

	// Demo 2: GitHub Theme
	fmt.Println("\n🐙 GitHub Theme:")
	githubBar := progress.Start(
		progress.WithTotal(60),
		progress.WithLabel("Git Operations"),
		progress.WithProgressTheme(progress.GitHubTheme),
		progress.WithProgressWidth(40),
	)
	for i := 0; i <= 60; i += 6 {
		githubBar.Set(i)
		time.Sleep(60 * time.Millisecond)
	}
	githubBar.Complete("Repository cloned successfully!")

	time.Sleep(300 * time.Millisecond)

	// Demo 3: VS Code Theme
	fmt.Println("\n💻 VS Code Theme:")
	vscodeBar := progress.Start(
		progress.WithTotal(35),
		progress.WithLabel("Code Compilation"),
		progress.WithProgressTheme(progress.VSCodeTheme),
		progress.WithProgressWidth(30),
	)
	for i := 0; i <= 35; i += 5 {
		vscodeBar.Set(i)
		time.Sleep(90 * time.Millisecond)
	}
	vscodeBar.Complete("Build successful!")

	time.Sleep(300 * time.Millisecond)

	// Demo 4: Neon Theme with Effects
	fmt.Println("\n✨ Neon Theme:")
	neonBar := progress.Start(
		progress.WithTotal(20),
		progress.WithLabel("Neon Glow"),
		progress.WithProgressTheme(progress.NeonTheme),
		progress.WithProgressEffect(progress.EffectGlow),
		progress.WithProgressWidth(25),
	)
	for i := 0; i <= 20; i += 2 {
		neonBar.Set(i)
		time.Sleep(200 * time.Millisecond)
	}
	neonBar.Complete("Neon effect complete!")

	time.Sleep(300 * time.Millisecond)

	// Demo 5: All Available Themes Showcase
	fmt.Println("\n🎨 All Available Themes:")
	for _, theme := range progress.AllProgressThemes {
		fmt.Printf("Testing %s theme: ", theme.Name)
		quickBar := progress.Start(
			progress.WithTotal(10),
			progress.WithLabel(theme.Name),
			progress.WithProgressTheme(theme),
			progress.WithProgressWidth(15),
		)
		for i := 0; i <= 10; i += 2 {
			quickBar.Set(i)
			time.Sleep(50 * time.Millisecond)
		}
		quickBar.Complete("✓")
	}

	time.Sleep(300 * time.Millisecond)

	// Demo 6: Advanced Spinner Examples
	fmt.Println("\n🌀 Advanced Spinner Examples:")

	// Custom frames spinner
	customFramesSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Custom animation..."),
		progress.WithSpinnerFrames([]string{"🌱", "🌿", "🍃", "🌳", "🌲"}),
		progress.WithSpinnerMaterialTheme(),
	)
	customFramesSpinner.Start()
	time.Sleep(2 * time.Second)
	customFramesSpinner.Stop("Growth complete!")

	// Emoji loading spinner
	emojiSpinner := progress.NewSpinnerWith(
		progress.WithMessage("Loading with style..."),
		progress.WithSpinnerFrames([]string{"🚀", "✨", "🌟", "⭐", "🌠", "💫"}),
		progress.WithSpinnerDraculaTheme(),
	)
	emojiSpinner.Start()
	time.Sleep(2 * time.Second)
	emojiSpinner.Stop("Mission accomplished!")

	time.Sleep(300 * time.Millisecond)

	// Demo 7: Complex Workflow Simulation
	fmt.Println("\n🔄 Complex Workflow Simulation:")

	workflowSteps := []struct {
		name     string
		duration time.Duration
		total    int
	}{
		{"Initializing", 800 * time.Millisecond, 15},
		{"Loading config", 1200 * time.Millisecond, 25},
		{"Processing data", 2000 * time.Millisecond, 40},
		{"Finalizing", 600 * time.Millisecond, 10},
	}

	for i, step := range workflowSteps {
		stepBar := progress.Start(
			progress.WithTotal(step.total),
			progress.WithLabel(step.name),
			progress.WithProgressWidth(30),
		)

		// Simulate incremental progress
		increment := step.total / 5
		for j := 0; j <= step.total; j += increment {
			if j > step.total {
				j = step.total
			}
			stepBar.Set(j)
			time.Sleep(step.duration / 5)
		}

		if i == len(workflowSteps)-1 {
			stepBar.Complete("Workflow completed successfully!")
		} else {
			stepBar.Complete("✓")
		}

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n🎉 Showcase Complete!")
	fmt.Println("TFX System Integration Features:")
	fmt.Println("• Custom theme creation with color system")
	fmt.Println("• Professional pre-built themes")
	fmt.Println("• Advanced visual effects")
	fmt.Println("• Complex workflow simulations")
	fmt.Println("• Emoji and custom frame animations")
	fmt.Println("• Seamless progress and spinner integration")
}
