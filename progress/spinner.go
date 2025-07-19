package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Spinner renders an animated spinner with an optional message.
type Spinner struct {
	frames   []string
	interval time.Duration
	message  string
	theme    ProgressTheme
	writer   io.Writer
	mu       sync.Mutex
	active   bool
	stopCh   chan struct{}
}

// Default spinner frames and interval.
var (
	DefaultSpinnerFrames   = []string{"|", "/", "-", "\\"}
	DefaultSpinnerInterval = 100 * time.Millisecond
)

type SpinnerOption func(*Spinner)

// WithSpinnerWriter sets a custom writer for spinner output.
func WithSpinnerWriter(w io.Writer) SpinnerOption { return func(s *Spinner) { s.writer = w } }

// WithSpinnerFrames sets custom frames for the spinner animation.
func WithSpinnerFrames(frames []string) SpinnerOption { return func(s *Spinner) { s.frames = frames } }

// WithSpinnerInterval sets the frame interval duration.
func WithSpinnerInterval(d time.Duration) SpinnerOption { return func(s *Spinner) { s.interval = d } }

// WithSpinnerTheme applies a ProgressTheme to colorize the spinner.
func WithSpinnerTheme(t ProgressTheme) SpinnerOption { return func(s *Spinner) { s.theme = t } }

// NewSpinner creates a new Spinner instance with defaults and any options.
func NewSpinner(message string, opts ...SpinnerOption) *Spinner {
	s := &Spinner{
		frames:   DefaultSpinnerFrames,
		interval: DefaultSpinnerInterval,
		message:  message,
		theme:    defaultTheme,
		writer:   os.Stdout,
		stopCh:   make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start begins the spinner in a new goroutine.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()
	go func() {
		i := 0
		for {
			select {
			case <-s.stopCh:
				return
			default:
				s.mu.Lock()
				frame := s.frames[i%len(s.frames)]
				i++
				s.mu.Unlock()
				// Render frame + message
				text := fmt.Sprintf("\r%s %s", frame, s.message)
				fmt.Fprint(s.writer, s.theme.LabelColor+text+"\033[0m")
				time.Sleep(s.interval)
			}
		}
	}()
}

// Stop stops the spinner and prints a success message.
func (s *Spinner) Stop(msg string) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.stopCh <- struct{}{}
	// Clear line and print done message
	text := fmt.Sprintf("\r%s %s\n", s.frames[0], msg)
	fmt.Fprint(s.writer, s.theme.CompleteColor+text+"\033[0m")
	s.mu.Unlock()
}

// Fail stops the spinner and prints a failure message.
func (s *Spinner) Fail(msg string) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.stopCh <- struct{}{}
	// Clear line and print failure message
	text := fmt.Sprintf("\r%s %s\n", s.frames[0], msg)
	fmt.Fprint(s.writer, s.theme.IncompleteColor+text+"\033[0m")
	s.mu.Unlock()
}
