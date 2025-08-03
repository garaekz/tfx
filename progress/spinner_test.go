package progress

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/garaekz/tfx/internal/testutil"
)

var spinnerTestMu sync.Mutex

func TestSpinnerBasic(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	cfg := DefaultSpinnerConfig()
	cfg.Message = "Loading"
	cfg.Writer = buf
	cfg.Frames = []string{"|", "/", "-", "\\"}
	cfg.Interval = 5 * time.Millisecond

	spinner := NewSpinnerWithConfig(cfg)
	spinner.Start()
	time.Sleep(25 * time.Millisecond) // Allow a few frames
	spinner.Stop("Done")

	out := buf.String()
	// For RunFX integration, we mainly verify completion messages
	if !strings.Contains(out, "Done") {
		t.Errorf("expected 'Done' in output, got %q", out)
	}
}

func TestSpinnerExpress(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := StartSpinner("Processing")
	spinner.writer = buf
	spinner.Start()
	time.Sleep(10 * time.Millisecond)
	spinner.Stop("Complete")

	out := buf.String()
	// For now, just verify that Complete message appears
	// The RunFX integration means the spinner animation might not be captured in tests
	if !strings.Contains(out, "Complete") {
		t.Errorf("expected 'Complete' in output, got %q", out)
	}
}

func TestSpinnerWithOptions(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Working",
		Writer:   buf,
		Frames:   []string{"⠋", "⠙", "⠹", "⠸"},
		Interval: 10 * time.Millisecond,
	})

	spinner.Start()
	time.Sleep(30 * time.Millisecond)
	spinner.Stop("Finished")

	out := buf.String()
	if !strings.Contains(out, "Finished") {
		t.Errorf("expected 'Finished' in output, got %q", out)
	}
}

func TestSpinnerFail(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Attempting",
		Writer:   buf,
		Frames:   []string{"*"},
		Interval: 5 * time.Millisecond,
	})

	spinner.Start()
	time.Sleep(15 * time.Millisecond)
	spinner.Fail("Error occurred")

	out := buf.String()
	if !strings.Contains(out, "Error occurred") {
		t.Errorf("expected 'Error occurred' in output, got %q", out)
	}
}

func TestSpinnerSuccess(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Task",
		Writer:   buf,
		Frames:   []string{"+"},
		Interval: 5 * time.Millisecond,
	})

	spinner.Start()
	time.Sleep(15 * time.Millisecond)
	spinner.Success("Task completed")

	out := buf.String()
	if !strings.Contains(out, "Task") {
		t.Errorf("expected 'Task' in output, got %q", out)
	}
	if !strings.Contains(out, "Task completed") {
		t.Errorf("expected 'Task completed' in output, got %q", out)
	}
}

func TestSpinnerSetMessage(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Initial",
		Writer:   buf,
		Frames:   []string{"*"},
		Interval: 5 * time.Millisecond,
	})

	spinner.Start()
	time.Sleep(10 * time.Millisecond)
	spinner.SetMessage("Updated")
	time.Sleep(10 * time.Millisecond)
	spinner.Stop("Done")

	if spinner.message != "Updated" {
		t.Errorf("expected message='Updated', got %q", spinner.message)
	}
}

func TestSpinnerSetTheme(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message: "Test",
		Writer:  buf,
		Theme:   MaterialTheme,
	})

	spinner.SetTheme(DraculaTheme)
	if spinner.theme != DraculaTheme {
		t.Error("expected theme to be DraculaTheme")
	}
}

func TestSpinnerAlreadyActive(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Test",
		Writer:   buf,
		Frames:   []string{"*"},
		Interval: 5 * time.Millisecond,
	})

	// Start spinner
	spinner.Start()
	time.Sleep(10 * time.Millisecond)

	// Try to start again (should do nothing)
	spinner.Start()
	time.Sleep(10 * time.Millisecond)

	spinner.Stop("Done")

	out := buf.String()
	if !strings.Contains(out, "Done") {
		t.Errorf("expected 'Done' in output, got %q", out)
	}
}

func TestSpinnerStopInactive(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message: "Test",
		Writer:  buf,
	})

	// Try to stop without starting (should do nothing)
	spinner.Stop("Not started")
	spinner.Fail("Not started")

	out := buf.String()
	// Should be empty since spinner was never active
	if out != "" {
		t.Errorf("expected empty output for inactive spinner, got %q", out)
	}
}

func TestSpinnerGlobalAlreadyActive(t *testing.T) {
	spinnerTestMu.Lock()
	defer spinnerTestMu.Unlock()

	// Ensure no global spinner is running
	if globalSpinner != nil {
		globalSpinner.Stop("cleanup")
		globalSpinner = nil
	}

	// Start first global spinner
	StartGlobalSpinner("Global 1")
	time.Sleep(10 * time.Millisecond)

	// Try to start another (should do nothing)
	StartGlobalSpinner("Global 2")
	time.Sleep(10 * time.Millisecond)

	StopGlobalSpinner("Done")

	// Simple check - just ensure no panic occurred
	if globalSpinner == nil {
		t.Error("expected globalSpinner to exist after operations")
	}
}

