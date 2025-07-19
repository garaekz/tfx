package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/internal/core"
)

func TestConsoleWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Level:      core.LevelInfo,
		Format:     core.FormatBadge,
		Timestamp:  false,
		Theme:      color.DefaultTheme,
		ForceColor: false,
		BadgeWidth: 4,
	}

	cw := NewConsoleWriter(buf, opts)

	entry := &core.Entry{
		Level:   core.LevelInfo,
		Message: "hello",
		Fields:  core.Fields{},
	}

	if err := cw.Write(entry); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	want := "[INFO] hello"
	if got != want {
		t.Errorf("unexpected output: %q != %q", got, want)
	}
}
