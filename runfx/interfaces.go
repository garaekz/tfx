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

// Loop defines the runtime loop for mounting and managing visuals.
type Loop interface {
	Mount(v Visual) (unmount func(), err error)
	Run(ctx context.Context) error
	Stop() error
	IsRunning() bool
}
