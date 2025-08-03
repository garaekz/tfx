package progress

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockWriter is a simple mock for io.Writer to capture output.
type mockWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return mw.buf.Write(p)
}

func (mw *mockWriter) String() string {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return mw.buf.String()
}

func (mw *mockWriter) Reset() {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.buf.Reset()
}

// TestProgress_Start tests the Start method.
func TestProgress_Start(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 10
	cfg.Label = "Starting"
	p := newProgress(cfg)

	p.Start()

	if !p.isStarted {
		t.Error("Expected progress to be started")
	}
	output := mw.String()
	if !strings.Contains(output, "Starting") {
		t.Errorf("Expected output to contain label, got %q", output)
	}

	// Test starting again (should do nothing)
	mw.Reset()
	p.Start()
	output = mw.String()
	if output != "" {
		t.Errorf("Expected no output on second start, got %q", output)
	}
}

// TestProgress_Set tests the Set method.
func TestProgress_Set(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 100
	cfg.Label = "Setting"
	p := newProgress(cfg)

	p.Set(50)
	if p.current != 50 {
		t.Errorf("Expected current to be 50, got %d", p.current)
	}
	output := mw.String()
	if !strings.Contains(output, "Setting") || !strings.Contains(output, "50%") {
		t.Errorf("Expected output to reflect 50%% progress, got %q", output)
	}

	// Test setting beyond total
	mw.Reset()
	p.Set(150)
	if p.current != 100 {
		t.Errorf("Expected current to be 100, got %d", p.current)
	}
	output = mw.String()
	if !strings.Contains(output, "100%") {
		t.Errorf("Expected output to reflect 100%% progress, got %q", output)
	}

	// Test setting before start
	mw.Reset()
	cfg2 := DefaultProgressConfig()
	cfg2.Writer = mw
	cfg2.Total = 10
	p2 := newProgress(cfg2)
	p2.Set(5)
	if !p2.isStarted {
		t.Error("Expected progress to be started after Set")
	}
}

// TestProgress_Increment tests the Increment method.
func TestProgress_Increment(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 5
	p := newProgress(cfg)

	p.Set(0)
	mw.Reset()
	p.Increment()
	if p.current != 1 {
		t.Errorf("Expected current to be 1, got %d", p.current)
	}
	output := mw.String()
	if !strings.Contains(output, "20%") {
		t.Errorf("Expected output to reflect 20%% progress, got %q", output)
	}

	// Increment past total
	mw.Reset()
	p.Set(4)
	p.Increment()
	p.Increment()
	if p.current != 5 {
		t.Errorf("Expected current to be 5, got %d", p.current)
	}
	output = mw.String()
	if !strings.Contains(output, "100%") {
		t.Errorf("Expected output to reflect 100%% progress, got %q", output)
	}
}

// TestProgress_Add tests the Add method.
func TestProgress_Add(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 10
	p := newProgress(cfg)

	p.Set(0)
	mw.Reset()
	p.Add(3)
	if p.current != 3 {
		t.Errorf("Expected current to be 3, got %d", p.current)
	}
	output := mw.String()
	if !strings.Contains(output, "30%") {
		t.Errorf("Expected output to reflect 30%% progress, got %q", output)
	}

	// Add past total
	mw.Reset()
	p.Set(8)
	p.Add(5)
	if p.current != 10 {
		t.Errorf("Expected current to be 10, got %d", p.current)
	}
	output = mw.String()
	if !strings.Contains(output, "100%") {
		t.Errorf("Expected output to reflect 100%% progress, got %q", output)
	}
}

// TestProgress_Complete tests the Complete method.
func TestProgress_Complete(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 100
	cfg.Label = "Completing"
	p := newProgress(cfg)

	p.Set(50)
	mw.Reset()
	p.Complete("Done!")

	if p.current != 100 {
		t.Errorf("Expected current to be 100, got %d", p.current)
	}
	output := mw.String()
	if !strings.Contains(output, "Done!") || !strings.Contains(output, "100%") ||
		!strings.Contains(output, "✅") {
		t.Errorf("Expected output to contain completion message and 100%%, got %q", output)
	}
}

// TestProgress_Fail tests the Fail method.
func TestProgress_Fail(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 100
	cfg.Label = "Failing"
	p := newProgress(cfg)

	p.Set(50)
	mw.Reset()
	p.Fail("Failed!")

	if p.current != 100 {
		t.Errorf("Expected current to be 100, got %d", p.current)
	}
	output := mw.String()
	if !strings.Contains(output, "Failed!") || !strings.Contains(output, "100%") ||
		!strings.Contains(output, "❌") {
		t.Errorf("Expected output to contain failure message and 100%%, got %q", output)
	}
}

