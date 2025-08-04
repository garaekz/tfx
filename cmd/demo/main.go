package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	version = "1.0.0"
	banner  = `
████████ ███████ ██   ██     ██████  ███████ ███    ███  ██████  
   ██    ██       ██ ██      ██   ██ ██      ████  ████ ██    ██ 
   ██    █████     ███       ██   ██ █████   ██ ████ ██ ██    ██ 
   ██    ██       ██ ██      ██   ██ ██      ██  ██  ██ ██    ██ 
   ██    ██      ██   ██     ██████  ███████ ██      ██  ██████  

TFX Demo - Showcase terminal effects and structured output
`
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "progress", "--progress", "-p":
		// runProgressDemo()
	case "color", "--color", "-c":
		runColorDemo()
	case "color-encoding", "--color-encoding", "-ce":
		runColorEncodingDemo()
	case "spinner", "--spinner", "-s":
		// runSpinnerDemo()
	case "logfx", "--logfx", "-l":
		runlogfxDemo()
	case "formfx", "--formfx", "-f":
		runFormFXDemo()
	case "multipath", "--multipath", "-m":
		// runMultipathDemo()
	case "showcase", "--showcase", "-sh":
		// runShowcaseDemo()
	case "all", "--all", "-a":
		// runAllDemos()
	case "version", "--version", "-v":
		fmt.Printf("TFX Demo v%s\n", version)
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Print(banner)
	fmt.Println("Usage: demo <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  progress, -p     Show progress bar demonstrations")
	fmt.Println("  color, -c        Show color system demonstrations")
	fmt.Println("  color-encoding -ce   Show new color encoding system")
	fmt.Println("  spinner, -s      Show spinner demonstrations")
	fmt.Println("  logfx, -l         Show logging system demonstrations")
	fmt.Println("  formfx, -f       Show interactive form demonstrations")
	fmt.Println("  multipath, -m    Show multipath API demonstrations")
	fmt.Println("  showcase, -sh    Show complete system showcase")
	fmt.Println("  all, -a          Run all demonstrations")
	fmt.Println("  version, -v      Show version information")
	fmt.Println("  help, -h         Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./demo progress")
	fmt.Println("  ./demo --color")
	fmt.Println("  ./demo color-encoding")
	fmt.Println("  ./demo -s")
	fmt.Println("  ./demo showcase")
	fmt.Println("  ./demo logfx")
}
