package main

import (
	"fmt"

	"github.com/garaekz/tfx/color"
)

func runColorEncodingDemo() {
	fmt.Println("🎨 TFX Color Encoding System - Clean API Demonstration")
	fmt.Println("=====================================================")

	// 1. Building blocks (hardcore API - still available)
	fmt.Println("\n1. BUILDING BLOCKS (Hardcore API):")
	fmt.Printf("Hex Color: %s%s%s\n", color.NewHex("#FF5733").Render(color.ModeANSI), "Orange Text", color.Reset)
	fmt.Printf("RGB Color: %s%s%s\n", color.NewRGB(255, 87, 51).Render(color.ModeANSI), "Same Orange", color.Reset)
	fmt.Printf("ANSI Code: %s%s%s\n", color.NewANSI(1).Render(color.ModeANSI), "Red Text", color.Reset)
	fmt.Printf("256 Color: %s%s%s\n", color.NewColor256(196).Render(color.Mode256Color), "Bright Red", color.Reset)

	fmt.Println("\n✨ NEW CLEAN API - REVOLUTIONARY COLOR ENCODING SYSTEM ✨")
	fmt.Println("========================================================")

	// 2. ENCODING SYSTEMS - The elegant new API!
	fmt.Println("\n2. ENCODING SYSTEMS (color.ENCODING.Color):")
	
	// ANSI Colors (16 colors) - CLEAN!
	fmt.Println("\n   ANSI Colors (16-color support) - CLEAN API:")
	fmt.Printf("   color.ANSI.Red:     %s%s%s\n", color.ANSI.Red.Render(color.ModeANSI), "ANSI Red Text", color.Reset)
	fmt.Printf("   color.ANSI.Blue:    %s%s%s\n", color.ANSI.Blue.Render(color.ModeANSI), "ANSI Blue Text", color.Reset)
	fmt.Printf("   color.ANSI.Green:   %s%s%s\n", color.ANSI.Green.Render(color.ModeANSI), "ANSI Green Text", color.Reset)
	fmt.Printf("   color.ANSI.Yellow:  %s%s%s\n", color.ANSI.Yellow.Render(color.ModeANSI), "ANSI Yellow Text", color.Reset)

	// TrueColor Colors (24-bit RGB) - CLEAN!
	fmt.Println("\n   TrueColor Colors (24-bit RGB support) - CLEAN API:")
	fmt.Printf("   color.TrueColor.Red:     %s%s%s\n", color.TrueColor.Red.Render(color.ModeTrueColor), "TrueColor Red Text", color.Reset)
	fmt.Printf("   color.TrueColor.Blue:    %s%s%s\n", color.TrueColor.Blue.Render(color.ModeTrueColor), "TrueColor Blue Text", color.Reset)
	fmt.Printf("   color.TrueColor.Green:   %s%s%s\n", color.TrueColor.Green.Render(color.ModeTrueColor), "TrueColor Green Text", color.Reset)
	fmt.Printf("   color.TrueColor.Yellow:  %s%s%s\n", color.TrueColor.Yellow.Render(color.ModeTrueColor), "TrueColor Yellow Text", color.Reset)

	// 256 Colors - CLEAN!
	fmt.Println("\n   256-Color Colors (256-color support) - CLEAN API:")
	fmt.Printf("   color.Color256.Red:      %s%s%s\n", color.Color256.Red.Render(color.Mode256Color), "256-Color Red Text", color.Reset)
	fmt.Printf("   color.Color256.Blue:     %s%s%s\n", color.Color256.Blue.Render(color.Mode256Color), "256-Color Blue Text", color.Reset)
	fmt.Printf("   color.Color256.Green:    %s%s%s\n", color.Color256.Green.Render(color.Mode256Color), "256-Color Green Text", color.Reset)
	fmt.Printf("   color.Color256.Yellow:   %s%s%s\n", color.Color256.Yellow.Render(color.Mode256Color), "256-Color Yellow Text", color.Reset)

	// 3. THEME SYSTEMS - The beautiful new API!
	fmt.Println("\n3. THEME SYSTEMS (color.THEME.Color) - CLEAN API:")
	
	// Material Design Theme - CLEAN!
	fmt.Println("\n   Material Design Theme - CLEAN API:")
	fmt.Printf("   color.Material.Red:    %s%s%s\n", color.Material.Red.Render(color.ModeTrueColor), "Material Red", color.Reset)
	fmt.Printf("   color.Material.Blue:   %s%s%s\n", color.Material.Blue.Render(color.ModeTrueColor), "Material Blue", color.Reset)
	fmt.Printf("   color.Material.Green:  %s%s%s\n", color.Material.Green.Render(color.ModeTrueColor), "Material Green", color.Reset)
	fmt.Printf("   color.Material.Purple: %s%s%s\n", color.Material.Purple.Render(color.ModeTrueColor), "Material Purple", color.Reset)
	fmt.Printf("   color.Material.Orange: %s%s%s\n", color.Material.Orange.Render(color.ModeTrueColor), "Material Orange", color.Reset)

	// Dracula Theme - CLEAN!
	fmt.Println("\n   Dracula Theme - CLEAN API:")
	fmt.Printf("   color.Dracula.Red:    %s%s%s\n", color.Dracula.Red.Render(color.ModeTrueColor), "Dracula Red", color.Reset)
	fmt.Printf("   color.Dracula.Blue:   %s%s%s\n", color.Dracula.Blue.Render(color.ModeTrueColor), "Dracula Blue", color.Reset)
	fmt.Printf("   color.Dracula.Green:  %s%s%s\n", color.Dracula.Green.Render(color.ModeTrueColor), "Dracula Green", color.Reset)
	fmt.Printf("   color.Dracula.Purple: %s%s%s\n", color.Dracula.Purple.Render(color.ModeTrueColor), "Dracula Purple", color.Reset)
	fmt.Printf("   color.Dracula.Pink:   %s%s%s\n", color.Dracula.Pink.Render(color.ModeTrueColor), "Dracula Pink", color.Reset)

	// Nord Theme - CLEAN!
	fmt.Println("\n   Nord Theme - CLEAN API:")
	fmt.Printf("   color.Nord.Red:    %s%s%s\n", color.Nord.Red.Render(color.ModeTrueColor), "Nord Red", color.Reset)
	fmt.Printf("   color.Nord.Blue:   %s%s%s\n", color.Nord.Blue.Render(color.ModeTrueColor), "Nord Blue", color.Reset)
	fmt.Printf("   color.Nord.Green:  %s%s%s\n", color.Nord.Green.Render(color.ModeTrueColor), "Nord Green", color.Reset)
	fmt.Printf("   color.Nord.Purple: %s%s%s\n", color.Nord.Purple.Render(color.ModeTrueColor), "Nord Purple", color.Reset)
	fmt.Printf("   color.Nord.Orange: %s%s%s\n", color.Nord.Orange.Render(color.ModeTrueColor), "Nord Orange", color.Reset)

	// GitHub Theme - CLEAN!
	fmt.Println("\n   GitHub Theme - CLEAN API:")
	fmt.Printf("   color.GitHub.Red:    %s%s%s\n", color.GitHub.Red.Render(color.ModeTrueColor), "GitHub Red", color.Reset)
	fmt.Printf("   color.GitHub.Blue:   %s%s%s\n", color.GitHub.Blue.Render(color.ModeTrueColor), "GitHub Blue", color.Reset)
	fmt.Printf("   color.GitHub.Green:  %s%s%s\n", color.GitHub.Green.Render(color.ModeTrueColor), "GitHub Green", color.Reset)
	fmt.Printf("   color.GitHub.Purple: %s%s%s\n", color.GitHub.Purple.Render(color.ModeTrueColor), "GitHub Purple", color.Reset)
	fmt.Printf("   color.GitHub.Orange: %s%s%s\n", color.GitHub.Orange.Render(color.ModeTrueColor), "GitHub Orange", color.Reset)

	// 4. DEFAULT COLORS - The most elegant API!
	fmt.Println("\n4. DEFAULT COLORS (Follow Active Theme) - MOST ELEGANT:")
	fmt.Printf("Current theme: %s\n", color.GetDefaultTheme())
	fmt.Printf("Current encoding: %s\n", color.GetDefaultEncoding())
	
	fmt.Printf("color.Red:    %s%s%s\n", color.Red.Render(color.GetDefaultEncoding()), "Default Red (follows theme)", color.Reset)
	fmt.Printf("color.Blue:   %s%s%s\n", color.Blue.Render(color.GetDefaultEncoding()), "Default Blue (follows theme)", color.Reset)
	fmt.Printf("color.Green:  %s%s%s\n", color.Green.Render(color.GetDefaultEncoding()), "Default Green (follows theme)", color.Reset)
	fmt.Printf("color.Purple: %s%s%s\n", color.Purple.Render(color.GetDefaultEncoding()), "Default Purple (follows theme)", color.Reset)

	// 5. THEME SWITCHING DEMONSTRATION
	fmt.Println("\n5. THEME SWITCHING - WATCH THE MAGIC:")
	
	// Switch to Dracula
	fmt.Println("\n   Switching to Dracula theme...")
	color.UseDracula()
	fmt.Printf("   New theme: %s\n", color.GetDefaultTheme())
	fmt.Printf("   color.Red:  %s%s%s\n", color.Red.Render(color.GetDefaultEncoding()), "Now Dracula Red", color.Reset)
	fmt.Printf("   color.Blue: %s%s%s\n", color.Blue.Render(color.GetDefaultEncoding()), "Now Dracula Blue", color.Reset)

	// Switch to Nord
	fmt.Println("\n   Switching to Nord theme...")
	color.UseNord()
	fmt.Printf("   New theme: %s\n", color.GetDefaultTheme())
	fmt.Printf("   color.Red:  %s%s%s\n", color.Red.Render(color.GetDefaultEncoding()), "Now Nord Red", color.Reset)
	fmt.Printf("   color.Blue: %s%s%s\n", color.Blue.Render(color.GetDefaultEncoding()), "Now Nord Blue", color.Reset)

	// Switch back to Material
	fmt.Println("\n   Switching back to Material theme...")
	color.UseMaterial()
	fmt.Printf("   New theme: %s\n", color.GetDefaultTheme())
	fmt.Printf("   color.Red:  %s%s%s\n", color.Red.Render(color.GetDefaultEncoding()), "Back to Material Red", color.Reset)
	fmt.Printf("   color.Blue: %s%s%s\n", color.Blue.Render(color.GetDefaultEncoding()), "Back to Material Blue", color.Reset)

	// 6. ENCODING SWITCHING
	fmt.Println("\n6. ENCODING SWITCHING:")
	
	// Switch encoding
	fmt.Println("\n   Current encoding: ANSI")
	fmt.Printf("   color.Material.Blue (ANSI): %s%s%s\n", color.Material.Blue.Render(color.ModeANSI), "ANSI Mode", color.Reset)
	
	fmt.Println("\n   Switching to TrueColor encoding...")
	color.SetDefaultEncoding(color.ModeTrueColor)
	fmt.Printf("   color.Material.Blue (True): %s%s%s\n", color.Material.Blue.Render(color.ModeTrueColor), "TrueColor Mode", color.Reset)

	// 7. ELEGANT API COMPARISON
	fmt.Println("\n7. API ELEGANCE COMPARISON:")
	fmt.Println("\n   🚫 OLD UGLY API (eliminated):")
	fmt.Println("   ❌ color.ANSIColors.Blue")
	fmt.Println("   ❌ color.TrueColorColors.Blue") 
	fmt.Println("   ❌ color.MaterialColors.Blue")
	fmt.Println("   ❌ color.DBlue")
	
	fmt.Println("\n   ✅ NEW CLEAN API (revolutionary):")
	fmt.Println("   ✨ color.ANSI.Blue")
	fmt.Println("   ✨ color.TrueColor.Blue")
	fmt.Println("   ✨ color.Material.Blue")
	fmt.Println("   ✨ color.Blue")

	fmt.Println("\n8. USAGE PATTERNS:")
	fmt.Println("\n   // Building blocks (hardcore)")
	fmt.Println("   color.NewHex(\"#FF5733\")")
	fmt.Println("   color.NewRGB(255, 87, 51)")
	fmt.Println("   color.NewANSI(1)")
	fmt.Println("   color.NewColor256(196)")
	
	fmt.Println("\n   // Explicit encoding (clean!)")
	fmt.Println("   color.ANSI.Blue")
	fmt.Println("   color.TrueColor.Blue")
	fmt.Println("   color.Color256.Blue")
	
	fmt.Println("\n   // Explicit theme (beautiful!)")
	fmt.Println("   color.Material.Blue")
	fmt.Println("   color.Dracula.Blue")
	fmt.Println("   color.Nord.Blue")
	
	fmt.Println("\n   // Default (most elegant!)")
	fmt.Println("   color.Blue")
	fmt.Println("   color.Red")
	fmt.Println("   color.Green")

	fmt.Println("\n✅ Clean Color Encoding System - Revolutionary API Completed!")
	fmt.Printf("🎯 No more ugly names like %sANSIColors%s, %sTrueColorColors%s, %sDBlue%s!\n", 
		color.Red.Render(color.GetDefaultEncoding()), color.Reset,
		color.Red.Render(color.GetDefaultEncoding()), color.Reset,
		color.Red.Render(color.GetDefaultEncoding()), color.Reset)
	fmt.Printf("✨ Welcome to the clean, elegant API: %scolor.ANSI.Blue%s, %scolor.Material.Blue%s, %scolor.Blue%s!\n",
		color.Green.Render(color.GetDefaultEncoding()), color.Reset,
		color.Green.Render(color.GetDefaultEncoding()), color.Reset,
		color.Green.Render(color.GetDefaultEncoding()), color.Reset)
}