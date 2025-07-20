package main

import (
	"fmt"
	"strings"
	"time"
)

func runAllDemos() {
	fmt.Print(banner)
	fmt.Println("Running all TFX demonstrations...")
	fmt.Println("==================================")

	// Run each demo with separators
	runProgressDemo()
	
	time.Sleep(1 * time.Second)
	fmt.Println("\n" + strings.Repeat("─", 50))
	
	runSpinnerDemo()
	
	time.Sleep(1 * time.Second)
	fmt.Println("\n" + strings.Repeat("─", 50))
	
	runColorDemo()
	
	time.Sleep(1 * time.Second)
	fmt.Println("\n" + strings.Repeat("─", 50))
	
	runLogxDemo()
	
	time.Sleep(1 * time.Second)
	fmt.Println("\n" + strings.Repeat("─", 50))
	
	runMultipathDemo()

	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("🎉 ALL TFX DEMONSTRATIONS COMPLETED! 🎉")
	fmt.Println(strings.Repeat("═", 50))
}