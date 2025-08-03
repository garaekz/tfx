// Package progress provides beautiful progress bars and spinners with multipath API support.
//
// This package follows TFX multipath pattern (see MULTIPATH.md):
//   - BEGINNER: Start() / Start(config)    // Zero-config or config struct
//   - HARDCORE: New().Method().Method()    // DSL builder pattern
//   - EXPERIMENTAL: StartWith(opts...)     // Functional options (not in use)
//
// Progress bars are powered by RunFX for robust terminal management.
package progress

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
	"github.com/garaekz/tfx/terminal"
)

// Progress tracks and displays progress of operations and implements runfx.Visual.
type Progress struct {
	total     int
	current   int
	label     string
	width     int
	startTime time.Time
	theme     ProgressTheme
	style     ProgressStyle
	effect    ProgressEffect
	writer    io.Writer
	detector  *terminal.Detector
	mu        sync.Mutex
	ShowETA   bool

	// RunFX integration
	loop      runfx.Loop
	unmount   func()
	isStarted bool
}

// ProgressConfig provides structured configuration for Progress.
type ProgressConfig struct {
	Total   int
	Label   string
	Width   int
	Theme   ProgressTheme
	Style   ProgressStyle
	Effect  ProgressEffect
	Writer  io.Writer
	ShowETA bool
}

// --- INTERNAL IMPLEMENTATION ---

// newProgress is the internal implementation.
func newProgress(cfg ProgressConfig) *Progress {
	return &Progress{
		total:    cfg.Total,
		label:    cfg.Label,
		width:    cfg.Width,
		theme:    cfg.Theme,
		style:    cfg.Style,
		effect:   cfg.Effect,
		writer:   cfg.Writer,
		detector: terminal.NewDetector(cfg.Writer),
		ShowETA:  cfg.ShowETA,
	}
}

// --- RUNFX VISUAL INTERFACE ---

// Render implements runfx.Visual interface
func (p *Progress) Render(w share.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		return
	}

	rendered := p.renderProgress()
	entry := &share.Entry{
		Message: rendered,
	}
	w.Write(entry)
}

// Tick implements runfx.Visual interface
func (p *Progress) Tick(now time.Time) {
	// Progress bars typically don't need tick updates unless they have animations
	// This could be used for effects like pulsing or color cycling
}

// OnResize implements runfx.Visual interface
func (p *Progress) OnResize(cols, rows int) {
	// Adjust progress bar width based on terminal size if needed
	p.mu.Lock()
	defer p.mu.Unlock()

	// Could adjust width to fit terminal
	if cols > 0 && cols < p.width+20 { // Leave some margin
		p.width = max(cols-20, 10)
	}
}

// --- PROGRESS METHODS ---

// Start begins the progress tracking and mounts to RunFX
func (p *Progress) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isStarted {
		return
	}

	p.isStarted = true
	p.startTime = time.Now()

	// Create RunFX loop and mount this progress bar
	p.loop = runfx.Start()
	var err error
	p.unmount, err = p.loop.Mount(p)
	if err != nil {
		// Fallback to direct rendering if RunFX fails
		p.loop = nil
		p.renderDirect()
		return
	}

	// Start the RunFX loop in background
	go func() {
		ctx := context.Background()
		p.loop.Run(ctx)
	}()

	// Initial render
	p.renderDirect()
}

// Set updates the current progress value
func (p *Progress) Set(current int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		p.isStarted = true
		p.startTime = time.Now()
		// Don't auto-start in a goroutine to avoid race conditions
		// The caller should call Start() explicitly
	}

	// Clamp current to not exceed total
	p.current = min(current, p.total)

	// Always render to show the update
	p.renderDirect()
}

// SetLabel updates the progress label
func (p *Progress) SetLabel(label string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.label = label

	// Render if already started
	if p.isStarted {
		p.renderDirect()
	}
}

// SetTotal updates the total value
func (p *Progress) SetTotal(total int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.total = total
}

