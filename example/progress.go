package main

import (
	"time"

	"github.com/garaekz/tfx/progress"
)

func main() {
	bar := progress.New(100, "Downloading",
		progress.WithStyle(progress.ProgressStyleDots),
		progress.WithTheme(progress.DraculaTheme),
		progress.WithWidth(30),
	)
	bar.Start()
	for i := 0; i <= 100; i++ {
		bar.Set(i)
		time.Sleep(25 * time.Millisecond)
	}
	bar.Complete("Done!")

	sugarBar := progress.Start(50, "ASCII", progress.WithStyle(progress.ProgressStyleAscii))
	for i := 0; i <= 50; i++ {
		sugarBar.Set(i)
		time.Sleep(20 * time.Millisecond)
	}
	sugarBar.Complete("OK")

	// Example: spinner with default frames and theme
	spin := progress.NewSpinner("Connecting...")
	spin.Start()
	time.Sleep(2 * time.Second)
	spin.Stop("✅ Connected")

	// Example: custom frames and interval
	customFrames := []string{"◐", "◓", "◑", "◒"}
	spin2 := progress.NewSpinner("Processing...",
		progress.WithSpinnerFrames(customFrames),
		progress.WithSpinnerInterval(200*time.Millisecond),
	)
	spin2.Start()
	time.Sleep(1 * time.Second)
	spin2.Fail("❌ Failed")
}
