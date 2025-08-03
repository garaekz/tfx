package runfx

import (
	"context"
	"sync"
	"time"
)

// EventLoop manages high-performance animation ticks
type EventLoop struct {
	interval   time.Duration
	stopCh     chan struct{}
	isRunning  bool
	mu         sync.Mutex
	visuals    []Visual
	dirtyFlags []bool
}

func NewEventLoop(interval time.Duration) *EventLoop {
	return &EventLoop{
		interval:   interval,
		stopCh:     make(chan struct{}),
		visuals:    []Visual{},
		dirtyFlags: []bool{},
	}
}

func (l *EventLoop) Mount(v Visual) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.visuals = append(l.visuals, v)
	l.dirtyFlags = append(l.dirtyFlags, true)
}

func (l *EventLoop) Run(ctx context.Context) {
	l.mu.Lock()
	l.isRunning = true
	l.mu.Unlock()
	for {
		select {
		case <-l.stopCh:
			return
		case <-ctx.Done():
			return
		default:
			l.tick()
			time.Sleep(l.interval)
		}
	}
}

// tick calls Tick and Render only for dirty visuals
func (l *EventLoop) tick() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, v := range l.visuals {
		if l.dirtyFlags[i] {
			v.Tick(time.Now())
			// Assume visuals Render to global writer (or pass as arg)
			// v.Render(writer)
			l.dirtyFlags[i] = false
		}
	}
}

// MarkDirty marks a visual as dirty for next tick
func (l *EventLoop) MarkDirty(idx int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if idx >= 0 && idx < len(l.dirtyFlags) {
		l.dirtyFlags[idx] = true
	}
}

func (l *EventLoop) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.isRunning {
		close(l.stopCh)
		l.isRunning = false
	}
}
