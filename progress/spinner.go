package progress

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
	"github.com/garaekz/tfx/terminal"
)

// Spinner renders an animated spinner with an optional message and implements runfx.Visual.
type Spinner struct {
	frames   []string
	interval time.Duration
	message  string
	theme    ProgressTheme
	effect   SpinnerEffect
	writer   io.Writer
	detector *terminal.Detector
	mu       sync.Mutex

	// RunFX integration
	loop       runfx.Loop
	unmount    func()
	isActive   bool
	frameIndex int
	startTime  time.Time

	// Legacy fields for compatibility
	stopCh chan struct{}
	active bool // alias for isActive
}

// SpinnerConfig provides structured configuration for Spinner
type SpinnerConfig struct {
	Message  string
	Frames   []string
	Interval time.Duration
	Theme    ProgressTheme
	Effect   SpinnerEffect
	Writer   io.Writer
}

// SpinnerEffect defines visual effects for spinners
type SpinnerEffect int

const (
	SpinnerEffectNone SpinnerEffect = iota
	SpinnerEffectRainbow
	SpinnerEffectPulse
	SpinnerEffectGlow
)

// --- INTERNAL IMPLEMENTATION ---

// 2. INSTANTIATED: Config struct
func NewSpinnerWithConfig(cfg SpinnerConfig) *Spinner {
	// Ensure defaults are set
	if len(cfg.Frames) == 0 {
		cfg.Frames = []string{"|", "/", "-", "\\"}
	}
	if cfg.Interval == 0 {
		cfg.Interval = 100 * time.Millisecond
	}
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	s := &Spinner{
		frames:   cfg.Frames,
		interval: cfg.Interval,
		message:  cfg.Message,
		theme:    cfg.Theme,
		effect:   cfg.Effect,
		writer:   cfg.Writer,
		stopCh:   make(chan struct{}, 1),
	}

	// Create terminal detector
	s.detector = terminal.NewDetector(s.writer)

	// Setup resize handling with signal handler
	signalHandler := terminal.NewSignalHandler()
	signalHandler.OnResize(s.Redraw)
	// Note: signalHandler.Start() would need to be called to actually listen

	return s
}

// --- SPINNER METHODS ---

// Start begins the spinner in a new goroutine
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.isActive {
		s.mu.Unlock()
		return
	}
	s.isActive = true
	s.startTime = time.Now()
	s.mu.Unlock()

	// Create RunFX loop with custom writer and mount this spinner
	cfg := runfx.Config{
		Output: s.writer,
	}
	s.loop = runfx.Start(cfg)
	var err error
	s.unmount, err = s.loop.Mount(s)
	if err != nil {
		// Fallback to direct rendering if RunFX fails
		go s.spin()
		return
	}

	// Start the RunFX loop in background
	go func() {
		ctx := context.Background()
		s.loop.Run(ctx)
	}()
}

// Stop stops the spinner and prints a success message
func (s *Spinner) Stop(msg string) {
	s.mu.Lock()
	if !s.isActive {
		s.mu.Unlock()
		return
	}
	s.isActive = false
	s.active = false // Keep legacy field updated
	if s.stopCh != nil {
		select {
		case s.stopCh <- struct{}{}:
		default:
		}
	}
	s.mu.Unlock()

	// Unmount from RunFX if mounted
	if s.unmount != nil {
		s.unmount()
	}

	// Clear line completely and print completion message
	fmt.Fprint(s.writer, "\r\033[K") // Clear entire line
	completion := RenderCompletion(s.theme, msg, true, s.detector)
	fmt.Fprint(s.writer, completion+"\n")
}

// Fail stops the spinner and prints a failure message
func (s *Spinner) Fail(msg string) {
	s.mu.Lock()
	if !s.isActive {
		s.mu.Unlock()
		return
	}
	s.isActive = false
	s.active = false // Keep legacy field updated
	if s.stopCh != nil {
		select {
		case s.stopCh <- struct{}{}:
		default:
		}
	}
	s.mu.Unlock()

	// Unmount from RunFX if mounted
	if s.unmount != nil {
		s.unmount()
	}

	// Clear line completely and print failure message
	fmt.Fprint(s.writer, "\r\033[K") // Clear entire line
	completion := RenderCompletion(s.theme, msg, false, s.detector)
	fmt.Fprint(s.writer, completion+"\n")
}

// Success is an alias for Stop for better readability
func (s *Spinner) Success(msg string) {
	s.Stop(msg)
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

// SetTheme updates the spinner theme
func (s *Spinner) SetTheme(theme ProgressTheme) {
	s.mu.Lock()
	s.theme = theme
	s.mu.Unlock()
}

// Redraw redraws the spinner.
func (s *Spinner) Redraw() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isActive {
		// This is a simplified redraw. A more robust implementation would
		// clear the line and re-render the current frame.
		fmt.Fprint(s.writer, "\r")
	}
}

// Close stops the spinner and unmounts from RunFX
func (s *Spinner) Close() {
	if s.unmount != nil {
		s.unmount()
	}
}

