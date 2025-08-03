package flowfx

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

// Task represents a single unit of work in a flow.
// Its creation uses functional options to avoid a giant struct.
// The task logic is defined in the Run closure.
type Task struct {
	Name         string
	Run          func(ctx context.Context, tracker *ProgressTracker) error
	progressPtr  *int64
	outputWriter io.Writer // Add this field
}

// TaskOption is a functional option for configuring a Task.
type TaskOption func(*Task)

// NewTask creates a new Task with the given options.
func NewTask(opts ...TaskOption) *Task {
	t := &Task{
		outputWriter: os.Stdout, // Default to os.Stdout
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// WithName sets the name of the task.
func WithName(name string) TaskOption {
	return func(t *Task) {
		t.Name = name
	}
}

// WithRun sets the execution logic of the task.
func WithRun(run func(ctx context.Context, tracker *ProgressTracker) error) TaskOption {
	return func(t *Task) {
		t.Run = run
	}
}

// WithProgress sets the progress pointer for the task.
func WithProgress(progress *int64) TaskOption {
	return func(t *Task) {
		t.progressPtr = progress
	}
}

// WithOutputWriter sets the output writer for the task.
func WithOutputWriter(writer io.Writer) TaskOption {
	return func(t *Task) {
		t.outputWriter = writer
	}
}

// ProgressTracker allows decoupled progress reporting from business logic.
// The task updates the progress via a pointer, and the flow orchestrator
// can read it to update the UI.
type ProgressTracker struct {
	Current *int64
	Total   int64
}

// autoUpdate is a hook that runs in a goroutine to update the UI.
func autoUpdate(ctx context.Context, task *Task) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if task.progressPtr != nil {
				fmt.Fprintf(task.outputWriter, "Task '%s': Progress %d\n", task.Name, *task.progressPtr)
			}
		}
	}
}

// Sequence runs a series of tasks sequentially.
// It stops at the first error.
func Sequence(ctx context.Context, tasks ...*Task) error {
	for _, task := range tasks {
		ctx, cancel := context.WithCancel(ctx)
		go autoUpdate(ctx, task)

		if err := task.Run(ctx, &ProgressTracker{Current: task.progressPtr}); err != nil {
			cancel()
			return err
		}
		cancel()
	}
	return nil
}

// Parallel runs a series of tasks in parallel.
// It waits for all tasks to complete and collects all errors.
func Parallel(ctx context.Context, tasks ...*Task) error {
	var g errgroup.Group
	var mu sync.Mutex
	var allErrors error

	for _, task := range tasks {
		t := task // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			ctx, cancel := context.WithCancel(ctx)
			go autoUpdate(ctx, t)

			err := t.Run(ctx, &ProgressTracker{Current: t.progressPtr})
			if err != nil {
				mu.Lock()
				allErrors = multierr.Append(allErrors, err)
				mu.Unlock()
			}
			cancel()
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return allErrors
	}

	return nil
}
