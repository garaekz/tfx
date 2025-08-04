package main

import (
	"fmt"
	"time"

	"github.com/garaekz/tfx/progress"
)

func runFormFXDemo() {
	// Create a progress bar using the express path.
	progressBar := progress.Start(progress.ProgressConfig{
		Total: 100,
		Label: "Downloading files...",
	})

	fmt.Println("Starting main task...")
	for i := 0; i <= 100; i++ {
		progressBar.Set(i)
		fmt.Print(progressBar.Render())
		time.Sleep(30 * time.Millisecond)
	}

	progressBar.Finish()
	fmt.Println("\nTask completed.")

	// Demonstrate spinner usage.
	spinner := progress.StartSpinner(progress.SpinnerConfig{
		Label: "Processing",
	})
	for i := 0; i < 20; i++ {
		spinner.Tick()
		fmt.Print(spinner.Render())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("\nDone.")
}
