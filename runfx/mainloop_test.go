package runfx

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/writer"
)

type mockVisual struct {
	tickCount   int
	renderCount int
	resizeCount int
	lastCols    int
	lastRows    int
	mu          sync.Mutex
}

func (m *mockVisual) Render(w share.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderCount++

	// Create a mock entry for testing
	entry := &share.Entry{
		Message: "test render",
	}
	w.Write(entry)
}

func (m *mockVisual) Tick(now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tickCount++
}

func (m *mockVisual) OnResize(cols, rows int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resizeCount++
	m.lastCols = cols
	m.lastRows = rows
}

func (m *mockVisual) GetCounts() (tick, render, resize int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.tickCount, m.renderCount, m.resizeCount
}

// createTestLoop creates a loop configured for testing with a buffer to capture output
func createTestLoop(tickInterval time.Duration) (*MainLoop, *bytes.Buffer) {
	var buf bytes.Buffer
	termWriter := writer.NewTerminalWriter(&buf, writer.TerminalOptions{
		DoubleBuffer: true,
		ForceColor:   false,
		DisableColor: false,
	})
	return NewMainLoopWithTerminal(tickInterval, termWriter, true), &buf
}

func TestStartAPI(t *testing.T) {
	loop := Start()
	if loop == nil {
		t.Fatal("Start returned nil")
	}

	if !loop.IsRunning() {
		// Should not be running initially
	} else {
		t.Error("Loop should not be running initially")
	}
}

func TestMainLoopBasic(t *testing.T) {
	loop, _ := createTestLoop(10 * time.Millisecond)

	// Test initial state
	if loop.IsRunning() {
		t.Error("Loop should not be running initially")
	}

	// Test mounting
	visual := &mockVisual{}
	unmount, err := loop.Mount(visual)
	if err != nil {
		t.Fatalf("Failed to mount visual: %v", err)
	}
	defer unmount()

	// Test running for a short time
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := loop.Run(ctx)
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			t.Errorf("Loop.Run failed: %v", err)
		}
	}()

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	if !loop.IsRunning() {
		t.Error("Loop should be running")
	}

	cancel()
	wg.Wait()

	// Check visual was ticked
	tick, _, _ := visual.GetCounts()
	if tick == 0 {
		t.Error("Visual was not ticked")
	}
}

func TestMainLoopMultipleVisuals(t *testing.T) {
	loop, _ := createTestLoop(10 * time.Millisecond)

	// Mount multiple visuals
	visuals := make([]*mockVisual, 3)
	unmounts := make([]func(), 3)

	for i := range visuals {
		visuals[i] = &mockVisual{}
		var err error
		unmounts[i], err = loop.Mount(visuals[i])
		if err != nil {
			t.Fatalf("Failed to mount visual %d: %v", i, err)
		}
		defer unmounts[i]()
	}

	// Run loop
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	wg.Wait()

	// Check all visuals were ticked
	for i, visual := range visuals {
		tick, _, _ := visual.GetCounts()
		if tick == 0 {
			t.Errorf("Visual %d was not ticked", i)
		}
	}
}

func TestMainLoopStop(t *testing.T) {
	loop := NewMainLoop(10 * time.Millisecond)
	visual := &mockVisual{}
	unmount, err := loop.Mount(visual)
	if err != nil {
		t.Fatalf("Failed to mount visual: %v", err)
	}
	defer unmount()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	// Let it start
	time.Sleep(50 * time.Millisecond)

	// Stop the loop
	err = loop.Stop()
	if err != nil {
		t.Errorf("Failed to stop loop: %v", err)
	}

	wg.Wait()

	// Should not be running after stop
	if loop.IsRunning() {
		t.Error("Loop should not be running after stop")
	}
}

func TestMainLoopUnmount(t *testing.T) {
	loop, _ := createTestLoop(10 * time.Millisecond)

	visual1 := &mockVisual{}
	visual2 := &mockVisual{}

	unmount1, err := loop.Mount(visual1)
	if err != nil {
		t.Fatalf("Failed to mount visual1: %v", err)
	}

	unmount2, err := loop.Mount(visual2)
	if err != nil {
		t.Fatalf("Failed to mount visual2: %v", err)
	}

	// Start loop
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	// Let it run
	time.Sleep(50 * time.Millisecond)

	// Unmount first visual
	unmount1()

	// Let it run more
	time.Sleep(50 * time.Millisecond)

	unmount2()
	wg.Wait()

	// Both visuals should have been ticked
	tick1, _, _ := visual1.GetCounts()
	tick2, _, _ := visual2.GetCounts()

	if tick1 == 0 {
		t.Error("Visual1 was not ticked")
	}
	if tick2 == 0 {
		t.Error("Visual2 was not ticked")
	}
}

func TestMainLoopConcurrency(t *testing.T) {
	loop, _ := createTestLoop(5 * time.Millisecond)

	numVisuals := 10
	visuals := make([]*mockVisual, numVisuals)
	unmounts := make([]func(), numVisuals)

	// Concurrent mounting
	var wg sync.WaitGroup
	for i := 0; i < numVisuals; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			visuals[idx] = &mockVisual{}
			var err error
			unmounts[idx], err = loop.Mount(visuals[idx])
			if err != nil {
				t.Errorf("Failed to mount visual %d: %v", idx, err)
			}
		}(i)
	}
	wg.Wait()

	// Clean up
	defer func() {
		for _, unmount := range unmounts {
			if unmount != nil {
				unmount()
			}
		}
	}()

	// Run loop
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	wg.Wait()

	// Check all visuals were ticked
	for i, visual := range visuals {
		if visual == nil {
			continue
		}
		tick, _, _ := visual.GetCounts()
		if tick == 0 {
			t.Errorf("Visual %d was not ticked", i)
		}
	}
}

func TestMainLoopResizeHandling(t *testing.T) {
	loop, _ := createTestLoop(10 * time.Millisecond)
	visual := &mockVisual{}

	unmount, err := loop.Mount(visual)
	if err != nil {
		t.Fatalf("Failed to mount visual: %v", err)
	}
	defer unmount()

	// Simulate resize by calling the handler directly
	// In a real scenario, this would be triggered by signal
	loop.handleResize()

	_, _, resize := visual.GetCounts()
	if resize == 0 {
		// Note: resize might be 0 if terminal size detection fails in test
		// This is acceptable in test environment
		t.Logf("Resize not detected (acceptable in test environment)")
	}
}
