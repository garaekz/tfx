package form

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/garaekz/tfx/form/validate"
)

func TestInput(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		defaultVal     string
		validateFn     func(string) error
		expected       string
		expectedError  error
		expectedOutput string
	}{
		{
			name:           "Basic input",
			input:          []string{"hello"},
			defaultVal:     "",
			validateFn:     nil,
			expected:       "hello",
			expectedError:  nil,
			expectedOutput: "Enter value: ",
		},
		{
			name:           "Empty input, with default",
			input:          []string{""},
			defaultVal:     "default_val",
			validateFn:     nil,
			expected:       "default_val",
			expectedError:  nil,
			expectedOutput: "Enter value (default: default_val): ",
		},
		{
			name:           "Validation fail then success",
			input:          []string{"a", "valid"},
			defaultVal:     "",
			validateFn:     validate.MinLen(3),
			expected:       "valid",
			expectedError:  nil,
			expectedOutput: "Enter value: Validation error: input must be at least 3 characters long\nEnter value: ",
		},
		{
			name:           "EOF error",
			input:          []string{},
			defaultVal:     "",
			validateFn:     nil,
			expected:       "",
			expectedError:  io.EOF,
			expectedOutput: "Enter value: ",
		},
		{
			name:           "Regex validation",
			input:          []string{"abc", "test@example.com"},
			defaultVal:     "",
			validateFn:     validate.Matches(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)),
			expected:       "test@example.com",
			expectedError:  nil,
			expectedOutput: "Enter value: Validation error: input does not match required pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\nEnter value: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			reader := &MockReader{lines: tt.input}

			cfg := DefaultInputConfig()
			cfg.Label = "Enter value"
			cfg.Default = tt.defaultVal
			cfg.Writer = buf
			cfg.Reader = reader
			if tt.validateFn != nil {
				cfg.Validate = tt.validateFn
			}

			result, err := InputWithConfig(cfg)

			if !errors.Is(err, tt.expectedError) && err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && result != tt.expected {
				t.Errorf("Expected result %q, got %q", tt.expected, result)
			}
			if buf.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, buf.String())
			}
		})
	}
}
