package main

import (
    "fmt"
    "time"

    // progrefx replaces the legacy progress package.  It provides
    // progress bars and spinners integrated with runfx.
    "github.com/garaekz/tfx/progrefx"
)

func runFormFXDemo() {
	// Create a progress bar using the express path.
    progressBar := progrefx.Start(progrefx.ProgressConfig{
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
    spinner := progrefx.StartSpinner(progrefx.SpinnerConfig{
		Label: "Processing",
    })
	for i := 0; i < 20; i++ {
		spinner.Tick()
		fmt.Print(spinner.Render())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("\nDone.")
}
