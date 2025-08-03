package runfx

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/terminal"
	"github.com/garaekz/tfx/writer"
)

// MainLoop is the central orchestrator that implements the Loop interface
type MainLoop struct {
	mux       *Multiplexer
	terminal  *writer.TerminalWriter
	cursor    *CursorManager
	screen    *ScreenManager
	signals   terminal.SignalHandler
	eventLoop *EventLoop
	ttyInfo   *TTYInfo
	isRunning int32 // atomic
	isStopped int32 // atomic guard for stopCh
	stopCh    chan struct{}

	// Context state protected by ctxMu
	ctxMu  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	mu         sync.RWMutex
	output     io.Writer
	nextRegion int
	testMode   bool // Set to true in tests to force ticking even in fallback mode
}

// VisualWriter implements share.Writer for capturing visual output
type VisualWriter struct {
	buffer []byte
}

func (vw *VisualWriter) Write(entry *share.Entry) error {
	// Simple implementation - just capture the formatted text
	text := fmt.Sprintf("%s\n", entry.Message)
	vw.buffer = append(vw.buffer, []byte(text)...)
	return nil
}

func (vw *VisualWriter) Close() error {
	return nil
}

func (vw *VisualWriter) Bytes() []byte {
	return vw.buffer
}

func (vw *VisualWriter) Reset() {
	vw.buffer = nil
}

// NewMainLoop creates a new main loop with all integrated components
func NewMainLoop(tickInterval time.Duration) *MainLoop {
	ttyInfo := DetectTTY()

	// Create terminal writer with double-buffering enabled
	terminalOpts := writer.TerminalOptions{
		DoubleBuffer: true,
		ForceColor:   false,
		DisableColor: ttyInfo.NoColor,
	}
	termWriter := writer.NewTerminalWriter(os.Stdout, terminalOpts)

	return &MainLoop{
		mux:        NewMultiplexer(),
		terminal:   termWriter,
		cursor:     &CursorManager{},
		screen:     NewScreenManager(),
		signals:    *terminal.NewSignalHandler(),
		eventLoop:  NewEventLoop(tickInterval),
		ttyInfo:    &ttyInfo,
		stopCh:     make(chan struct{}),
		output:     os.Stdout,
		nextRegion: 0,
	}
}

// NewMainLoopWithTerminal creates a main loop with a custom terminal writer for testing
func NewMainLoopWithTerminal(tickInterval time.Duration, termWriter *writer.TerminalWriter, testMode bool) *MainLoop {
	ttyInfo := DetectTTY()

	return &MainLoop{
		mux:        NewMultiplexer(),
		terminal:   termWriter,
		cursor:     &CursorManager{},
		screen:     NewScreenManager(),
		signals:    *terminal.NewSignalHandler(),
		eventLoop:  NewEventLoop(tickInterval),
		ttyInfo:    &ttyInfo,
		stopCh:     make(chan struct{}),
		output:     termWriter,
		nextRegion: 0,
		testMode:   testMode,
	}
}

// Mount implements the Loop interface
func (ml *MainLoop) Mount(v Visual) (unmount func(), err error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	if !ml.ttyInfo.IsTTY && !ml.testMode {
		LogMount("fallback")
		FallbackOutput("Visual mounted (fallback mode)")
		return func() {}, nil
	}

	// Generate unique name for the visual
	visualName := fmt.Sprintf("visual_%d", ml.nextRegion)
	ml.nextRegion++

	// Mount in multiplexer
	ml.mux.Mount(visualName, v)

	// Allocate screen region (assuming each visual gets 3 lines for now)
	regionStart := ml.nextRegion * 3
	regionEnd := regionStart + 2
	ml.screen.AllocateRegion(visualName, regionStart, regionEnd)

	// Add to event loop
	ml.eventLoop.Mount(v)

	LogMount(visualName)

	// Return unmount function
	unmount = func() {
		ml.mu.Lock()
		defer ml.mu.Unlock()

		// Clear the region
		ml.screen.ClearRegion(visualName)

		// Remove from multiplexer
		ml.mux.Unmount(visualName)

		LogUnmount(visualName)
	}

	return unmount, nil
}

