package runfx

import (
	"context"
	"time"

	"github.com/garaekz/tfx/writer"
)

// Visual represents a renderable component that can react to terminal events.
//
// Render writes the visual representation to the provided writer.
// Tick allows the visual to update any internal state on each loop cycle.
// OnResize notifies the visual of terminal size changes.
type Visual interface {
	Render(w writer.Writer)
	Tick(now time.Time)
	OnResize(cols, rows int)
}

// Interactive is a Visual that can respond to keyboard input.
type Interactive interface {
	Visual
	OnKey(key Key) bool // Returns true to stop the loop.
}

// Loop defines the runtime loop for mounting and managing visuals.
type Loop interface {
	Mount(v Visual) (unmount func(), err error)
	Run(ctx context.Context) error
	Stop() error
	IsRunning() bool
}
