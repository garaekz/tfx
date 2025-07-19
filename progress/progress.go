package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Progress represents a terminal progress bar.
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
	writer    io.Writer
	mu        sync.Mutex
}

var (
	// Default settings for new progress bars.
	defaultTheme            = DraculaTheme
	defaultStyle            = ProgressStyleBar
	defaultWidth            = 40
	defaultWriter io.Writer = os.Stdout
)

type Option func(*Progress)

// WithTheme applies a ProgressTheme to the bar.
func WithTheme(t ProgressTheme) Option { return func(p *Progress) { p.theme = t } }

// WithStyle applies a ProgressStyle to the bar.
func WithStyle(s ProgressStyle) Option { return func(p *Progress) { p.style = s } }

// WithWidth sets the width of the bar.
func WithWidth(w int) Option { return func(p *Progress) { p.width = w } }

// WithWriter sets the io.Writer for rendering the bar.
func WithWriter(w io.Writer) Option { return func(p *Progress) { p.writer = w } }

// New creates a new Progress instance with defaults and any provided options.
func New(total int, label string, opts ...Option) *Progress {
	p := &Progress{
		total:  total,
		label:  label,
		width:  defaultWidth,
		theme:  defaultTheme,
		style:  defaultStyle,
		writer: defaultWriter,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Start is sugar: creates, starts, and returns a Progress bar in one call.
func Start(total int, label string, opts ...Option) *Progress {
	p := New(total, label, opts...)
	p.Start()
	return p
}

// Start begins rendering the progress bar.
func (p *Progress) Start() {
	p.mu.Lock()
	if !p.started {
		p.started = true
		p.startTime = time.Now()
		p.render()
	}
	p.mu.Unlock()
}

// Set moves the progress bar to the specified value.
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

// Add increments the progress bar by delta.
func (p *Progress) Add(delta int) {
	p.Set(p.current + delta)
}

// Complete finishes the progress bar and prints the trailing message.
func (p *Progress) Complete(msg string) {
	p.mu.Lock()
	p.current = p.total
	p.done = true
	p.render()
	fmt.Fprintln(p.writer, " "+msg)
	p.mu.Unlock()
}

// render delegates the visual construction to RenderBar and writes to the writer.
func (p *Progress) render() {
	output := RenderBar(p)
	fmt.Fprint(p.writer, output)
	if p.done {
		fmt.Fprint(p.writer, "\n")
	}
}