// Finish completes the progress and unmounts from RunFX
func (p *Progress) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		return
	}

	p.current = p.total

	// Render final state
	if p.loop == nil {
		p.renderDirect()
	}

	// Unmount from RunFX
	if p.unmount != nil {
		p.unmount()
		p.unmount = nil
	}

	// Stop the loop
	if p.loop != nil {
		p.loop.Stop()
		p.loop = nil
	}
}

// Close cleans up resources
func (p *Progress) Close() {
	if p.unmount != nil {
		p.unmount()
	}
	if p.loop != nil {
		p.loop.Stop()
	}
}

// --- RENDERING ---

// renderProgress generates the progress bar string
func (p *Progress) renderProgress() string {
	percentage := float64(p.current) / float64(p.total)
	if percentage > 1.0 {
		percentage = 1.0
	}

	// Use the existing render logic from render.go
	return RenderBar(p, p.detector)
}

// renderDirect renders directly to writer (fallback when RunFX not available)
func (p *Progress) renderDirect() {
	rendered := p.renderProgress()
	fmt.Fprint(p.writer, "\r"+rendered)
}

// --- ADDITIONAL METHODS FOR COMPATIBILITY ---

// Increment increases the progress by 1
func (p *Progress) Increment() {
	p.Add(1)
}

// Add increases the progress by the specified amount
func (p *Progress) Add(amount int) {
	p.Set(p.current + amount)
}

// Complete sets progress to 100% and finishes
func (p *Progress) Complete(message ...string) {
	p.Set(p.total)

	// Show completion message with icon
	if len(message) > 0 && len(message[0]) > 0 {
		completion := RenderCompletion(p.theme, message[0], true, p.detector)
		if p.loop == nil {
			fmt.Fprint(p.writer, completion)
		}
	} else {
		// Default completion message
		completion := RenderCompletion(p.theme, "Complete", true, p.detector)
		if p.loop == nil {
			fmt.Fprint(p.writer, completion)
		}
	}

	p.Finish()
}

// Fail marks progress as failed and finishes
func (p *Progress) Fail(message ...string) {
	p.Set(p.total) // Set to 100% to show completion

	// Show failure message with icon
	if len(message) > 0 && len(message[0]) > 0 {
		failure := RenderCompletion(p.theme, message[0], false, p.detector)
		if p.loop == nil {
			fmt.Fprint(p.writer, failure)
		}
	} else {
		// Default failure message
		failure := RenderCompletion(p.theme, "Failed", false, p.detector)
		if p.loop == nil {
			fmt.Fprint(p.writer, failure)
		}
	}

	p.Finish()
}

// SetTheme updates the progress theme
func (p *Progress) SetTheme(theme ProgressTheme) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.theme = theme
}

// SetEffect updates the progress effect
func (p *Progress) SetEffect(effect ProgressEffect) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.effect = effect
	// Enable or disable effects based on the effect type
	if effect == EffectNone {
		p.theme.EffectEnabled = false
	} else {
		p.theme.EffectEnabled = true
	}
}

// GetPercent returns the current percentage (0.0 to 100.0)
func (p *Progress) GetPercent() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.total == 0 {
		return 0.0
	}
	percent := float64(p.current) / float64(p.total) * 100.0
	if percent > 100.0 {
		return 100.0
	}
	return percent
}

// GetElapsed returns the elapsed time since progress started
func (p *Progress) GetElapsed() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.isStarted {
		return 0
	}
	return time.Since(p.startTime)
}

// GetETA estimates the time remaining based on current rate
func (p *Progress) GetETA() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted || p.current == 0 {
		return 0
	}

	elapsed := time.Since(p.startTime)
	rate := float64(p.current) / elapsed.Seconds()

	if rate > 0 {
		remaining := float64(p.total-p.current) / rate
		return time.Duration(remaining * float64(time.Second))
	}

	return 0
}

// Redraw forces a redraw of the progress bar
func (p *Progress) Redraw() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		return // Don't render if not started
	}

	// Always render directly for immediate feedback
	p.renderDirect()
}
