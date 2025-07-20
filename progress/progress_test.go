package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// Test Progress Bar Functionality
func TestProgressBasic(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := DefaultProgressConfig()
	cfg.Total = 4
	cfg.Label = "Test"
	cfg.Writer = buf
	cfg.Width = 4
	p := NewProgress(cfg)
	p.Start()
	p.Set(2)
	p.Complete("Done")
	out := buf.String()
	if !strings.Contains(out, "Test") {
		t.Errorf("expected label 'Test' in output, got %q", out)
	}
	if !strings.Contains(out, "50%") {
		t.Errorf("expected '50%%' in output, got %q", out)
	}
	if !strings.Contains(out, "Done") {
		t.Errorf("expected 'Done' in completion output, got %q", out)
	}
}

func TestProgressExpress(t *testing.T) {
	buf := &bytes.Buffer{}
	p := StartProgress(10, "Loading")
	p.writer = buf
	p.Start()
	p.Set(5)
	p.Complete("Finished")
	out := buf.String()
	if !strings.Contains(out, "Loading") {
		t.Errorf("expected label 'Loading' in output, got %q", out)
	}
	if !strings.Contains(out, "50%") {
		t.Errorf("expected '50%%' in output, got %q", out)
	}
}

func TestProgressWithOptions(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(20),
		WithLabel("Processing"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
	)
	p.Start()
	p.Set(10)
	p.Complete("Success")
	out := buf.String()
	if !strings.Contains(out, "Processing") {
		t.Errorf("expected label 'Processing' in output, got %q", out)
	}
	if !strings.Contains(out, "50%") {
		t.Errorf("expected '50%%' in output, got %q", out)
	}
}

func TestProgressIncrement(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(5), WithProgressWriter(buf))
	p.Start()
	p.Increment()
	p.Increment()
	p.Increment()
	if p.current != 3 {
		t.Errorf("expected current=3, got %d", p.current)
	}
	if p.GetPercent() != 60.0 {
		t.Errorf("expected 60%%, got %.1f%%", p.GetPercent())
	}
}

func TestProgressAdd(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.Add(3)
	p.Add(2)
	if p.current != 5 {
		t.Errorf("expected current=5, got %d", p.current)
	}
	if p.GetPercent() != 50.0 {
		t.Errorf("expected 50%%, got %.1f%%", p.GetPercent())
	}
}

func TestProgressFail(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.Set(3)
	p.Fail("Error occurred")
	out := buf.String()
	if !strings.Contains(out, "Error occurred") {
		t.Errorf("expected failure message in output, got %q", out)
	}
	if !p.done {
		t.Error("expected progress to be marked as done")
	}
}

func TestProgressSetLabel(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.SetLabel("New Label")
	if p.label != "New Label" {
		t.Errorf("expected label='New Label', got %q", p.label)
	}
}

func TestProgressSetTheme(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.SetTheme(DraculaTheme)
	if p.theme != DraculaTheme {
		t.Error("expected theme to be DraculaTheme")
	}
}

func TestProgressSetEffect(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.SetEffect(EffectRainbow)
	if p.effect != EffectRainbow {
		t.Error("expected effect to be EffectRainbow")
	}
	if !p.theme.EffectEnabled {
		t.Error("expected theme.EffectEnabled to be true")
	}
}

func TestProgressGetElapsed(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	
	// Before start, elapsed should be 0
	if p.GetElapsed() != 0 {
		t.Error("expected elapsed to be 0 before start")
	}
	
	p.Start()
	time.Sleep(10 * time.Millisecond)
	elapsed := p.GetElapsed()
	if elapsed < 10*time.Millisecond {
		t.Errorf("expected elapsed >= 10ms, got %v", elapsed)
	}
}

func TestProgressGetETA(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	
	// Before start, ETA should be 0
	if p.GetETA() != 0 {
		t.Error("expected ETA to be 0 before start")
	}
	
	p.Start()
	time.Sleep(10 * time.Millisecond)
	p.Set(5) // 50% complete
	eta := p.GetETA()
	if eta <= 0 {
		t.Errorf("expected positive ETA, got %v", eta)
	}
}

func TestProgressOverflow(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(10), WithProgressWriter(buf))
	p.Start()
	p.Set(15) // Over total
	if p.current != 10 {
		t.Errorf("expected current to be clamped to total=10, got %d", p.current)
	}
}

func TestProgressGlobalAPI(t *testing.T) {
	buf := &bytes.Buffer{}
	
	// Test global progress functions
	StartGlobalProgress(10, "Global")
	// Override writer for testing
	globalProgress.writer = buf
	Set(5)
	Increment()
	Complete("Global Done")
	
	out := buf.String()
	if !strings.Contains(out, "Global Done") {
		t.Errorf("expected 'Global Done' in output, got %q", out)
	}
}

func TestProgressGetPercentZeroTotal(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(WithTotal(0), WithProgressWriter(buf))
	percent := p.GetPercent()
	if percent != 0 {
		t.Errorf("expected 0%% for zero total, got %.1f%%", percent)
	}
}

func TestProgressStyleOptions(t *testing.T) {
	buf := &bytes.Buffer{}
	
	// Test WithGradientEffect
	p1 := NewProgressWith(
		WithTotal(10),
		WithProgressWriter(buf),
		WithGradientEffect(),
	)
	if p1.effect != EffectGradient {
		t.Error("expected EffectGradient")
	}
	
	// Test WithDotsStyle
	p2 := NewProgressWith(
		WithTotal(10),
		WithProgressWriter(buf),
		WithDotsStyle(),
	)
	if p2.style != ProgressStyleDots {
		t.Error("expected ProgressStyleDots")
	}
	
	// Test WithArrowsStyle
	p3 := NewProgressWith(
		WithTotal(10),
		WithProgressWriter(buf),
		WithArrowsStyle(),
	)
	if p3.style != ProgressStyleArrows {
		t.Error("expected ProgressStyleArrows")
	}
}

func TestProgressPresets(t *testing.T) {
	tests := []struct {
		name string
		constructor func(int, string) *Progress
	}{
		{"Material", NewMaterialProgress},
		{"Dracula", NewDraculaProgress},
		{"Nord", NewNordProgress},
		{"Rainbow", NewRainbowProgress},
		{"Minimal", NewMinimalProgress},
		{"Wide", NewWideProgress},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			p := tt.constructor(10, "Test")
			p.writer = buf
			p.Start()
			p.Set(5)
			p.Complete("Done")
			out := buf.String()
			if !strings.Contains(out, "Test") {
				t.Errorf("expected 'Test' in output, got %q", out)
			}
		})
	}
}

func TestSpinnerStop(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := DefaultSpinnerConfig()
	cfg.Message = "Wait"
	cfg.Writer = buf
	cfg.Frames = []string{"*"}
	cfg.Interval = 10 * time.Millisecond
	spin := NewSpinner(cfg)
	spin.Start()
	// allow a few frames to render
	time.Sleep(30 * time.Millisecond)
	spin.Stop("OK")
	out := buf.String()
	// Should contain the stop message
	if !strings.Contains(out, "OK") {
		t.Errorf("expected 'OK' in spinner output, got %q", out)
	}
}
