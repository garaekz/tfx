package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProgressBasic(t *testing.T) {
	buf := &bytes.Buffer{}
	p := New(4, "Test", WithWriter(buf), WithWidth(4))
	p.Start()
	p.Set(2)
	p.Complete("Done")
	out := buf.String()
	if !strings.Contains(out, "Test") {
		t.Errorf("expected label 'Test' in output, got %q", out)
	}
	if !strings.Contains(out, "50%") {
		t.Errorf("expected '50%%' in output, got %q", out)
	}
	if !strings.HasSuffix(out, "\n Done\n") {
		t.Errorf("expected suffix '\\n Done\\n', got %q", out)
	}
}

func TestSpinnerStop(t *testing.T) {
	buf := &bytes.Buffer{}
	spin := NewSpinner("Wait", WithSpinnerWriter(buf), WithSpinnerFrames([]string{"*"}), WithSpinnerInterval(10*time.Millisecond))
	spin.Start()
	// allow a few frames to render
	time.Sleep(30 * time.Millisecond)
	spin.Stop("OK")
	out := buf.String()
	// Should contain the stop message
	if !strings.Contains(out, "OK") {
		t.Errorf("expected 'OK' in spinner output, got %q", out)
	}
}
