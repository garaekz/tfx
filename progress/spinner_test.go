package progress

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

// safeBuffer is a thread-safe wrapper for bytes.Buffer
// to avoid data races in concurrent tests
type safeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (s *safeBuffer) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *safeBuffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}

var spinnerTestMu sync.Mutex

func TestSpinnerBasic(t *testing.T) {
	buf := &safeBuffer{}
	cfg := DefaultSpinnerConfig()
	cfg.Message = "Loading"
	cfg.Writer = buf
	cfg.Frames = []string{"|", "/", "-", "\\"}
	cfg.Interval = 5 * time.Millisecond

	spinner := NewSpinner(cfg)
	spinner.Start()
	time.Sleep(25 * time.Millisecond) // Allow a few frames
	spinner.Stop("Done")

	out := buf.String()
	if !strings.Contains(out, "Loading") {
		t.Errorf("expected 'Loading' in output, got %q", out)
	}
	if !strings.Contains(out, "Done") {
		t.Errorf("expected 'Done' in output, got %q", out)
	}
}

func TestSpinnerExpress(t *testing.T) {
	buf := &safeBuffer{}
	spinner := StartSpinner("Processing")
	spinner.writer = buf
	spinner.Start()
	time.Sleep(10 * time.Millisecond)
	spinner.Stop("Complete")

	out := buf.String()
	if !strings.Contains(out, "Processing") {
		t.Errorf("expected 'Processing' in output, got %q", out)
	}
	if !strings.Contains(out, "Complete") {
		t.Errorf("expected 'Complete' in output, got %q", out)
	}
}

func TestSpinnerWithOptions(t *testing.T) {
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Working"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"⠋", "⠙", "⠹", "⠸"}),
		WithSpinnerInterval(10*time.Millisecond),
	)

	spinner.Start()
	time.Sleep(30 * time.Millisecond)
	spinner.Stop("Finished")

	out := buf.String()
	if !strings.Contains(out, "Working") {
		t.Errorf("expected 'Working' in output, got %q", out)
	}
	if !strings.Contains(out, "Finished") {
		t.Errorf("expected 'Finished' in output, got %q", out)
	}
}

func TestSpinnerFail(t *testing.T) {
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Attempting"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"*"}),
		WithSpinnerInterval(5*time.Millisecond),
	)

	spinner.Start()
	time.Sleep(15 * time.Millisecond)
	spinner.Fail("Error occurred")

	out := buf.String()
	if !strings.Contains(out, "Attempting") {
		t.Errorf("expected 'Attempting' in output, got %q", out)
	}
	if !strings.Contains(out, "Error occurred") {
		t.Errorf("expected 'Error occurred' in output, got %q", out)
	}
}

func TestSpinnerSuccess(t *testing.T) {
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Task"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"+"}),
		WithSpinnerInterval(5*time.Millisecond),
	)

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
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Initial"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"*"}),
		WithSpinnerInterval(5*time.Millisecond),
	)

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
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Test"),
		WithSpinnerWriter(buf),
		WithSpinnerMaterialTheme(),
	)

	spinner.SetTheme(DraculaTheme)
	if spinner.theme != DraculaTheme {
		t.Error("expected theme to be DraculaTheme")
	}
}

func TestSpinnerAlreadyActive(t *testing.T) {
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Test"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"*"}),
		WithSpinnerInterval(5*time.Millisecond),
	)

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
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Test"),
		WithSpinnerWriter(buf),
	)

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
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Test"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"*"}),
		WithSpinnerInterval(5*time.Millisecond),
	)

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
	buf := &safeBuffer{}
	spinner := NewSpinnerWith(
		WithMessage("Rainbow"),
		WithSpinnerWriter(buf),
		WithSpinnerFrames([]string{"*"}),
		WithSpinnerEffect(SpinnerEffectRainbow),
		WithSpinnerInterval(5*time.Millisecond),
	)

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

func TestSpinnerPresetFrames(t *testing.T) {
	tests := []struct {
		name     string
		optFunc  func() SpinnerConfig
		expected []string
	}{
		{
			"Dots",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithDotsFrames()(&cfg)
				return cfg
			},
			[]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		},
		{
			"Arrow",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithArrowFrames()(&cfg)
				return cfg
			},
			[]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
		},
		{
			"Bounce",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithBounceFrames()(&cfg)
				return cfg
			},
			[]string{"⠁", "⠂", "⠄", "⠂"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.optFunc()
			buf := &safeBuffer{}
			cfg.Writer = buf
			if len(cfg.Frames) != len(tt.expected) {
				t.Errorf("expected %d frames, got %d", len(tt.expected), len(cfg.Frames))
			}
			for i, frame := range cfg.Frames {
				if i < len(tt.expected) && frame != tt.expected[i] {
					t.Errorf("expected frame[%d]=%q, got %q", i, tt.expected[i], frame)
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
			buf := &safeBuffer{}
			spinner := tt.constructor("Test Message")
			spinner.writer = buf
			spinner.Start()
			time.Sleep(10 * time.Millisecond)
			spinner.Stop("Done")

			out := buf.String()
			if !strings.Contains(out, "Test Message") {
				t.Errorf("expected 'Test Message' in output, got %q", out)
			}
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
		optFunc  func() SpinnerConfig
		expected ProgressTheme
	}{
		{
			"Material",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithSpinnerMaterialTheme()(&cfg)
				return cfg
			},
			MaterialTheme,
		},
		{
			"Dracula",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithSpinnerDraculaTheme()(&cfg)
				return cfg
			},
			DraculaTheme,
		},
		{
			"Nord",
			func() SpinnerConfig {
				cfg := DefaultSpinnerConfig()
				WithSpinnerNordTheme()(&cfg)
				return cfg
			},
			NordTheme,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.optFunc()
			if cfg.Theme != tt.expected {
				t.Errorf("expected theme %v, got %v", tt.expected, cfg.Theme)
			}
		})
	}
}

func TestSpinnerRainbowOption(t *testing.T) {
	cfg := DefaultSpinnerConfig()
	WithSpinnerRainbow()(&cfg)

	if cfg.Effect != SpinnerEffectRainbow {
		t.Errorf("expected effect SpinnerEffectRainbow, got %v", cfg.Effect)
	}
}