// TestProgress_SetLabel tests the SetLabel method.
func TestProgress_SetLabel(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	p.SetLabel("New Label")
	if p.label != "New Label" {
		t.Errorf("Expected label to be 'New Label', got %q", p.label)
	}

	// Test setting label when started
	p.Start()
	mw.Reset()
	p.SetLabel("Another Label")
	output := mw.String()
	if !strings.Contains(output, "Another Label") {
		t.Errorf("Expected output to contain new label when started, got %q", output)
	}
}

// TestProgress_SetTheme tests the SetTheme method.
func TestProgress_SetTheme(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	newTheme := DraculaTheme
	p.SetTheme(newTheme)
	if p.theme != newTheme {
		t.Errorf("Expected theme to be DraculaTheme, got %v", p.theme)
	}

	// Test setting theme when started
	p.Start()
	mw.Reset()
	p.SetTheme(NordTheme)
	output := mw.String()
	// Cannot easily assert theme change from output, but ensure no crash
	if output == "" {
		// This is just to use output and prevent unused variable error
	}
}

// TestProgress_SetEffect tests the SetEffect method.
func TestProgress_SetEffect(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	p.SetEffect(EffectRainbow)
	if p.effect != EffectRainbow {
		t.Errorf("Expected effect to be EffectRainbow, got %v", p.effect)
	}
	if !p.theme.EffectEnabled {
		t.Error("Expected effect to be enabled")
	}

	// Test setting effect to None
	p.SetEffect(EffectNone)
	if p.effect != EffectNone {
		t.Errorf("Expected effect to be EffectNone, got %v", p.effect)
	}
	if p.theme.EffectEnabled {
		t.Error("Expected effect to be disabled")
	}
}

// TestProgress_GetPercent tests the GetPercent method.
func TestProgress_GetPercent(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 100
	p := newProgress(cfg)

	p.Set(0)
	if p.GetPercent() != 0 {
		t.Errorf("Expected 0%%, got %f", p.GetPercent())
	}

	p.Set(50)
	if p.GetPercent() != 50 {
		t.Errorf("Expected 50%%, got %f", p.GetPercent())
	}

	p.Set(100)
	if p.GetPercent() != 100 {
		t.Errorf("Expected 100%%, got %f", p.GetPercent())
	}

	// Test with total 0
	cfg2 := DefaultProgressConfig()
	cfg2.Writer = mw
	cfg2.Total = 0
	p2 := newProgress(cfg2)
	if p2.GetPercent() != 0 {
		t.Errorf("Expected 0%% for total 0, got %f", p2.GetPercent())
	}
}

// TestProgress_GetElapsed tests the GetElapsed method.
func TestProgress_GetElapsed(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	// Before start
	if p.GetElapsed() != 0 {
		t.Errorf("Expected 0 elapsed before start, got %v", p.GetElapsed())
	}

	p.Start()
	time.Sleep(100 * time.Millisecond)
	elapsed := p.GetElapsed()
	if elapsed < 100*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Expected elapsed time around 100ms, got %v", elapsed)
	}
}

// TestProgress_GetETA tests the GetETA method.
func TestProgress_GetETA(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	cfg.Total = 100
	cfg.ShowETA = true
	p := newProgress(cfg)

	// Before start or current is 0
	if p.GetETA() != 0 {
		t.Errorf("Expected 0 ETA before start or current is 0, got %v", p.GetETA())
	}

	p.Start()
	time.Sleep(50 * time.Millisecond)
	p.Set(10) // 10% done in 50ms
	// Rate: 10 / 0.05s = 200 units/sec
	// Remaining: 90 units
	// ETA: 90 / 200 = 0.45s = 450ms
	eta := p.GetETA()
	if eta < 400*time.Millisecond || eta > 500*time.Millisecond {
		t.Errorf("Expected ETA around 450ms, got %v", eta)
	}

	// Test with total 0
	cfg2 := DefaultProgressConfig()
	cfg2.Writer = mw
	cfg2.Total = 0
	p2 := newProgress(cfg2)
	if p2.GetETA() != 0 {
		t.Errorf("Expected 0 ETA for total 0, got %v", p2.GetETA())
	}
}

// TestProgress_Redraw tests the Redraw method.
func TestProgress_Redraw(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	// Redraw before start (should do nothing)
	p.Redraw()
	if mw.String() != "" {
		t.Errorf("Expected no output on redraw before start, got %q", mw.String())
	}

	// Redraw after start
	p.Start()
	mw.Reset()
	p.Redraw()
	output := mw.String()
	if !strings.Contains(output, p.label) {
		t.Errorf("Expected output on redraw after start, got %q", output)
	}
}

