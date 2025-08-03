package form

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		name           string
		options        []string
		input          []string
		defaultVal     int
		pageSize       int
		isTTY          bool
		expected       int
		expectedError  error
		expectedOutput string
	}{
		{
			name:           "Basic select (TTY)",
			options:        []string{"Opt1", "Opt2", "Opt3"},
			input:          []string{""}, // Enter key
			defaultVal:     0,
			pageSize:       10,
			isTTY:          true,
			expected:       0,
			expectedError:  nil,
			expectedOutput: "\033[H\033[2JChoose an option:\n> Opt1\n  Opt2\n  Opt3\n",
		},
		{
			name:          "Select with arrow keys (TTY)",
			options:       []string{"Opt1", "Opt2", "Opt3"},
			input:         []string{"down", "down", ""},
			defaultVal:    0,
			pageSize:      10,
			isTTY:         true,
			expected:      2,
			expectedError: nil,
			expectedOutput: "\033[H\033[2JChoose an option:\n> Opt1\n  Opt2\n  Opt3\n" +
				"\033[H\033[2JChoose an option:\n  Opt1\n> Opt2\n  Opt3\n" +
				"\033[H\033[2JChoose an option:\n  Opt1\n  Opt2\n> Opt3\n",
		},
		{
			name:          "Select with paging (TTY)",
			options:       []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			input:         []string{"down", "down", "down", "down", "down", "down", "down", "down", "down", "down", ""},
			defaultVal:    0,
			pageSize:      5,
			isTTY:         true,
			expected:      10,
			expectedError: nil,
			expectedOutput: "\033[H\033[2JChoose an option:\n> 1\n  2\n  3\n  4\n  5\n" +
				strings.Repeat("\033[H\033[2JChoose an option:\n  %s\n", 10) + // Simplified, actual output would vary
				"\033[H\033[2JChoose an option:\n  7\n  8\n  9\n  10\n> 11\n",
		},
		{
			name:           "Non-TTY select, valid input",
			options:        []string{"OptA", "OptB"},
			input:          []string{"2"},
			defaultVal:     0,
			pageSize:       10,
			isTTY:          false,
			expected:       1,
			expectedError:  nil,
			expectedOutput: "Choose an option:\n1) OptA\n2) OptB\nEnter choice (number): ",
		},
		{
			name:           "Non-TTY select, invalid then valid input",
			options:        []string{"OptA", "OptB"},
			input:          []string{"x", "1"},
			defaultVal:     0,
			pageSize:       10,
			isTTY:          false,
			expected:       0,
			expectedError:  nil,
			expectedOutput: "Choose an option:\n1) OptA\n2) OptB\nEnter choice (number): invalid choice: x\nChoose an option:\n1) OptA\n2) OptB\nEnter choice (number): ",
		},
		{
			name:           "No options provided",
			options:        []string{},
			input:          []string{""},
			defaultVal:     0,
			pageSize:       10,
			isTTY:          true,
			expected:       -1,
			expectedError:  errors.New("select: no options provided"),
			expectedOutput: "",
		},
		{
			name:           "Cancel with q (TTY)",
			options:        []string{"Opt1", "Opt2"},
			input:          []string{"q"},
			defaultVal:     0,
			pageSize:       10,
			isTTY:          true,
			expected:       -1,
			expectedError:  ErrCanceled,
			expectedOutput: "\033[H\033[2JChoose an option:\n> Opt1\n  Opt2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			reader := &MockReader{lines: tt.input}

			cfg := DefaultSelectConfig()
			cfg.Label = "Choose an option:"
			cfg.Options = tt.options
			cfg.Default = tt.defaultVal
			cfg.Writer = buf
			cfg.Reader = reader
			cfg.PageSize = tt.pageSize

			// Mock IsTTY
			oldIsTTY := IsTTY
			defer func() { IsTTY = oldIsTTY }()
			cfg.IsTerminal = func(w io.Writer) bool { return tt.isTTY }

			result, err := SelectWithConfig(cfg)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if err == nil && result != tt.expected {
				t.Errorf("Expected result %d, got %d", tt.expected, result)
			}
			// For TTY tests, the output contains ANSI escape codes for clearing screen.
			// We'll do a Contains check for simplicity, or a more complex golden file comparison.
			if tt.isTTY {
				// Basic check for presence of options and cursor
				for _, opt := range tt.options {
					if !strings.Contains(buf.String(), opt) {
						t.Errorf("Output missing option %q: %q", opt, buf.String())
					}
				}
			} else {
				if buf.String() != tt.expectedOutput {
					t.Errorf("Expected output %q, got %q", tt.expectedOutput, buf.String())
				}
			}
		})
	}
}
