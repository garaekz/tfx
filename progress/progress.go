// Package progress provides beautiful progress bars and spinners with multipath API support.
//
// This package enforces TermFX multipath pattern (see MULTIPATH.md):
//   - EXPRESS: Start(args...)             // default, config struct, functional options
//   - CONFIG:  StartWith(cfg ProgressConfig) // typed config entry for IDE autocompletion
//   - FLUENT:  Start(opts...)
//
// Object lifecycle uses typed constructors:
//   - New()                                 // default instance
//   - NewWithConfig(cfg ProgressConfig)     // explicit config instance
package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/terminal"
)

// Progress tracks and displays progress of operations.
type Progress struct {
	total     int
	current   int
	label     string
	width     int
	started   bool
	startTime time.Time
	theme     ProgressTheme
	style     ProgressStyle
	effect    ProgressEffect
	writer    io.Writer
	detector  *terminal.Detector
	mu        sync.Mutex
}

// ProgressConfig provides structured configuration for Progress.
type ProgressConfig struct {
	Total  int
	Label  string
	Width  int
	Theme  ProgressTheme
	Style  ProgressStyle
	Effect ProgressEffect
	Writer io.Writer
}

// DefaultProgressConfig returns default config for Progress.
func DefaultProgressConfig() ProgressConfig {
	return ProgressConfig{
		Total:  100,
		Label:  "Progress",
		Width:  40,
		Theme:  MaterialTheme,
		Style:  ProgressStyleBar,
		Effect: EffectNone,
		Writer: os.Stdout,
	}
}

// Start creates and starts a progress bar (EXPRESS API).
// Supports multipath: Start(), Start(cfg), Start(opts...).
func Start(args ...any) *Progress {
	p := newProgress(args...)
	p.Start()
	return p
}

// StartWith creates and starts a progress bar with explicit config (IDE SUPPORT API).
func StartWith(cfg ProgressConfig) *Progress {
	return Start(cfg)
}

// New creates a new Progress instance with defaults (OBJECT API).
func New() *Progress {
	return newProgress()
}

// NewWithConfig creates a new Progress instance with explicit config (OBJECT API).
func NewWithConfig(cfg ProgressConfig) *Progress {
	return newProgress(cfg)
}

// newProgress is the internal implementation supporting multipath overload using OverloadWithOptions.
func newProgress(args ...any) *Progress {
	// Separate functional options from other args
	var opts []share.Option[ProgressConfig]
	var cfgArgs []any
	for _, arg := range args {
		switch v := arg.(type) {
		case share.Option[ProgressConfig]:
			opts = append(opts, v)
		default:
			cfgArgs = append(cfgArgs, v)
		}
	}

	// Merge config and options
	cfg := share.OverloadWithOptions(cfgArgs, DefaultProgressConfig(), opts...)

	// Ensure defaults for zero values
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	if cfg.Total == 0 {
		cfg.Total = 100
	}
	if cfg.Label == "" {
		cfg.Label = "Progress"
	}
	if cfg.Width == 0 {
		cfg.Width = 40
	}

	p := &Progress{
		total:    cfg.Total,
		label:    cfg.Label,
		width:    cfg.Width,
		theme:    cfg.Theme,
		style:    cfg.Style,
		effect:   cfg.Effect,
		writer:   cfg.Writer,
		detector: terminal.NewDetector(cfg.Writer),
	}
	return p
}

// --- FUNCTIONAL OPTIONS ---

// WithTotal sets the total value for progress tracking
func WithTotal(total int) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Total = total
	}
}

// WithLabel sets the progress label
func WithLabel(label string) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Label = label
	}
}

// WithProgressWidth sets the width of the progress bar
func WithProgressWidth(width int) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Width = width
	}
}

// WithProgressTheme applies a ProgressTheme
func WithProgressTheme(theme ProgressTheme) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Theme = theme
	}
}

// WithProgressStyle applies a ProgressStyle
func WithProgressStyle(style ProgressStyle) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Style = style
	}
}

// WithProgressEffect applies visual effects
func WithProgressEffect(effect ProgressEffect) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Effect = effect
	}
}

// WithProgressWriter sets a custom writer for progress output
func WithProgressWriter(writer io.Writer) share.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Writer = writer
	}
}

// --- CONVENIENCE OPTIONS ---

// WithMaterialTheme applies Material Design theme
func WithMaterialTheme() share.Option[ProgressConfig] {
	return WithProgressTheme(MaterialTheme)
}

// WithDraculaTheme applies Dracula theme
func WithDraculaTheme() share.Option[ProgressConfig] {
	return WithProgressTheme(DraculaTheme)
}

// WithNordTheme applies Nord theme
func WithNordTheme() share.Option[ProgressConfig] {
	return WithProgressTheme(NordTheme)
}

// WithRainbowEffect enables rainbow effect
func WithRainbowEffect() share.Option[ProgressConfig] {
	return WithProgressEffect(EffectRainbow)
}

// WithGradientEffect enables gradient effect
func WithGradientEffect() share.Option[ProgressConfig] {
	return WithProgressEffect(EffectGradient)
}

// --- STYLE OPTIONS ---

// WithBarStyle uses classic progress bar style
func WithBarStyle() share.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleBar)
}

// WithDotsStyle uses dots progress style
func WithDotsStyle() share.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleDots)
}

// WithArrowsStyle uses arrows progress style
func WithArrowsStyle() share.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleArrows)
}

