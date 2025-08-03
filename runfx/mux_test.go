package runfx

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/garaekz/tfx/internal/share"
)

type dummyVisual struct {
	id string
}

func (d *dummyVisual) Render(_ share.Writer) {}
func (d *dummyVisual) Tick(_ time.Time)      {}
func (d *dummyVisual) OnResize(_, _ int)     {}

func TestMultiplexerBasic(t *testing.T) {
	mux := NewMultiplexer()
	v := &dummyVisual{id: "test"}

	// Test mounting
	mux.Mount("v1", v)
	if len(mux.ListVisuals()) != 1 {
		t.Error("Visual not mounted")
	}

	// Test retrieval
	retrieved, ok := mux.GetVisual("v1")
	if !ok {
		t.Error("Visual not found after mounting")
	}
	if retrieved != v {
		t.Error("Retrieved visual is not the same as mounted")
	}

	// Test unmounting
	mux.Unmount("v1")
	if len(mux.ListVisuals()) != 0 {
		t.Error("Visual not unmounted")
	}

	// Test retrieval after unmount
	_, ok = mux.GetVisual("v1")
	if ok {
		t.Error("Visual still found after unmounting")
	}
}

func TestMultiplexerConcurrency(t *testing.T) {
	mux := NewMultiplexer()
	numGoroutines := 10
	numVisuals := 5

	var wg sync.WaitGroup

	// Concurrent mounting
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numVisuals; j++ {
				name := fmt.Sprintf("visual_%d_%d", id, j)
				v := &dummyVisual{id: name}
				mux.Mount(name, v)
			}
		}(i)
	}

	wg.Wait()

	// Check all visuals were mounted
	expected := numGoroutines * numVisuals
	if len(mux.ListVisuals()) != expected {
		t.Errorf("Expected %d visuals, got %d", expected, len(mux.ListVisuals()))
	}

	// Concurrent access and unmounting
	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			// Read operations
			for j := 0; j < numVisuals; j++ {
				name := fmt.Sprintf("visual_%d_%d", id, j)
				mux.GetVisual(name)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			// Unmount operations
			for j := 0; j < numVisuals; j++ {
				name := fmt.Sprintf("visual_%d_%d", id, j)
				mux.Unmount(name)
			}
		}(i)
	}

	wg.Wait()

	// All should be unmounted
	if len(mux.ListVisuals()) != 0 {
		t.Errorf("Expected 0 visuals after unmounting, got %d", len(mux.ListVisuals()))
	}
}

func TestMultiplexerReuseGap(t *testing.T) {
	mux := NewMultiplexer()
	v1 := &dummyVisual{id: "v1"}
	v2 := &dummyVisual{id: "v2"}

	// Mount first visual
	mux.Mount("test1", v1)

	// Unmount to create a gap
	mux.Unmount("test1")

	// Try to reuse gap
	name, reused := mux.ReuseGap(v2)
	if !reused {
		t.Error("Gap was not reused")
	}
	if name == "" {
		t.Error("No name returned for reused gap")
	}

	// Verify the visual is accessible
	retrieved, ok := mux.GetVisual(name)
	if !ok {
		t.Error("Reused visual not accessible")
	}
	if retrieved != v2 {
		t.Error("Reused visual is not the correct one")
	}
}

func TestMultiplexerDuplicateMount(t *testing.T) {
	mux := NewMultiplexer()
	v1 := &dummyVisual{id: "v1"}
	v2 := &dummyVisual{id: "v2"}

	// Mount first visual
	mux.Mount("test", v1)

	// Try to mount different visual with same name
	mux.Mount("test", v2)

	// Should still have only one visual
	if len(mux.ListVisuals()) != 1 {
		t.Errorf("Expected 1 visual after duplicate mount, got %d", len(mux.ListVisuals()))
	}

	// Should be the original visual (not replaced)
	retrieved, ok := mux.GetVisual("test")
	if !ok {
		t.Error("Visual not found")
	}
	if retrieved != v1 {
		t.Error("Visual was replaced, should have been kept")
	}
}