// TestProgress_Close tests the Close method.
func TestProgress_Close(t *testing.T) {
	mw := &mockWriter{}
	cfg := DefaultProgressConfig()
	cfg.Writer = mw
	p := newProgress(cfg)

	// Close before start (should not panic)
	p.Close()

	// Close after start
	p.Start()
	p.Close()
	// No direct assertion for Close, but ensure no panics or errors
}

// TestProgress_Constructors tests all constructor functions.
func TestProgress_Constructors(t *testing.T) {
	mw := &mockWriter{}

	// Test Start()
	p1 := Start()
	if p1 == nil || !p1.isStarted {
		t.Error("Start() failed to create and start progress")
	}
	p1.Close()

	// Test StartWith() - Now we'll use the DSL instead
	p2 := New().Label("Fluent").Total(50).Writer(mw).Start()
	if p2 == nil || !p2.isStarted || p2.label != "Fluent" || p2.total != 50 || p2.writer != mw {
		t.Error("New().Label().Total().Writer().Start() failed to create and start progress")
	}
	p2.Close()

	// Test Start(config)
	cfg3 := DefaultProgressConfig()
	cfg3.Label = "Configured"
	cfg3.Total = 25
	cfg3.Writer = mw
	p3 := Start(cfg3)
	if p3 == nil || !p3.isStarted || p3.label != "Configured" || p3.total != 25 || p3.writer != mw {
		t.Error("Start(config) failed to create and start progress with config")
	}
	p3.Close()

	// Test New() builder without starting
	p4Builder := New()
	if p4Builder == nil {
		t.Error("New() failed to create progress builder")
	}

	// Build without starting
	p4 := p4Builder.Build()
	if p4 == nil || p4.isStarted {
		t.Error("Build() failed to create progress or started unexpectedly")
	}

	// Test using config struct directly with newProgress
	cfg5 := DefaultProgressConfig()
	cfg5.Label = "New Config"
	cfg5.Total = 75
	cfg5.Writer = mw
	p5 := newProgress(cfg5)
	if p5 == nil || p5.isStarted || p5.label != "New Config" || p5.total != 75 || p5.writer != mw {
		t.Error("newProgress() failed to create progress with config")
	}
	p5.Close()

	// Test DSL builder
	p6 := New().Label("New Fluent").Total(90).Writer(mw).Build()
	if p6 == nil || p6.isStarted || p6.label != "New Fluent" || p6.total != 90 || p6.writer != mw {
		t.Error("New().Label().Total().Writer().Build() failed to create progress")
	}
	p6.Close()

	// Test DSL with themes
	p7 := New().Total(10).Label("Simple").Build()
	if p7 == nil || p7.total != 10 || p7.label != "Simple" {
		t.Error("New().Total().Label().Build() failed")
	}
	p7.Close()

	p8 := New().Total(20).Label("Material").MaterialTheme().Build()
	if p8 == nil || p8.total != 20 || p8.label != "Material" || p8.theme != MaterialTheme {
		t.Error("New().MaterialTheme() failed")
	}
	p8.Close()

	p9 := New().Total(30).Label("Dracula").DraculaTheme().Build()
	if p9 == nil || p9.total != 30 || p9.label != "Dracula" || p9.theme != DraculaTheme {
		t.Errorf("New().DraculaTheme() failed: %+v", p9)
	}
	p9.Close()

	p10 := New().Total(40).Label("Nord").NordTheme().Build()
	if p10 == nil || p10.total != 40 || p10.label != "Nord" || p10.theme != NordTheme {
		t.Error("New().NordTheme() failed")
	}
	p10.Close()

	// Test with effects using DSL
	p11 := New().Total(50).Label("Rainbow").Effect(EffectRainbow).Build()
	if p11 == nil || p11.total != 50 || p11.label != "Rainbow" || p11.effect != EffectRainbow {
		t.Error("New().Effect() failed")
	}
	p11.Close()
}

// TestBasicProgressFunctionality tests basic progress functionality.
func TestBasicProgressFunctionality(t *testing.T) {
	mw := &mockWriter{}

	// Test creating progress with Start() function
	p := Start()
	if p == nil {
		t.Error("Start() failed to create progress")
	}
	p.Close()

	// Test creating progress with config
	cfg := DefaultProgressConfig()
	cfg.Total = 10
	cfg.Label = "Global Task"
	cfg.Writer = mw

	p2 := Start(cfg)
	if p2 == nil || !p2.isStarted {
		t.Error("Start(config) failed to create and start progress")
	}

	p2.Set(5)
	output := mw.String()
	if !strings.Contains(output, "Global Task") {
		t.Errorf("Expected output to contain 'Global Task', got %q", output)
	}

	p2.Close()
}
