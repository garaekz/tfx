package tfx

import (
	"bytes"
	"strings"
	"testing"
	"time"
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
