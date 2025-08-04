package runfx

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/term"

	"github.com/garaekz/tfx/terminal"
	"github.com/garaekz/tfx/writer"
)

// --- Event Types ---
type (
	keyEvent    Key
	tickEvent   struct{ time time.Time }
	resizeEvent struct{ cols, rows int }
	errorEvent  error
)

// --- Loop Definition ---

// MainLoop is handles terminal I/O, signals, and the render/tick cycle.
type MainLoop struct {
	// Configuration and Dependencies
	writer  *writer.TerminalWriter
	reader  *KeyReader
	signals *terminal.SignalHandler
	mux     *Multiplexer
	ticker  *time.Ticker

	// Internal State
	events   chan any // Central event channel
	cancel   context.CancelFunc
	rawState *term.State
	running  atomic.Bool
}

// --- Public API Methods ---

// Mount registers a visual component with the loop's multiplexer.
func (ml *MainLoop) Mount(v Visual) (unmount func(), err error) {
	if v == nil {
		return nil, ErrMountFailed
	}

	// Enforce a maximum number of mounted visuals to prevent stack
	// exhaustion when rendering large numbers of components. The check is
	// done before mounting, and verified after to handle concurrent mounts.
	if ml.mux.Count() >= MaxVisuals {
		return nil, ErrTooManyVisuals
	}

	// Get the current terminal size and inform the new visual immediately.
	// This ensures the component has its layout calculated before the first render.
	if cols, rows, err := ml.writer.GetSize(); err == nil {
		v.OnResize(cols, rows)
	}

	// Mount the visual in the multiplexer and get its unique ID.
	id := ml.mux.Mount(v)

	if ml.mux.Count() > MaxVisuals {
		ml.mux.Unmount(id)
		return nil, ErrTooManyVisuals
	}

	// Return a closure that captures the ID to unmount the visual later.
	return func() { ml.mux.Unmount(id) }, nil
}

// Run starts the main loop and blocks until the context is canceled or Stop() is called.
func (ml *MainLoop) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if !ml.running.CompareAndSwap(false, true) {
		return ErrLoopAlreadyRunning
	}
	defer ml.running.Store(false)

	// Setup terminal
	if state, err := ml.writer.EnableRawMode(); err == nil {
		ml.rawState = state
		defer ml.writer.RestoreMode(ml.rawState)
	}
	ml.writer.HideCursor()
	defer ml.writer.ShowCursor()
	defer ml.writer.Clear()

	// Create a cancellable context for the loop's goroutines
	loopCtx, cancel := context.WithCancel(ctx)
	ml.cancel = cancel
	defer ml.ticker.Stop()

	// Start event producers
	go ml.produceKeyEvents(loopCtx)
	go ml.produceTickEvents(loopCtx)
	go ml.produceSignalEvents(loopCtx)

	// Initial render
	ml.renderFrame()

	// Main event processing loop
	for {
		select {
		case <-loopCtx.Done():
			return loopCtx.Err()
		case e := <-ml.events:
			shouldStop, shouldRender := ml.handleEvent(e)
			if shouldStop {
				return nil
			}
			if shouldRender {
				ml.renderFrame()
			}
		}
	}
}

// Stop gracefully shuts down the main loop.
func (ml *MainLoop) Stop() error {
	if !ml.running.Load() {
		return ErrLoopNotRunning
	}
	if ml.cancel != nil {
		ml.cancel()
		return nil
	}
	return ErrLoopClosed
}

// IsRunning checks if the loop is currently active.
func (ml *MainLoop) IsRunning() bool {
	return ml.running.Load()
}

// --- Internal Event Producers ---

func (ml *MainLoop) produceKeyEvents(ctx context.Context) {
	for {
		key, err := ml.reader.ReadKey(ctx)
		if err != nil {
			select {
			case ml.events <- errorEvent(err):
			case <-ctx.Done():
				return
			}
			return
		}
		select {
		case ml.events <- keyEvent(key):
		case <-ctx.Done():
			return
		}
	}
}

func (ml *MainLoop) produceTickEvents(ctx context.Context) {
	for {
		select {
		case t := <-ml.ticker.C:
			select {
			case ml.events <- tickEvent{time: t}:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (ml *MainLoop) produceSignalEvents(ctx context.Context) {
	ml.signals.OnResize(func() {
		cols, rows, err := ml.writer.GetSize()
		if err == nil {
			select {
			case ml.events <- resizeEvent{cols: cols, rows: rows}:
			case <-ctx.Done():
			}
		}
	})
	ml.signals.OnStop(func() {
		ml.Stop()
	})
	ml.signals.Listen(ctx)
}

// --- Internal Event Handler and Renderer ---

// handleEvent processes a single event from the central channel.
// It returns (shouldStop, shouldRender).
func (ml *MainLoop) handleEvent(e any) (bool, bool) {
	switch event := e.(type) {
	case keyEvent:
		// Dispatch key to all interactive visuals using their IDs.
		for _, id := range ml.mux.ListVisuals() {
			if v, ok := ml.mux.GetVisual(id); ok {
				if i, isInteractive := v.(Interactive); isInteractive {
					if i.OnKey(Key(event)) {
						// Stop if OnKey returns true. Render one last time.
						return true, true
					}
				}
			}
		}
		// If no component stopped the loop, we assume a state change and re-render.
		return false, true
	case tickEvent:
		// Dispatch tick to all updatable visuals.
		for _, id := range ml.mux.ListVisuals() {
			if v, ok := ml.mux.GetVisual(id); ok {
				if u, isUpdatable := v.(Updatable); isUpdatable {
					u.Tick(event.time)
				}
			}
		}
		// A tick always implies a potential visual change.
		return false, true
	case resizeEvent:
		// Dispatch resize to all visuals.
		ml.mux.OnResize(event.cols, event.rows)
		// A resize always requires a full re-render.
		return false, true
	case errorEvent:
		// Log or handle error, for now we stop.
		fmt.Fprintf(os.Stderr, "runfx error: %v\n", event)
		return true, false // Stop, no need to render.
	}
	// Default case for unknown events.
	return false, false
}

// renderFrame clears the screen and renders all mounted visuals.
func (ml *MainLoop) renderFrame() {
	ml.writer.Clear()
	// The multiplexer now handles aggregating the render output
	renderBytes := ml.mux.Render()
	ml.writer.Write(renderBytes)
	ml.writer.Flush()
}
