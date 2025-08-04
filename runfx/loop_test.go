package runfx

import (
	"context"
	"io"
	"testing"
	"time"
)

type dummyVisual struct{}

func (d dummyVisual) Render() []byte    { return nil }
func (d dummyVisual) OnResize(int, int) {}

// TestMountVisualLimit ensures that mounting more than MaxVisuals visuals
// returns ErrTooManyVisuals and does not exceed the allowed limit.
func TestMountVisualLimit(t *testing.T) {
	loop := Start()
	for i := 0; i < MaxVisuals; i++ {
		if _, err := loop.Mount(dummyVisual{}); err != nil {
			t.Fatalf("mount failed at %d: %v", i, err)
		}
	}

	if _, err := loop.Mount(dummyVisual{}); err != ErrTooManyVisuals {
		t.Fatalf("expected ErrTooManyVisuals, got %v", err)
	}
}

// TestMountNilVisual ensures mounting a nil visual returns ErrMountFailed
func TestMountNilVisual(t *testing.T) {
	loop := Start()
	if _, err := loop.Mount(nil); err != ErrMountFailed {
		t.Fatalf("expected ErrMountFailed, got %v", err)
	}
}

// TestRunNilContext verifies Run handles a nil context without panic
func TestRunNilContext(t *testing.T) {
	loop := StartWith(Config{Output: io.Discard, TickInterval: time.Millisecond})
	ml := loop.(*MainLoop)

	pr, pw := io.Pipe()
	ml.reader = NewKeyReader(pr)

	done := make(chan error, 1)
	go func() {
		done <- ml.Run(nil)
	}()

	time.Sleep(5 * time.Millisecond)
	if err := ml.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}
	pw.Close()
	if err := <-done; err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

// TestStopWithoutRun ensures Stop returns ErrLoopNotRunning when loop hasn't started
func TestStopWithoutRun(t *testing.T) {
	loop := Start()
	if err := loop.Stop(); err != ErrLoopNotRunning {
		t.Fatalf("expected ErrLoopNotRunning, got %v", err)
	}
}

// TestRunAlreadyRunning ensures Run cannot be called concurrently
func TestRunAlreadyRunning(t *testing.T) {
	loop := StartWith(Config{Output: io.Discard, TickInterval: time.Millisecond})
	ml := loop.(*MainLoop)

	pr, pw := io.Pipe()
	ml.reader = NewKeyReader(pr)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		ml.Run(ctx)
		close(done)
	}()

	// Give the loop time to start
	time.Sleep(5 * time.Millisecond)

	if err := ml.Run(context.Background()); err != ErrLoopAlreadyRunning {
		t.Fatalf("expected ErrLoopAlreadyRunning, got %v", err)
	}

	cancel()
	pw.Close()
	<-done
}
