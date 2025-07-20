package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/garaekz/tfx/internal/shared"
	"github.com/garaekz/tfx/terminal"
)

// Progress tracks and displays progress of operations with color themes and effects
type Progress struct {
	total     int
	current   int
	label     string
	width     int
	done      bool
	started   bool
	startTime time.Time
	theme     ProgressTheme
	style     ProgressStyle
	effect    ProgressEffect
	writer    io.Writer
	detector  *terminal.Detector
	mu        sync.Mutex
}

// ProgressConfig provides structured configuration for Progress
type ProgressConfig struct {
	Total  int
	Label  string
	Width  int
	Theme  ProgressTheme
	Style  ProgressStyle
	Effect ProgressEffect
	Writer io.Writer
}


// DefaultProgressConfig returns default configuration for Progress
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

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick default
func StartProgress(total int, label string) *Progress {
	cfg := DefaultProgressConfig()
	cfg.Total = total
	cfg.Label = label
	return NewProgress(cfg)
}

// 2. INSTANTIATED: Config struct
func NewProgress(cfg ProgressConfig, opts ...shared.Option[ProgressConfig]) *Progress {
	// Apply functional options to config
	shared.ApplyOptions(&cfg, opts...)

	p := &Progress{
		total:  cfg.Total,
		label:  cfg.Label,
		width:  cfg.Width,
		theme:  cfg.Theme,
		style:  cfg.Style,
		effect: cfg.Effect,
		writer: cfg.Writer,
	}

	// Create terminal detector
	p.detector = terminal.NewDetector(p.writer)

	return p
}

// 3. FLUENT: Functional options
func NewProgressWith(opts ...shared.Option[ProgressConfig]) *Progress {
	cfg := DefaultProgressConfig()
	return NewProgress(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

// WithTotal sets the total value for progress tracking
func WithTotal(total int) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Total = total
	}
}

// WithLabel sets the progress label
func WithLabel(label string) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Label = label
	}
}

// WithProgressWidth sets the width of the progress bar
func WithProgressWidth(width int) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Width = width
	}
}

// WithProgressTheme applies a ProgressTheme
func WithProgressTheme(theme ProgressTheme) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Theme = theme
	}
}

// WithProgressStyle applies a ProgressStyle
func WithProgressStyle(style ProgressStyle) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Style = style
	}
}

// WithProgressEffect applies visual effects
func WithProgressEffect(effect ProgressEffect) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Effect = effect
	}
}

// WithProgressWriter sets a custom writer for progress output
func WithProgressWriter(writer io.Writer) shared.Option[ProgressConfig] {
	return func(cfg *ProgressConfig) {
		cfg.Writer = writer
	}
}

// --- CONVENIENCE OPTIONS ---

// WithMaterialTheme applies Material Design theme
func WithMaterialTheme() shared.Option[ProgressConfig] {
	return WithProgressTheme(MaterialTheme)
}

// WithDraculaTheme applies Dracula theme
func WithDraculaTheme() shared.Option[ProgressConfig] {
	return WithProgressTheme(DraculaTheme)
}

// WithNordTheme applies Nord theme
func WithNordTheme() shared.Option[ProgressConfig] {
	return WithProgressTheme(NordTheme)
}

// WithRainbowEffect enables rainbow effect
func WithRainbowEffect() shared.Option[ProgressConfig] {
	return WithProgressEffect(EffectRainbow)
}

// WithGradientEffect enables gradient effect
func WithGradientEffect() shared.Option[ProgressConfig] {
	return WithProgressEffect(EffectGradient)
}

// --- STYLE OPTIONS ---

// WithBarStyle uses classic progress bar style
func WithBarStyle() shared.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleBar)
}

// WithDotsStyle uses dots progress style
func WithDotsStyle() shared.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleDots)
}

// WithArrowsStyle uses arrows progress style
func WithArrowsStyle() shared.Option[ProgressConfig] {
	return WithProgressStyle(ProgressStyleArrows)
}

// --- PROGRESS METHODS ---

// Start initializes the progress tracking
func (p *Progress) Start() {
	p.mu.Lock()
	if !p.started {
		p.started = true
		p.startTime = time.Now()
		p.render()
	}
	p.mu.Unlock()
}

// Set updates the current progress value
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

// Increment increases progress by 1
func (p *Progress) Increment() {
	p.Set(p.current + 1)
}

// Add increases progress by specified amount
func (p *Progress) Add(amount int) {
	p.Set(p.current + amount)
}

// Complete marks progress as completed with success message
func (p *Progress) Complete(msg string) {
	p.mu.Lock()
	p.current = p.total
	p.done = true
	p.render()

	// Render completion message with theme
	completion := RenderCompletion(p.theme, msg, true, p.detector)
	fmt.Fprintln(p.writer, completion)
	p.mu.Unlock()
}

// Fail marks progress as failed with error message
func (p *Progress) Fail(msg string) {
	p.mu.Lock()
	p.done = true
	p.render()

	// Render failure message
	completion := RenderCompletion(p.theme, msg, false, p.detector)
	fmt.Fprintln(p.writer, completion)
	p.mu.Unlock()
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

// --- INTERNAL METHODS ---

// render displays the current progress state
func (p *Progress) render() {
	output := RenderBar(p)
	fmt.Fprint(p.writer, output)
	if p.done {
		fmt.Fprint(p.writer, "\n")
	}
}

// --- GLOBAL CONVENIENCE FUNCTIONS ---

var globalProgress *Progress
var globalProgressMu sync.Mutex

// StartGlobalProgress starts global progress tracking (EXPRESS API)
func StartGlobalProgress(total int, label string) {
	globalProgressMu.Lock()
	globalProgress = StartProgress(total, label)
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

// NewMaterialProgress creates a progress bar with Material Design theme
func NewMaterialProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithMaterialTheme(),
	)
}

// NewDraculaProgress creates a progress bar with Dracula theme and rainbow effect
func NewDraculaProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithDraculaTheme(),
		WithRainbowEffect(),
	)
}

// NewNordProgress creates a progress bar with Nord theme
func NewNordProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithNordTheme(),
	)
}

// NewRainbowProgress creates a progress bar with rainbow effect
func NewRainbowProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithProgressTheme(RainbowTheme),
		WithRainbowEffect(),
	)
}

// NewMinimalProgress creates a clean, minimal progress bar
func NewMinimalProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithMaterialTheme(),
		WithBarStyle(),
		WithProgressWidth(30),
	)
}

// NewWideProgress creates a wide progress bar for detailed tracking
func NewWideProgress(total int, label string) *Progress {
	return NewProgressWith(
		WithTotal(total),
		WithLabel(label),
		WithMaterialTheme(),
		WithProgressWidth(60),
	)
}