// Run implements the Loop interface
func (ml *MainLoop) Run(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&ml.isRunning, 0, 1) {
		return ErrLoopAlreadyRunning
	}
	defer atomic.StoreInt32(&ml.isRunning, 0)

	ml.ctxMu.Lock()
	ml.ctx, ml.cancel = context.WithCancel(ctx)
	cancel := ml.cancel
	ml.ctxMu.Unlock()
	defer cancel()

	if !ml.ttyInfo.IsTTY && !ml.testMode {
		DebugLog("Running in fallback mode (non-TTY)")
		<-ml.ctx.Done()
		return nil
	}

	// Initialize terminal state (only if we have a real TTY)
	if ml.ttyInfo.IsTTY {
		ml.terminal.HideCursor()
		defer ml.terminal.ShowCursor()
	}

	// Start signal handler
	go ml.handleSignals()

	// Start event loop
	go ml.eventLoop.Run(ml.ctx)

	// Main render loop
	ticker := time.NewTicker(16 * time.Millisecond) // ~60fps
	defer ticker.Stop()

	for {
		select {
		case <-ml.ctx.Done():
			return ml.ctx.Err()
		case <-ml.stopCh:
			return nil
		case <-ticker.C:
			ml.renderFrame()
		}
	}
}

// Stop implements the Loop interface
func (ml *MainLoop) Stop() error {
	if atomic.LoadInt32(&ml.isRunning) == 0 {
		return ErrLoopNotRunning
	}

	// Use atomic operation to ensure channel is only closed once
	if atomic.CompareAndSwapInt32(&ml.isStopped, 0, 1) {
		close(ml.stopCh)
	}

	ml.ctxMu.Lock()
	cancel := ml.cancel
	ml.ctxMu.Unlock()

	if cancel != nil {
		cancel()
	}

	ml.eventLoop.Stop()
	ml.signals.Stop()
	return nil
}

// IsRunning implements the Loop interface
func (ml *MainLoop) IsRunning() bool {
	return atomic.LoadInt32(&ml.isRunning) != 0
}

// renderFrame renders all mounted visuals
func (ml *MainLoop) renderFrame() {
	ml.mu.RLock()
	visuals := ml.mux.ListVisuals()
	ml.mu.RUnlock()

	if len(visuals) == 0 {
		return
	}

	// Build frame buffer
	var frameBuffer []byte

	for _, visualName := range visuals {
		visual, exists := ml.mux.GetVisual(visualName)
		if !exists {
			continue
		}

		// Create a buffer writer for this visual
		bufWriter := &VisualWriter{}
		visual.Render(bufWriter)

		// Position cursor for this visual's region
		region, hasRegion := ml.screen.GetRegion(visualName)
		if hasRegion {
			// Move cursor to region start
			cursorPos := fmt.Sprintf("\033[%d;1H", region[0]+1)
			frameBuffer = append(frameBuffer, cursorPos...)
		}

		frameBuffer = append(frameBuffer, bufWriter.Bytes()...)
		LogRender(visualName)
	}

	// Render the complete frame using terminal writer
	ml.terminal.Write(frameBuffer)
}

// handleSignals manages resize and termination signals
func (ml *MainLoop) handleSignals() {
	// Set up the callbacks
	ml.signals.OnResize(ml.handleResize)
	ml.signals.OnStop(ml.handleStop)

	// Start listening
	ml.signals.Listen(ml.ctx)
}

// handleResize handles terminal resize events
func (ml *MainLoop) handleResize() {
	DebugLog("Terminal resize detected")

	// Get new terminal size
	cols, rows, err := terminal.GetSize()
	if err != nil {
		DebugLog("Failed to get terminal size: %v", err)
		return
	}

	// Notify all visuals of resize
	ml.mu.RLock()
	visuals := ml.mux.ListVisuals()
	ml.mu.RUnlock()

	for _, visualName := range visuals {
		visual, exists := ml.mux.GetVisual(visualName)
		if exists {
			visual.OnResize(cols, rows)
		}
	}

	// Clear and reallocate regions if needed
	ml.reallocateRegions(rows)
}

// handleStop handles termination signals
func (ml *MainLoop) handleStop() {
	DebugLog("Stop signal received")
	ml.Stop()
}

// reallocateRegions reallocates screen regions based on new terminal size
func (ml *MainLoop) reallocateRegions(rows int) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	visuals := ml.mux.ListVisuals()
	if len(visuals) == 0 {
		return
	}

	// Simple reallocation: divide available rows evenly
	linesPerVisual := rows / len(visuals)
	if linesPerVisual < 1 {
		linesPerVisual = 1
	}

	for i, visualName := range visuals {
		start := i * linesPerVisual
		end := start + linesPerVisual - 1
		if end >= rows {
			end = rows - 1
		}
		ml.screen.Reallocate(visualName, start, end)
	}
}

// BufferWriter implements share.Writer for capturing visual output
type BufferWriter struct {
	buffer []byte
}

func (bw *BufferWriter) Write(p []byte) (n int, err error) {
	bw.buffer = append(bw.buffer, p...)
	return len(p), nil
}

func (bw *BufferWriter) Bytes() []byte {
	return bw.buffer
}
