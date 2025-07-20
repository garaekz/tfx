package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/shared"
	"github.com/garaekz/tfx/terminal"
)

// Spinner renders an animated spinner with an optional message.
type Spinner struct {
	frames   []string
	interval time.Duration
	message  string
	theme    ProgressTheme
	effect   SpinnerEffect
	writer   io.Writer
	detector *terminal.Detector
	mu       sync.Mutex
	active   bool
	stopCh   chan struct{}
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

// DefaultSpinnerConfig returns default configuration for Spinner
func DefaultSpinnerConfig() SpinnerConfig {
	return SpinnerConfig{
		Message:  "",
		Frames:   []string{"|", "/", "-", "\\"},
		Interval: 100 * time.Millisecond,
		Theme:    MaterialTheme,
		Effect:   SpinnerEffectNone,
		Writer:   os.Stdout,
	}
}

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick default
func StartSpinner(message string) *Spinner {
	return NewSpinner(DefaultSpinnerConfig(), WithMessage(message))
}

// 2. INSTANTIATED: Config struct
func NewSpinner(cfg SpinnerConfig, opts ...shared.Option[SpinnerConfig]) *Spinner {
	// Apply functional options to config
	shared.ApplyOptions(&cfg, opts...)

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

	return s
}

// 3. FLUENT: Functional options
func NewSpinnerWith(opts ...shared.Option[SpinnerConfig]) *Spinner {
	cfg := DefaultSpinnerConfig()
	return NewSpinner(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

// WithMessage sets the spinner message
func WithMessage(message string) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Message = message
	}
}

// WithSpinnerFrames sets custom frames for the spinner animation
func WithSpinnerFrames(frames []string) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Frames = frames
	}
}

// WithSpinnerInterval sets the frame interval duration
func WithSpinnerInterval(interval time.Duration) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Interval = interval
	}
}

// WithSpinnerTheme applies a ProgressTheme to colorize the spinner
func WithSpinnerTheme(theme ProgressTheme) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Theme = theme
	}
}

// WithSpinnerEffect applies visual effects
func WithSpinnerEffect(effect SpinnerEffect) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Effect = effect
	}
}

// WithSpinnerWriter sets a custom writer for spinner output
func WithSpinnerWriter(writer io.Writer) shared.Option[SpinnerConfig] {
	return func(cfg *SpinnerConfig) {
		cfg.Writer = writer
	}
}

// --- CONVENIENCE OPTIONS ---

// WithSpinnerMaterialTheme applies Material Design theme
func WithSpinnerMaterialTheme() shared.Option[SpinnerConfig] {
	return WithSpinnerTheme(MaterialTheme)
}

// WithSpinnerDraculaTheme applies Dracula theme
func WithSpinnerDraculaTheme() shared.Option[SpinnerConfig] {
	return WithSpinnerTheme(DraculaTheme)
}

// WithSpinnerNordTheme applies Nord theme
func WithSpinnerNordTheme() shared.Option[SpinnerConfig] {
	return WithSpinnerTheme(NordTheme)
}

// WithSpinnerRainbow enables rainbow effect
func WithSpinnerRainbow() shared.Option[SpinnerConfig] {
	return WithSpinnerEffect(SpinnerEffectRainbow)
}

// --- PRESET FRAMES ---

// WithDotsFrames uses dots animation
func WithDotsFrames() shared.Option[SpinnerConfig] {
	return WithSpinnerFrames([]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"})
}

// WithArrowFrames uses arrow animation
func WithArrowFrames() shared.Option[SpinnerConfig] {
	return WithSpinnerFrames([]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"})
}

// WithBounceFrames uses bouncing animation
func WithBounceFrames() shared.Option[SpinnerConfig] {
	return WithSpinnerFrames([]string{"⠁", "⠂", "⠄", "⠂"})
}

// --- SPINNER METHODS ---

// Start begins the spinner in a new goroutine
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go s.spin()
}

// Stop stops the spinner and prints a success message
func (s *Spinner) Stop(msg string) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.stopCh <- struct{}{}
	s.mu.Unlock()

	// Clear line and print completion message
	completion := RenderCompletion(s.theme, msg, true, s.detector)
	fmt.Fprint(s.writer, "\r"+completion+"\n")
}

// Fail stops the spinner and prints a failure message
func (s *Spinner) Fail(msg string) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.stopCh <- struct{}{}
	s.mu.Unlock()

	// Clear line and print failure message
	completion := RenderCompletion(s.theme, msg, false, s.detector)
	fmt.Fprint(s.writer, "\r"+completion+"\n")
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

var globalSpinner *Spinner
var globalSpinnerMu sync.Mutex

// StartGlobalSpinner starts a global spinner (EXPRESS API)
func StartGlobalSpinner(message string) {
	globalSpinnerMu.Lock()
	if globalSpinner != nil && globalSpinner.active {
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
	return NewSpinnerWith(
		WithMessage(message),
		WithSpinnerMaterialTheme(),
	)
}

// NewDraculaSpinner creates a spinner with Dracula theme and rainbow effect
func NewDraculaSpinner(message string) *Spinner {
	return NewSpinnerWith(
		WithMessage(message),
		WithSpinnerDraculaTheme(),
		WithSpinnerRainbow(),
	)
}

// NewNordSpinner creates a spinner with Nord theme
func NewNordSpinner(message string) *Spinner {
	return NewSpinnerWith(
		WithMessage(message),
		WithSpinnerNordTheme(),
	)
}

// NewDotsSpinner creates a spinner with dots animation
func NewDotsSpinner(message string) *Spinner {
	return NewSpinnerWith(
		WithMessage(message),
		WithDotsFrames(),
	)
}