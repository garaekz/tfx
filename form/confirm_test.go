package form

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
)

// MockReader for testing ReadLine
type MockReader struct {
	lines []string
	index int
}

func (m *MockReader) ReadLine(ctx context.Context) (string, error) {
	if m.index >= len(m.lines) {
		return "", io.EOF // Simulate end of input
	}
	line := m.lines[m.index]
	m.index++
	return line, nil
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		defaultVal     bool
		expected       bool
		expectedError  error
		expectedOutput string
	}{
		{
			name:           "Yes input",
			input:          []string{"y"},
			defaultVal:     false,
			expected:       true,
			expectedError:  nil,
			expectedOutput: "Do you want to proceed? (y/N) ",
		},
		{
			name:           "No input",
			input:          []string{"n"},
			defaultVal:     true,
			expected:       false,
			expectedError:  nil,
			expectedOutput: "Do you want to proceed? (Y/n) ",
		},
		{
			name:           "Empty input, default false",
			input:          []string{""},
			defaultVal:     false,
			expected:       false,
			expectedError:  nil,
			expectedOutput: "Do you want to proceed? (y/N) ",
		},
		{
			name:           "Empty input, default true",
			input:          []string{""},
			defaultVal:     true,
			expected:       true,
			expectedError:  nil,
			expectedOutput: "Do you want to proceed? (Y/n) ",
		},
		{
			name:           "Invalid then yes",
			input:          []string{"abc", "yes"},
			defaultVal:     false,
			expected:       true,
			expectedError:  nil,
			expectedOutput: "Do you want to proceed? (y/N) Please enter 'y' or 'n'.\nDo you want to proceed? (y/N) ",
		},
		{
			name:           "EOF error",
			input:          []string{},
			defaultVal:     false,
			expected:       false,
			expectedError:  io.EOF,
			expectedOutput: "Do you want to proceed? (y/N) ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			reader := &MockReader{lines: tt.input}

			cfg := DefaultConfirmConfig()
			cfg.Label = "Do you want to proceed?"
			cfg.Default = tt.defaultVal
			cfg.Writer = buf
			cfg.Reader = reader

			result, err := ConfirmWithConfig(cfg)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && result != tt.expected {
				t.Errorf("Expected result %t, got %t", tt.expected, result)
			}
			if buf.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, buf.String())
			}
		})
	}
}