// --- runfx.Visual INTERFACE IMPLEMENTATION ---

// Render displays the current spinner frame to the provided writer
func (s *Spinner) Render(writer share.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isActive {
		return
	}

	// Safety check for empty frames
	if len(s.frames) == 0 {
		return
	}

	frame := s.frames[s.frameIndex%len(s.frames)]

	// Apply effects if enabled
	if s.effect == SpinnerEffectRainbow {
		frame = s.applyRainbowEffect(frame, s.frameIndex)
	}

	// Render spinner with theme
	text := s.renderFrame(frame)
	entry := &share.Entry{
		Message: text,
	}
	writer.Write(entry)
}

// Tick updates the spinner frame index based on the interval
func (s *Spinner) Tick(timestamp time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Increment frame index
	s.frameIndex = (s.frameIndex + 1) % len(s.frames)
}

// OnResize handles terminal resize events
func (s *Spinner) OnResize(width, height int) {
	// For spinner, we don't need special resize handling
	// The spinner is typically just one line and self-contained
}

// --- INTERNAL METHODS ---

// spin runs the spinner animation loop
func (s *Spinner) spin() {
	i := 0
	for {
		select {
		case <-s.stopCh:
			return
		default:
			s.mu.Lock()
			if !s.isActive {
				s.mu.Unlock()
				return
			}
			// Safety check for empty frames
			if len(s.frames) == 0 {
				s.mu.Unlock()
				return
			}
			frame := s.frames[i%len(s.frames)]

			// Apply effects if enabled
			if s.effect == SpinnerEffectRainbow {
				frame = s.applyRainbowEffect(frame, i)
			}

			// Render spinner with theme
			text := s.renderFrame(frame)
			fmt.Fprint(s.writer, text)

			i++
			s.mu.Unlock()
			time.Sleep(s.interval)
		}
	}
}

// renderFrame renders a single spinner frame with theme colors
func (s *Spinner) renderFrame(frame string) string {
	frameColor := s.theme.RenderColor(s.theme.CompleteColor, s.detector)
	labelColor := s.theme.RenderColor(s.theme.LabelColor, s.detector)

	styledFrame := frameColor + frame + color.Reset
	styledMessage := labelColor + s.message + color.Reset

	return fmt.Sprintf("\r%s %s", styledFrame, styledMessage)
}

// applyRainbowEffect applies rainbow colors to the spinner frame
func (s *Spinner) applyRainbowEffect(frame string, iteration int) string {
	colors := []color.Color{
		color.MaterialRed,
		color.MaterialOrange,
		color.MaterialYellow,
		color.MaterialGreen,
		color.MaterialBlue,
		color.MaterialPurple,
	}

	colorIndex := iteration % len(colors)
	coloredFrame := s.theme.RenderColor(colors[colorIndex], s.detector)
	return coloredFrame + frame + color.Reset
}

// --- GLOBAL CONVENIENCE FUNCTIONS ---

var (
	globalSpinner   *Spinner
	globalSpinnerMu sync.Mutex
)

// StartGlobalSpinner starts a global spinner (EXPRESS API)
func StartGlobalSpinner(message string) {
	globalSpinnerMu.Lock()
	if globalSpinner != nil && globalSpinner.isActive {
		globalSpinnerMu.Unlock()
		return
	}
	globalSpinner = StartSpinner(message)
	globalSpinnerMu.Unlock()
	globalSpinner.Start()
}

// StopGlobalSpinner stops the global spinner with success message
func StopGlobalSpinner(msg string) {
	globalSpinnerMu.Lock()
	if globalSpinner != nil {
		globalSpinner.Stop(msg)
	}
	globalSpinnerMu.Unlock()
}

// FailGlobalSpinner stops the global spinner with failure message
func FailGlobalSpinner(msg string) {
	globalSpinnerMu.Lock()
	if globalSpinner != nil {
		globalSpinner.Fail(msg)
	}
	globalSpinnerMu.Unlock()
}

// --- PRESET CONSTRUCTORS ---

// NewMaterialSpinner creates a spinner with Material Design theme
func NewMaterialSpinner(message string) *Spinner {
	return NewSpinnerWithConfig(SpinnerConfig{
		Message: message,
		Theme:   MaterialTheme,
	})
}

// NewDraculaSpinner creates a spinner with Dracula theme and rainbow effect
func NewDraculaSpinner(message string) *Spinner {
	return NewSpinnerWithConfig(SpinnerConfig{
		Message: message,
		Theme:   DraculaTheme,
		Effect:  SpinnerEffectRainbow,
	})
}

// NewNordSpinner creates a spinner with Nord theme
func NewNordSpinner(message string) *Spinner {
	return NewSpinnerWithConfig(SpinnerConfig{
		Message: message,
		Theme:   NordTheme,
	})
}

// NewDotsSpinner creates a spinner with dots animation
func NewDotsSpinner(message string) *Spinner {
	return NewSpinnerWithConfig(SpinnerConfig{
		Message: message,
		Frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	})
}
