package flowfx

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	name := "Test Task"
	runFunc := func(ctx context.Context, tracker *ProgressTracker) error { return nil }
	var progress int64

	task := NewTask(
		WithName(name),
		WithRun(runFunc),
		WithProgress(&progress),
	)

	if task.Name != name {
		t.Errorf("Expected task name %q, got %q", name, task.Name)
	}
	if task.Run == nil {
		t.Error("Expected Run function to be set, got nil")
	}
	if task.progressPtr != &progress {
		t.Error("Expected progress pointer to be set, got different pointer")
	}
}

func TestSequence(t *testing.T) {
	t.Run("Successful sequence", func(t *testing.T) {
		var order []string
		mu := &sync.Mutex{}

		task1 := NewTask(
			WithName("Task 1"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				mu.Lock()
				order = append(order, "Task 1")
				mu.Unlock()
				return nil
			}),
		)
		task2 := NewTask(
			WithName("Task 2"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				mu.Lock()
				order = append(order, "Task 2")
				mu.Unlock()
				return nil
			}),
		)

		err := Sequence(context.Background(), task1, task2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedOrder := []string{"Task 1", "Task 2"}
		if !compareStringSlices(order, expectedOrder) {
			t.Errorf("Expected order %v, got %v", expectedOrder, order)
		}
	})

	t.Run("Sequence with error", func(t *testing.T) {
		var order []string
		mu := &sync.Mutex{}
		expectedErr := errors.New("Task 1 failed")

		task1 := NewTask(
			WithName("Task 1"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				mu.Lock()
				order = append(order, "Task 1")
				mu.Unlock()
				return expectedErr
			}),
		)
		task2 := NewTask(
			WithName("Task 2"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				mu.Lock()
				order = append(order, "Task 2")
				mu.Unlock()
				return nil
			}),
		)

		err := Sequence(context.Background(), task1, task2)
		if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %q, got %v", expectedErr, err)
		}

		expectedOrder := []string{"Task 1"}
		if !compareStringSlices(order, expectedOrder) {
			t.Errorf("Expected order %v, got %v", expectedOrder, order)
		}
	})

	t.Run("Sequence with progress", func(t *testing.T) {
		var progress int64
		task := NewTask(
			WithName("Progress Task"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				for i := 0; i <= 10; i++ {
					atomic.StoreInt64(tracker.Current, int64(i))
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			}),
			WithProgress(&progress),
		)

		err := Sequence(context.Background(), task)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if atomic.LoadInt64(&progress) != 10 {
			t.Errorf("Expected final progress to be 10, got %d", atomic.LoadInt64(&progress))
		}
	})
}

func TestParallel(t *testing.T) {
	t.Run("Successful parallel", func(t *testing.T) {
		var order []string
		mu := &sync.Mutex{}

		task1 := NewTask(
			WithName("Task 1"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				order = append(order, "Task 1")
				mu.Unlock()
				return nil
			}),
		)
		task2 := NewTask(
			WithName("Task 2"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				mu.Lock()
				order = append(order, "Task 2")
				mu.Unlock()
				return nil
			}),
		)

		err := Parallel(context.Background(), task1, task2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(order) != 2 {
			t.Errorf("Expected 2 tasks to run, got %d", len(order))
		}
		// Order is not guaranteed in parallel, just check if both ran
		if !containsString(order, "Task 1") || !containsString(order, "Task 2") {
			t.Errorf("Expected both tasks to be in order, got %v", order)
		}
	})

	t.Run("Parallel with errors", func(t *testing.T) {
		expectedErr1 := errors.New("Task 1 failed")
		expectedErr2 := errors.New("Task 2 failed")

		task1 := NewTask(
			WithName("Task 1"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				return expectedErr1
			}),
		)
		task2 := NewTask(
			WithName("Task 2"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				return expectedErr2
			}),
		)

		err := Parallel(context.Background(), task1, task2)
		if err == nil {
			t.Error("Expected errors, got nil")
		}

		// Check if both errors are present in the aggregated error
		// This relies on multierr.Append behavior, which might not be strictly ordered
		errStr := err.Error()
		if !strings.Contains(errStr, expectedErr1.Error()) ||
			!strings.Contains(errStr, expectedErr2.Error()) {
			t.Errorf("Expected aggregated error to contain both errors, got %q", errStr)
		}
	})

	t.Run("Parallel with progress", func(t *testing.T) {
		var progress1, progress2 int64
		task1 := NewTask(
			WithName("Progress Task 1"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				for i := 0; i <= 5; i++ {
					atomic.StoreInt64(tracker.Current, int64(i))
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			}),
			WithProgress(&progress1),
		)
		task2 := NewTask(
			WithName("Progress Task 2"),
			WithRun(func(ctx context.Context, tracker *ProgressTracker) error {
				for i := 0; i <= 10; i++ {
					atomic.StoreInt64(tracker.Current, int64(i))
					time.Sleep(5 * time.Millisecond)
				}
				return nil
			}),
			WithProgress(&progress2),
		)

		err := Parallel(context.Background(), task1, task2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if atomic.LoadInt64(&progress1) != 5 {
			t.Errorf("Expected final progress1 to be 5, got %d", atomic.LoadInt64(&progress1))
		}
		if atomic.LoadInt64(&progress2) != 10 {
			t.Errorf("Expected final progress2 to be 10, got %d", atomic.LoadInt64(&progress2))
		}
	})
}

// Helper function to compare string slices (order-sensitive)
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Helper function to check if a string is present in a slice
func containsString(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func TestAutoUpdate(t *testing.T) {
	var progress int64 = 0
	var buf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())

	task := NewTask(
		WithName("AutoUpdate Task"),
		WithProgress(&progress),
		WithOutputWriter(&buf),
	)

	go autoUpdate(ctx, task)

	// Simulate progress updates
	atomic.StoreInt64(&progress, 5)
	time.Sleep(150 * time.Millisecond) // Give autoUpdate time to tick
	atomic.StoreInt64(&progress, 10)
	time.Sleep(150 * time.Millisecond)

	cancel()                          // Stop autoUpdate goroutine
	time.Sleep(50 * time.Millisecond) // Give it a moment to shut down

	output := buf.String()
	if !strings.Contains(output, "Task 'AutoUpdate Task': Progress 5") {
		t.Errorf("Expected output to contain 'Progress 5', got %q", output)
	}
	if !strings.Contains(output, "Task 'AutoUpdate Task': Progress 10") {
		t.Errorf("Expected output to contain 'Progress 10', got %q", output)
	}
}