// --- PROGRESS METHODS ---

// Start initializes the progress tracking.
func (p *Progress) Start() {
	p.mu.Lock()
	if !p.started {
		p.started = true
		p.startTime = time.Now()
		p.render()
	}
	p.mu.Unlock()
}

// Set updates the current progress value.
func (p *Progress) Set(value int) {
	p.mu.Lock()
	if !p.started {
		p.started = true
		p.startTime = time.Now()
	}
	if value > p.total {
		value = p.total
	}
	p.current = value
	p.render()
	p.mu.Unlock()
}

// Complete marks progress as completed with success message.
func (p *Progress) Complete(msg string) {
	p.mu.Lock()
	p.current = p.total
	// Clear current line and render final progress with completion message
	fmt.Fprint(p.writer, "\r")
	output := RenderBar(p)
	completion := RenderCompletion(p.theme, msg, true, p.detector)
	fmt.Fprintln(p.writer, output+completion)
	p.mu.Unlock()
}

// Fail marks progress as failed with error message.
func (p *Progress) Fail(msg string) {
	p.mu.Lock()
	p.current = p.total
	// Clear current line and render final progress with failure message
	fmt.Fprint(p.writer, "\r")
	output := RenderBar(p)
	completion := RenderCompletion(p.theme, msg, false, p.detector)
	fmt.Fprintln(p.writer, output+completion)
	p.mu.Unlock()
}

// render displays the current progress state.
func (p *Progress) render() {
	output := RenderBar(p)
	fmt.Fprint(p.writer, output)
	// Don't add newline here - let Complete()/Fail() handle final rendering
}

// Increment increases progress by 1
func (p *Progress) Increment() {
	p.Set(p.current + 1)
}

// Add increases progress by specified amount
func (p *Progress) Add(amount int) {
	p.Set(p.current + amount)
}

// SetLabel updates the progress label
func (p *Progress) SetLabel(label string) {
	p.mu.Lock()
	p.label = label
	if p.started {
		p.render()
	}
	p.mu.Unlock()
}

// SetTheme updates the progress theme
func (p *Progress) SetTheme(theme ProgressTheme) {
	p.mu.Lock()
	p.theme = theme
	if p.started {
		p.render()
	}
	p.mu.Unlock()
}

// SetEffect updates the visual effect
func (p *Progress) SetEffect(effect ProgressEffect) {
	p.mu.Lock()
	p.effect = effect
	p.theme.EffectEnabled = (effect != EffectNone)
	if p.started {
		p.render()
	}
	p.mu.Unlock()
}

// GetPercent returns the current progress percentage
func (p *Progress) GetPercent() float64 {
	if p.total == 0 {
		return 0
	}
	return float64(p.current) / float64(p.total) * 100
}

// GetElapsed returns elapsed time since start
func (p *Progress) GetElapsed() time.Duration {
	if !p.started {
		return 0
	}
	return time.Since(p.startTime)
}

// GetETA estimates remaining time
func (p *Progress) GetETA() time.Duration {
	if !p.started || p.current == 0 {
		return 0
	}

	elapsed := p.GetElapsed()
	rate := float64(p.current) / elapsed.Seconds()
	remaining := float64(p.total-p.current) / rate

	return time.Duration(remaining * float64(time.Second))
}

// --- GLOBAL CONVENIENCE FUNCTIONS ---

var globalProgress *Progress
var globalProgressMu sync.Mutex

// StartGlobalProgress starts global progress tracking (EXPRESS API)
func StartGlobalProgress(total int, label string) {
	globalProgressMu.Lock()
	globalProgress = Start(total, label)
	globalProgressMu.Unlock()
	globalProgress.Start()
}

// Set sets global progress value
func Set(value int) {
	globalProgressMu.Lock()
	if globalProgress != nil {
		globalProgress.Set(value)
	}
	globalProgressMu.Unlock()
}

// Increment increments global progress
func Increment() {
	globalProgressMu.Lock()
	if globalProgress != nil {
		globalProgress.Increment()
	}
	globalProgressMu.Unlock()
}

// Complete completes global progress
func Complete(msg string) {
	globalProgressMu.Lock()
	if globalProgress != nil {
		globalProgress.Complete(msg)
	}
	globalProgressMu.Unlock()
}

// --- PRESET CONSTRUCTORS ---

// NewProgress creates a progress bar with custom configuration
func NewProgress(total int, label string) *Progress {
	return newProgress(ProgressConfig{
		Total: total,
		Label: label,
	})
}

// NewMaterialProgress creates a progress bar with Material Design theme
func NewMaterialProgress(total int, label string) *Progress {
	return newProgress(
		ProgressConfig{Total: total, Label: label},
		WithMaterialTheme(),
	)
}

// NewDraculaProgress creates a progress bar with Dracula theme and rainbow effect
func NewDraculaProgress(total int, label string) *Progress {
	return newProgress(
		ProgressConfig{Total: total, Label: label},
		WithDraculaTheme(),
		WithRainbowEffect(),
	)
}

// NewNordProgress creates a progress bar with Nord theme
func NewNordProgress(total int, label string) *Progress {
	return newProgress(
		ProgressConfig{Total: total, Label: label},
		WithNordTheme(),
	)
}

// NewRainbowProgress creates a progress bar with rainbow effect
func NewRainbowProgress(total int, label string) *Progress {
	return newProgress(
		ProgressConfig{Total: total, Label: label},
		WithRainbowEffect(),
	)
}
