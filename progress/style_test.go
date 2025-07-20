package progress

import "testing"

func TestProgressStyleConstants(t *testing.T) {
	// Test that all style constants are defined and have different values
	styles := []ProgressStyle{
		ProgressStyleBar,
		ProgressStyleDots,
		ProgressStyleArrows,
		ProgressStyleAscii,
	}
	
	// Check that they have different values
	for i, style1 := range styles {
		for j, style2 := range styles {
			if i != j && style1 == style2 {
				t.Errorf("styles at index %d and %d have same value %d", i, j, style1)
			}
		}
	}
}

func TestStyleChars(t *testing.T) {
	tests := []struct {
		style    ProgressStyle
		filled   string
		empty    string
	}{
		{ProgressStyleBar, "█", "░"},
		{ProgressStyleDots, "●", "○"},
		{ProgressStyleArrows, ">", "-"},
		{ProgressStyleAscii, "=", "-"},
	}
	
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			filled := tt.style.FilledChar()
			empty := tt.style.EmptyChar()
			
			if filled != tt.filled {
				t.Errorf("expected filled char '%s', got '%s'", tt.filled, filled)
			}
			
			if empty != tt.empty {
				t.Errorf("expected empty char '%s', got '%s'", tt.empty, empty)
			}
		})
	}
}

func TestStyleCharsMapping(t *testing.T) {
	// Test that all styles have entries in styleChars map
	styles := []ProgressStyle{
		ProgressStyleBar,
		ProgressStyleDots,
		ProgressStyleArrows,
		ProgressStyleAscii,
	}
	
	for _, style := range styles {
		chars, exists := styleChars[style]
		if !exists {
			t.Errorf("style %d not found in styleChars map", style)
			continue
		}
		
		if chars.Filled == "" {
			t.Errorf("style %d has empty filled character", style)
		}
		
		if chars.Empty == "" {
			t.Errorf("style %d has empty empty character", style)
		}
	}
}

func TestStyleCharsSameForFilled(t *testing.T) {
	// Test that FilledChar() returns same value as direct map access
	for style, chars := range styleChars {
		filled := style.FilledChar()
		if filled != chars.Filled {
			t.Errorf("FilledChar() returned '%s', but map has '%s' for style %d", 
				filled, chars.Filled, style)
		}
	}
}

func TestStyleCharsSameForEmpty(t *testing.T) {
	// Test that EmptyChar() returns same value as direct map access
	for style, chars := range styleChars {
		empty := style.EmptyChar()
		if empty != chars.Empty {
			t.Errorf("EmptyChar() returned '%s', but map has '%s' for style %d", 
				empty, chars.Empty, style)
		}
	}
}

func TestStyleIntegrationWithProgress(t *testing.T) {
	// Test that styles can be used with progress bars
	styles := []struct {
		style ProgressStyle
		name  string
	}{
		{ProgressStyleBar, "Bar"},
		{ProgressStyleDots, "Dots"},
		{ProgressStyleArrows, "Arrows"},
		{ProgressStyleAscii, "Ascii"},
	}
	
	for _, tt := range styles {
		t.Run(tt.name, func(t *testing.T) {
			// This tests that the style can be used in configuration
			p := NewProgressWith(
				WithTotal(10),
				WithLabel("Style Test"),
				WithProgressStyle(tt.style),
			)
			
			if p.style != tt.style {
				t.Errorf("expected style %d, got %d", tt.style, p.style)
			}
		})
	}
}