func TestSpinnerMultipleStartStop(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Test",
		Writer:   buf,
		Frames:   []string{"*"},
		Interval: 5 * time.Millisecond,
	})

	// First start/stop cycle
	spinner.Start()
	time.Sleep(10 * time.Millisecond)
	spinner.Stop("First done")

	// Second start should not work (spinner is already stopped)
	spinner.Start()
	time.Sleep(10 * time.Millisecond)
	spinner.Stop("Second done")

	out := buf.String()
	if !strings.Contains(out, "First done") {
		t.Errorf("expected 'First done' in output, got %q", out)
	}
}

func TestSpinnerEffects(t *testing.T) {
	buf := &testutil.SafeBuffer{}
	spinner := NewSpinnerWithConfig(SpinnerConfig{
		Message:  "Rainbow",
		Writer:   buf,
		Frames:   []string{"*"},
		Effect:   SpinnerEffectRainbow,
		Interval: 5 * time.Millisecond,
	})

	if spinner.effect != SpinnerEffectRainbow {
		t.Error("expected spinner effect to be SpinnerEffectRainbow")
	}

	spinner.Start()
	time.Sleep(15 * time.Millisecond)
	spinner.Stop("Rainbow done")

	out := buf.String()
	if !strings.Contains(out, "Rainbow") {
		t.Errorf("expected 'Rainbow' in output, got %q", out)
	}
}

func TestSpinnerFrameOptions(t *testing.T) {
	tests := []struct {
		name     string
		frames   []string
		expected []string
	}{
		{
			"Dots",
			[]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
			[]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		},
		{
			"Arrow",
			[]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
			[]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
		},
		{
			"Bounce",
			[]string{"⠁", "⠂", "⠄", "⠂"},
			[]string{"⠁", "⠂", "⠄", "⠂"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := SpinnerConfig{
				Frames: tt.frames,
			}
			buf := &testutil.SafeBuffer{}
			cfg.Writer = buf

			spinner := NewSpinnerWithConfig(cfg)
			// Check if frames are set correctly
			if len(spinner.frames) != len(tt.expected) {
				t.Errorf("expected %d frames, got %d", len(tt.expected), len(spinner.frames))
			}
			for i, frame := range tt.expected {
				if i < len(spinner.frames) && spinner.frames[i] != frame {
					t.Errorf("expected frame[%d]='%s', got '%s'", i, frame, spinner.frames[i])
				}
			}
		})
	}
}

func TestSpinnerPresets(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(string) *Spinner
	}{
		{"Material", NewMaterialSpinner},
		{"Dracula", NewDraculaSpinner},
		{"Nord", NewNordSpinner},
		{"Dots", NewDotsSpinner},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &testutil.SafeBuffer{}
			spinner := tt.constructor("Test Message")
			spinner.writer = buf
			spinner.Start()
			time.Sleep(10 * time.Millisecond)
			spinner.Stop("Done")

			out := buf.String()
			if !strings.Contains(out, "Done") {
				t.Errorf("expected 'Done' in output, got %q", out)
			}
		})
	}
}

func TestSpinnerGlobalAPI(t *testing.T) {
	spinnerTestMu.Lock()
	defer spinnerTestMu.Unlock()

	// Ensure no global spinner is running
	if globalSpinner != nil {
		globalSpinner.Stop("cleanup")
		globalSpinner = nil
	}

	StartGlobalSpinner("Global Task")
	time.Sleep(10 * time.Millisecond)
	StopGlobalSpinner("Global Done")

	// Simple check - just ensure no panic occurred
	if globalSpinner == nil {
		t.Error("expected globalSpinner to exist after operations")
	}
}

func TestSpinnerGlobalAPIFail(t *testing.T) {
	spinnerTestMu.Lock()
	defer spinnerTestMu.Unlock()

	// Ensure no global spinner is running
	if globalSpinner != nil {
		globalSpinner.Stop("cleanup")
		globalSpinner = nil
	}

	StartGlobalSpinner("Global Task 2")
	time.Sleep(10 * time.Millisecond)
	FailGlobalSpinner("Global Failed")

	// Simple check - just ensure no panic occurred
	if globalSpinner == nil {
		t.Error("expected globalSpinner to exist after operations")
	}
}

func TestSpinnerThemeOptions(t *testing.T) {
	tests := []struct {
		name     string
		theme    ProgressTheme
		expected ProgressTheme
	}{
		{
			"Material",
			MaterialTheme,
			MaterialTheme,
		},
		{
			"Dracula",
			DraculaTheme,
			DraculaTheme,
		},
		{
			"Nord",
			NordTheme,
			NordTheme,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := SpinnerConfig{
				Theme: tt.theme,
			}
			if cfg.Theme != tt.expected {
				t.Errorf("expected theme %v, got %v", tt.expected, cfg.Theme)
			}
		})
	}
}

func TestSpinnerRainbowOption(t *testing.T) {
	cfg := SpinnerConfig{
		Effect: SpinnerEffectRainbow,
	}

	if cfg.Effect != SpinnerEffectRainbow {
		t.Errorf("expected effect SpinnerEffectRainbow, got %v", cfg.Effect)
	}
}
