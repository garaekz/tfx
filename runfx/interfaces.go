package runfx

import (
	"context"
	"time"

	"github.com/garaekz/tfx/internal/share"
)

// Visual defines a multiplexed terminal visual managed by RunFX.
type Visual interface {
	Render(w share.Writer)
	Tick(now time.Time)
	OnResize(cols, rows int)
}

// Interactive extends Visual with keyboard input handling.
// Use this interface for visuals that need to respond to user input.
type Interactive interface {
	Visual
	OnKey(key Key) bool // Returns true if the key was handled
}

// Loop defines the runtime loop for mounting and managing visuals.
type Loop interface {
	Mount(v Visual) (unmount func(), err error)
	Run(ctx context.Context) error
	Stop() error
	IsRunning() bool
}

// InteractiveLoop extends Loop with keyboard input support.
type InteractiveLoop interface {
	Loop
	MountInteractive(v Interactive) (unmount func(), err error)
}
