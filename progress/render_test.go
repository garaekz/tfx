package progress

import (
	"bytes"
	"strings"
	"testing"

	"github.com/garaekz/tfx/terminal"
)

func TestRenderBar(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("Test"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
	)
	p.Set(5) // 50%
	
	result := RenderBar(p)
	
	if !strings.Contains(result, "Test") {
		t.Errorf("expected 'Test' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "50%") {
		t.Errorf("expected '50%%' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
		t.Errorf("expected borders in rendered bar, got %q", result)
	}
}

func TestRenderBarWithEffects(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("Rainbow"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
		WithRainbowEffect(),
	)
	p.Set(7) // 70%
	
	result := RenderBar(p)
	
	if !strings.Contains(result, "Rainbow") {
		t.Errorf("expected 'Rainbow' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "70%") {
		t.Errorf("expected '70%%' in rendered bar, got %q", result)
	}
}

func TestRenderBarZeroProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("Empty"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
	)
	// Don't set any progress (0%)
	
	result := RenderBar(p)
	
	if !strings.Contains(result, "Empty") {
		t.Errorf("expected 'Empty' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "0%") {
		t.Errorf("expected '0%%' in rendered bar, got %q", result)
	}
}

func TestRenderBarCompleteProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("Complete"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
	)
	p.Set(10) // 100%
	
	result := RenderBar(p)
	
	if !strings.Contains(result, "Complete") {
		t.Errorf("expected 'Complete' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "100%") {
		t.Errorf("expected '100%%' in rendered bar, got %q", result)
	}
}

func TestRenderSpinner(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWith(
		WithMessage("Loading"),
		WithSpinnerWriter(buf),
	)
	
	result := RenderSpinner(s, "|")
	
	if !strings.Contains(result, "Loading") {
		t.Errorf("expected 'Loading' in rendered spinner, got %q", result)
	}
	if !strings.Contains(result, "|") {
		t.Errorf("expected '|' frame in rendered spinner, got %q", result)
	}
	// Should start with \r for overwriting
	if !strings.HasPrefix(result, "\r") {
		t.Errorf("expected result to start with \\r, got %q", result)
	}
}

func TestRenderSpinnerEmptyMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWith(
		WithMessage(""),
		WithSpinnerWriter(buf),
	)
	
	result := RenderSpinner(s, "*")
	
	if !strings.Contains(result, "*") {
		t.Errorf("expected '*' frame in rendered spinner, got %q", result)
	}
}

func TestRenderCompletion(t *testing.T) {
	tests := []struct {
		name     string
		success  bool
		message  string
		expected string
	}{
		{"Success", true, "Done!", "✅"},
		{"Failure", false, "Error", "❌"},
	}
	
	detector := terminal.NewDetector(&bytes.Buffer{})
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderCompletion(MaterialTheme, tt.message, tt.success, detector)
			
			if !strings.Contains(result, tt.message) {
				t.Errorf("expected message '%s' in result, got %q", tt.message, result)
			}
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected icon '%s' in result, got %q", tt.expected, result)
			}
		})
	}
}

func TestRenderAdvancedBar(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(100),
		WithLabel("Advanced"),
		WithProgressWriter(buf),
		WithProgressWidth(20),
	)
	
	p.Start()
	p.Set(25) // 25%
	
	result := RenderAdvancedBar(p)
	
	if !strings.Contains(result, "Advanced") {
		t.Errorf("expected 'Advanced' in rendered bar, got %q", result)
	}
	if !strings.Contains(result, "25%") {
		t.Errorf("expected '25%%' in rendered bar, got %q", result)
	}
	// Should contain ETA since progress is started and > 0
	if !strings.Contains(result, "ETA:") {
		t.Errorf("expected 'ETA:' in advanced bar, got %q", result)
	}
}

func TestRenderAdvancedBarNoProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(100),
		WithLabel("NoProgress"),
		WithProgressWriter(buf),
		WithProgressWidth(20),
	)
	
	// Don't start or set progress
	result := RenderAdvancedBar(p)
	
	if !strings.Contains(result, "NoProgress") {
		t.Errorf("expected 'NoProgress' in rendered bar, got %q", result)
	}
	// Should not contain ETA since not started
	if strings.Contains(result, "ETA:") {
		t.Errorf("unexpected 'ETA:' in bar without progress, got %q", result)
	}
}

func TestRenderMultiProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	
	// Create multiple progress bars
	bars := []*Progress{
		NewProgressWith(
			WithTotal(10),
			WithLabel("Task 1"),
			WithProgressWriter(buf),
		),
		NewProgressWith(
			WithTotal(20),
			WithLabel("Task 2"),
			WithProgressWriter(buf),
		),
		NewProgressWith(
			WithTotal(15),
			WithLabel("Task 3"),
			WithProgressWriter(buf),
		),
	}
	
	// Set some progress
	bars[0].Set(5)  // 50%
	bars[1].Set(10) // 50%
	bars[2].Set(3)  // 20%
	
	result := RenderMultiProgress(bars)
	
	// Check that all tasks are present
	if !strings.Contains(result, "Task 1") {
		t.Errorf("expected 'Task 1' in multi-progress, got %q", result)
	}
	if !strings.Contains(result, "Task 2") {
		t.Errorf("expected 'Task 2' in multi-progress, got %q", result)
	}
	if !strings.Contains(result, "Task 3") {
		t.Errorf("expected 'Task 3' in multi-progress, got %q", result)
	}
	
	// Check numbering
	if !strings.Contains(result, "1. ") {
		t.Errorf("expected '1. ' numbering in multi-progress, got %q", result)
	}
	if !strings.Contains(result, "2. ") {
		t.Errorf("expected '2. ' numbering in multi-progress, got %q", result)
	}
	if !strings.Contains(result, "3. ") {
		t.Errorf("expected '3. ' numbering in multi-progress, got %q", result)
	}
	
	// Check newlines (multi-line output)
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines in multi-progress, got %d lines", len(lines))
	}
}

func TestRenderMultiProgressEmpty(t *testing.T) {
	result := RenderMultiProgress([]*Progress{})
	if result != "" {
		t.Errorf("expected empty string for empty progress bars, got %q", result)
	}
}

func TestRenderASCIIBar(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("ASCII"),
		WithProgressWriter(buf),
		WithProgressWidth(10),
	)
	p.Set(6) // 60%
	
	tests := []struct {
		name     string
		style    ASCIIStyle
		expected string
	}{
		{"Classic", ASCIIClassic, ">"},
		{"Dots", ASCIIDots, "●"},
		{"Blocks", ASCIIBlocks, "█"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderASCIIBar(p, tt.style)
			
			if !strings.Contains(result, "ASCII") {
				t.Errorf("expected 'ASCII' in rendered bar, got %q", result)
			}
			if !strings.Contains(result, "60%") {
				t.Errorf("expected '60%%' in rendered bar, got %q", result)
			}
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected '%s' character in %s style, got %q", tt.expected, tt.name, result)
			}
		})
	}
}

func TestRenderResponsiveBar(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("Responsive"),
		WithProgressWriter(buf),
		WithProgressWidth(50), // Start with wide width
	)
	p.Set(5) // 50%
	
	// Test with small terminal width
	smallResult := RenderResponsiveBar(p, 30)
	if !strings.Contains(smallResult, "Responsive") {
		t.Errorf("expected 'Responsive' in small responsive bar, got %q", smallResult)
	}
	
	// Test with large terminal width
	largeResult := RenderResponsiveBar(p, 150)
	if !strings.Contains(largeResult, "Responsive") {
		t.Errorf("expected 'Responsive' in large responsive bar, got %q", largeResult)
	}
	
	// Original width should be preserved
	if p.width != 50 {
		t.Errorf("expected original width to be preserved (50), got %d", p.width)
	}
}

func TestRenderResponsiveBarMinimumWidth(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewProgressWith(
		WithTotal(10),
		WithLabel("VeryLongLabelThatTakesUpSpace"),
		WithProgressWriter(buf),
		WithProgressWidth(20),
	)
	p.Set(5)
	
	// Test with very small terminal width
	result := RenderResponsiveBar(p, 10)
	if !strings.Contains(result, "VeryLongLabelThatTakesUpSpace") {
		t.Errorf("expected long label in result, got %q", result)
	}
	
	// Should still render something despite small width
	if len(result) == 0 {
		t.Error("expected non-empty result even with small width")
	}
}

func TestRenderBarWithThemes(t *testing.T) {
	buf := &bytes.Buffer{}
	
	tests := []struct {
		name  string
		theme ProgressTheme
	}{
		{"Material", MaterialTheme},
		{"Dracula", DraculaTheme},
		{"Nord", NordTheme},
		{"GitHub", GitHubTheme},
		{"Tailwind", TailwindTheme},
		{"VSCode", VSCodeTheme},
		{"Rainbow", RainbowTheme},
		{"Neon", NeonTheme},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProgressWith(
				WithTotal(10),
				WithLabel(tt.name),
				WithProgressWriter(buf),
				WithProgressTheme(tt.theme),
			)
			p.Set(5)
			
			result := RenderBar(p)
			
			if !strings.Contains(result, tt.name) {
				t.Errorf("expected theme name '%s' in result, got %q", tt.name, result)
			}
			if !strings.Contains(result, "50%") {
				t.Errorf("expected '50%%' in result, got %q", result)
			}
		})
	}
}