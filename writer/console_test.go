package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/garaekz/tfx/color"
	. "github.com/garaekz/tfx/internal/share"
)

func TestConsoleWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := ConsoleOptions{
		Level:      LevelInfo,
		Format:     FormatBadge,
		Timestamp:  false,
		Theme:      color.DefaultTheme,
		ForceColor: false,
		BadgeWidth: 4,
	}

	cw := NewConsoleWriter(buf, opts)

	entry := &Entry{
		Level:   LevelInfo,
		Message: "hello",
		Fields:  Fields{},
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
