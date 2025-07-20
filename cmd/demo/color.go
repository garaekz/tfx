package main

import (
	"fmt"
	"strings"

	"github.com/garaekz/tfx/color"
)

func runColorDemo() {
	fmt.Println("🎨 TFX Color System Demonstrations")
	fmt.Println("==================================")

	// 1. Basic color creation methods
	fmt.Println("\n1. COLOR CREATION METHODS:")
	
	fmt.Print("   RGB Colors: ")
	red := color.RGB(255, 0, 0)
	green := color.RGB(0, 255, 0)
	blue := color.RGB(0, 0, 255)
	fmt.Printf("%s %s %s\n", 
		red.Apply("RED"), 
		green.Apply("GREEN"), 
		blue.Apply("BLUE"))

	fmt.Print("   Hex Colors: ")
	orange := color.Hex("#FF8C00")
	purple := color.Hex("#9932CC")
	cyan := color.Hex("#00FFFF")
	fmt.Printf("%s %s %s\n", 
		orange.Apply("ORANGE"), 
		purple.Apply("PURPLE"), 
		cyan.Apply("CYAN"))

	fmt.Print("   ANSI Colors: ")
	fmt.Printf("%s %s %s %s\n",
		color.ColorRed.Apply("RED"),
		color.ColorGreen.Apply("GREEN"),
		color.ColorBlue.Apply("BLUE"),
		color.ColorYellow.Apply("YELLOW"))

	// 2. Color configuration
	fmt.Println("\n2. COLOR CONFIGURATION:")
	customColor := color.NewColor(color.ColorConfig{
		R: 255, G: 100, B: 200,
		Name: "Custom Pink",
	})
	fmt.Printf("   Custom Color: %s\n", customColor.Apply("Custom Pink"))

	// 3. Functional options
	fmt.Println("\n3. FUNCTIONAL OPTIONS:")
	fluentColor := color.NewColorWith(
		color.WithRGB(100, 200, 255),
		color.WithName("Fluent Blue"),
	)
	fmt.Printf("   Fluent Color: %s\n", fluentColor.Apply("Fluent Blue"))

	// 4. Background colors
	fmt.Println("\n4. BACKGROUND COLORS:")
	bgRed := color.ColorRed.Bg()
	bgGreen := color.ColorGreen.Bg()
	bgBlue := color.ColorBlue.Bg()
	fmt.Printf("   %s %s %s\n",
		color.ApplyBg("RED BG", color.ColorWhite, bgRed, color.ModeTrueColor),
		color.ApplyBg("GREEN BG", color.ColorBlack, bgGreen, color.ModeTrueColor),
		color.ApplyBg("BLUE BG", color.ColorWhite, bgBlue, color.ModeTrueColor))

	// 5. Styling system
	fmt.Println("\n5. STYLING SYSTEM:")
	
	fmt.Print("   Text Styling: ")
	fmt.Printf("%s ", color.NewStyleWith(
		color.WithText("BOLD"),
		color.WithForeground(color.ColorRed),
		color.WithBold(),
	))
	fmt.Printf("%s ", color.NewStyleWith(
		color.WithText("ITALIC"),
		color.WithForeground(color.ColorGreen),
		color.WithItalic(),
	))
	fmt.Printf("%s\n", color.NewStyleWith(
		color.WithText("UNDERLINE"),
		color.WithForeground(color.ColorBlue),
		color.WithUnderline(),
	))

	// 6. Badges
	fmt.Println("\n6. BADGES:")
	fmt.Printf("   %s %s %s %s %s\n",
		color.SuccessBadge("SUCCESS"),
		color.ErrorBadge("ERROR"),
		color.WarningBadge("WARNING"),
		color.InfoBadge("INFO"),
		color.DebugBadge("DEBUG"))

	// 7. Material Design colors
	fmt.Println("\n7. MATERIAL DESIGN COLORS:")
	materialColors := []color.Color{
		color.MaterialRed, color.MaterialPink, color.MaterialPurple,
		color.MaterialBlue, color.MaterialCyan, color.MaterialTeal,
		color.MaterialGreen, color.MaterialLime, color.MaterialYellow,
		color.MaterialOrange,
	}
	fmt.Print("   ")
	for _, c := range materialColors {
		fmt.Printf("%s ", c.Apply("██"))
	}
	fmt.Println()

	// 8. Tailwind colors
	fmt.Println("\n8. TAILWIND COLORS:")
	tailwindColors := []color.Color{
		color.TailwindRed, color.TailwindOrange, color.TailwindAmber,
		color.TailwindYellow, color.TailwindLime, color.TailwindGreen,
		color.TailwindTeal, color.TailwindCyan, color.TailwindBlue,
		color.TailwindPurple, color.TailwindPink,
	}
	fmt.Print("   ")
	for _, c := range tailwindColors {
		fmt.Printf("%s ", c.Apply("██"))
	}
	fmt.Println()

	// 9. Dracula theme colors
	fmt.Println("\n9. DRACULA THEME COLORS:")
	draculaColors := []color.Color{
		color.DraculaPurple, color.DraculaPink, color.DraculaGreen,
		color.DraculaOrange, color.DraculaRed, color.DraculaYellow,
		color.DraculaCyan,
	}
	fmt.Print("   ")
	for _, c := range draculaColors {
		fmt.Printf("%s ", c.Apply("██"))
	}
	fmt.Println()

	// 10. Rainbow text
	fmt.Println("\n10. SPECIAL EFFECTS:")
	fmt.Printf("   Rainbow: %s\n", color.RainbowText("RAINBOW TEXT", color.ModeTrueColor))
	
	// 11. Gradient text
	gradientColors := []color.Color{color.MaterialRed, color.MaterialOrange, color.MaterialYellow}
	fmt.Printf("   Gradient: %s\n", color.GradientText("GRADIENT TEXT", gradientColors, color.ModeTrueColor))

	// 12. Progress bar
	fmt.Println("\n11. PROGRESS BAR:")
	fmt.Printf("   Progress: %s\n", color.ProgressBar(7, 10, 20, color.MaterialGreen, color.ColorBrightBlack))

	// 13. Bordered text
	fmt.Println("\n12. BORDERED TEXT:")
	borderText := "TFX Demo\nColor System\nShowcase"
	fmt.Println(color.Border(borderText, color.MaterialBlue))

	// 14. Palettes
	fmt.Println("\n13. COLOR PALETTES:")
	
	fmt.Println("   Status Palette:")
	statusPalette := color.StatusPalette()
	for _, name := range []string{"success", "error", "warning", "info", "debug"} {
		if c, exists := statusPalette.Get(name); exists {
			fmt.Printf("     %s: %s\n", strings.ToUpper(name), c.Apply(name))
		}
	}

	fmt.Println("\n✅ Color demonstration completed!")
}