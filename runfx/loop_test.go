package runfx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/garaekz/tfx/internal/share"
)

type testVisual struct {
	tickCount   int
	renderCount int
	resizeCount int
	lastCols    int
	lastRows    int
	mu          sync.Mutex
}

func (v *testVisual) Render(_ share.Writer) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.renderCount++
}

func (v *testVisual) Tick(_ time.Time) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.tickCount++
}

func (v *testVisual) OnResize(cols, rows int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.resizeCount++
	v.lastCols = cols
	v.lastRows = rows
}

func (v *testVisual) GetCounts() (tick, render, resize int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.tickCount, v.renderCount, v.resizeCount
}

func TestEventLoopBasic(t *testing.T) {
	loop := NewEventLoop(10 * time.Millisecond)
	v := &testVisual{}
	loop.Mount(v)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	wg.Wait()

	tick, _, _ := v.GetCounts()
	if tick == 0 {
		t.Error("Visual did not tick")
	}

	// In test environments timing can be unpredictable, so just verify it ticked
	t.Logf("Got %d ticks in 50ms (interval: 10ms)", tick)
}

func TestEventLoopConcurrency(t *testing.T) {
	loop := NewEventLoop(5 * time.Millisecond)
	visuals := make([]*testVisual, 5)

	// Mount multiple visuals
	for i := range visuals {
		visuals[i] = &testVisual{}
		loop.Mount(visuals[i])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	wg.Wait()

	// Check all visuals ticked
	for i, v := range visuals {
		tick, _, _ := v.GetCounts()
		if tick == 0 {
			t.Errorf("Visual %d did not tick", i)
		}
	}
}

func TestEventLoopStop(t *testing.T) {
	loop := NewEventLoop(10 * time.Millisecond)
	v := &testVisual{}
	loop.Mount(v)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	// Stop the loop after a short time
	time.Sleep(50 * time.Millisecond)
	loop.Stop()

	wg.Wait()

	tick, _, _ := v.GetCounts()
	if tick == 0 {
		t.Error("Visual did not tick before stop")
	}
}

func TestEventLoopDirtyFlags(t *testing.T) {
	loop := NewEventLoop(10 * time.Millisecond)
	v := &testVisual{}
	loop.Mount(v)

	// Mark as dirty
	loop.MarkDirty(0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop.Run(ctx)
	}()

	wg.Wait()

	tick, _, _ := v.GetCounts()
	if tick == 0 {
		t.Error("Visual did not tick when marked dirty")
	}
}
