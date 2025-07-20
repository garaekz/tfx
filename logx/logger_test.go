package logx

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/garaekz/tfx/color"
)

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	logger := New(opts)
	logger.Info("hello")
	// allow async write to complete
	time.Sleep(10 * time.Millisecond)

	got := strings.TrimSpace(buf.String())
	want := "[INFO] hello"
	if got != want {
		t.Errorf("unexpected output: %q != %q", got, want)
	}
}
func TestGlobalFunctions_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	Configure(opts)
	Info("global info message")
	time.Sleep(10 * time.Millisecond)

	got := strings.TrimSpace(buf.String())
	want := "[INFO] global info message"
	if got != want {
		t.Errorf("unexpected output: %q != %q", got, want)
	}
}

func TestGlobalFunctions_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false

	Configure(opts)
	SetLevel(100) // set to a high level so Info won't log
	Info("should not appear")
	time.Sleep(10 * time.Millisecond)

	got := strings.TrimSpace(buf.String())
	if got != "" {
		t.Errorf("expected no output, got: %q", got)
	}
}

func TestGlobalFunctions_SetOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = os.Stdout // set to default first
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	Configure(opts)
	SetOutput(buf) // then change output
	Info("output changed")
	Flush() // ensure async writes are flushed

	got := strings.TrimSpace(buf.String())
	want := "[INFO] output changed"
	if got != want {
		t.Errorf("unexpected output: %q != %q", got, want)
	}
}

func TestGlobalFunctions_Success(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	Configure(opts)
	Success("operation succeeded")
	time.Sleep(10 * time.Millisecond)

	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "operation succeeded") {
		t.Errorf("expected success message, got: %q", got)
	}
}

func TestGlobalFunctions_Badge(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	Configure(opts)
	Badge("TEST", "badge message", color.Hex("#FF5733"))
	time.Sleep(10 * time.Millisecond)

	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "badge message") {
		t.Errorf("expected badge message, got: %q", got)
	}
}

func TestGlobalFunctions_EnableDisableTimestamp(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := DefaultOptions()
	opts.Output = buf
	opts.Timestamp = false
	opts.ForceColor = false
	opts.BadgeWidth = 4

	Configure(opts)
	EnableTimestamp()
	Info("timestamp enabled")
	time.Sleep(10 * time.Millisecond)
	got := buf.String()
	if !strings.Contains(got, "timestamp enabled") {
		t.Errorf("expected message, got: %q", got)
	}
	// crude check for timestamp: should contain a number (year)
	if !strings.Contains(got, "20") {
		t.Errorf("expected timestamp in output, got: %q", got)
	}

	buf.Reset()
	DisableTimestamp()
	Info("timestamp disabled")
	time.Sleep(10 * time.Millisecond)
	got = buf.String()
	if !strings.Contains(got, "timestamp disabled") {
		t.Errorf("expected message, got: %q", got)
	}
	// crude check for timestamp: should not contain a number (year)
	if strings.Contains(got, "20") {
		t.Errorf("did not expect timestamp in output, got: %q", got)
	}
}
