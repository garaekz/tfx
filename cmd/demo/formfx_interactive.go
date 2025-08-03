package main

import (
	"fmt"
	"log"

	"github.com/garaekz/tfx/formfx"
)

// DemoInteractiveFormFX demonstrates the interactive capabilities of FormFX
func DemoInteractiveFormFX() {
	fmt.Println("🎯 FormFX Interactive Demo")
	fmt.Println("=========================")
	fmt.Println()

	// Demo 1: Interactive Select with arrow key navigation
	fmt.Println("Demo 1: Interactive Select")
	fmt.Println("- Use ↑↓ arrows or WASD to navigate")
	fmt.Println("- Press Enter to select")
	fmt.Println("- Press Esc or Q to cancel")
	fmt.Println()

	options := []string{
		"🚀 Start new project",
		"📁 Open existing project",
		"⚙️  Settings",
		"📊 View reports",
		"🔄 Sync data",
		"❌ Exit",
	}

	choice, err := formfx.Select("What would you like to do?", options)
	if err != nil {
		if err == formfx.ErrCanceled {
			fmt.Println("Operation canceled by user")
		} else {
			log.Printf("Select error: %v", err)
		}
		return
	}

	fmt.Printf("You selected: %s\n\n", options[choice])

	// Demo 2: Interactive Confirm with visual feedback
	fmt.Println("Demo 2: Interactive Confirm")
	fmt.Println("- Use ↑↓ arrows or WASD to choose Yes/No")
	fmt.Println("- Or press Y/N directly")
	fmt.Println("- Press Enter to confirm")
	fmt.Println("- Press Esc to cancel")
	fmt.Println()

	confirmed, err := formfx.Confirm("Do you want to proceed with this action?")
	if err != nil {
		if err == formfx.ErrCanceled {
			fmt.Println("Confirmation canceled by user")
		} else {
			log.Printf("Confirm error: %v", err)
		}
		return
	}

	if confirmed {
		fmt.Println("✅ Action confirmed!")
	} else {
		fmt.Println("❌ Action canceled.")
	}
	fmt.Println()

	// Demo 3: Combining multiple prompts
	fmt.Println("Demo 3: Multi-step Form")
	fmt.Println("=======================")
	fmt.Println()

	// Step 1: Project type selection
	projectTypes := []string{
		"Web Application",
		"Desktop Application",
		"Mobile App",
		"API Service",
		"CLI Tool",
	}

	projectType, err := formfx.Select("Select project type:", projectTypes)
	if err != nil {
		log.Printf("Project type selection error: %v", err)
		return
	}

	// Step 2: Confirm project type
	confirmMsg := fmt.Sprintf("Create a %s project?", projectTypes[projectType])
	proceed, err := formfx.Confirm(confirmMsg)
	if err != nil {
		log.Printf("Confirmation error: %v", err)
		return
	}

	if proceed {
		fmt.Printf("🎉 Creating %s project...\n", projectTypes[projectType])

		// Step 3: Final confirmation
		final, err := formfx.NewConfirm().
			Label("Are you absolutely sure?").
			Default(true).
			Show()

		if err != nil {
			log.Printf("Final confirmation error: %v", err)
			return
		}

		if final {
			fmt.Println("✅ Project creation confirmed!")
		} else {
			fmt.Println("❌ Project creation cancelled.")
		}
	} else {
		fmt.Println("❌ Project creation cancelled.")
	}

	fmt.Println()
	fmt.Println("🎯 FormFX Interactive Demo Complete!")
}
