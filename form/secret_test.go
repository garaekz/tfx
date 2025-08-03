package form

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"golang.org/x/term"
)

// MockTermReadPassword for testing term.ReadPassword
type MockTermReadPassword struct {
	passwords []string
	index     int
	// Simulate TTY behavior
	isTTY bool
}

func (m *MockTermReadPassword) ReadPassword(fd int) ([]byte, error) {
	if !m.isTTY {
		return nil, errors.New("not a TTY")
	}
	if m.index >= len(m.passwords) {
		return nil, io.EOF
	}
	password := m.passwords[m.index]
	m.index++
	return []byte(password), nil
}

// MockIsTTY for testing IsTTY
type MockIsTTY struct {
	isTTY bool
}

func (m *MockIsTTY) IsTerminal(w io.Writer) bool {
	return m.isTTY
}

// MockMakeRaw for testing terminal.MakeRaw
type MockMakeRaw struct {
	returnError error
}

func (m *MockMakeRaw) MakeRaw(fd uintptr) (*term.State, error) {
	return nil, m.returnError
}

// MockRestoreTerminal for testing terminal.RestoreTerminal
type MockRestoreTerminal struct {
	returnError error
}

func (m *MockRestoreTerminal) RestoreTerminal(fd uintptr, state *term.State) error {
	return m.returnError
}

func TestSecret(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		confirm        bool
		isTTY          bool
		expected       string
		expectedError  error
		expectedOutput string
		makeRawError   error
	}{
		{
			name:           "Basic secret input (TTY)",
			input:          []string{"mysecret"},
			confirm:        false,
			isTTY:          true,
			expected:       "mysecret",
			expectedError:  nil,
			expectedOutput: "Enter secret: \n",
			makeRawError:   nil,
		},
		{
			name:           "Secret with confirmation (TTY, match)",
			input:          []string{"mysecret", "mysecret"},
			confirm:        true,
			isTTY:          true,
			expected:       "mysecret",
			expectedError:  nil,
			expectedOutput: "Enter secret: \nConfirm secret: \n",
			makeRawError:   nil,
		},
		{
			name:           "Secret with confirmation (TTY, mismatch then match)",
			input:          []string{"secret1", "secret2", "secret3", "secret3"},
			confirm:        true,
			isTTY:          true,
			expected:       "secret3",
			expectedError:  nil,
			expectedOutput: "Enter secret: \nConfirm secret: \nSecrets do not match. Please try again.\nEnter secret: \nConfirm secret: \n",
			makeRawError:   nil,
		},
		{
			name:           "Basic secret input (non-TTY)",
			input:          []string{"mysecret"},
			confirm:        false,
			isTTY:          false,
			expected:       "mysecret",
			expectedError:  nil,
			expectedOutput: "Enter secret: ",
			makeRawError:   nil,
		},
		{
			name:           "EOF error (TTY)",
			input:          []string{},
			confirm:        false,
			isTTY:          true,
			expected:       "",
			expectedError:  io.EOF,
			expectedOutput: "Enter secret: ",
			makeRawError:   nil,
		},
		{
			name:           "EOF error (non-TTY)",
			input:          []string{},
			confirm:        false,
			isTTY:          false,
			expected:       "",
			expectedError:  io.EOF,
			expectedOutput: "Enter secret: ",
			makeRawError:   nil,
		},
		{
			name:           "MakeRaw error",
			input:          []string{"secret"},
			confirm:        false,
			isTTY:          true,
			expected:       "",
			expectedError:  errors.New("failed to set raw terminal mode: mock make raw error"),
			expectedOutput: "",
			makeRawError:   errors.New("mock make raw error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)

			cfg := DefaultSecretConfig()

			// Mock ReadPassword
			mockReadPassword := &MockTermReadPassword{passwords: tt.input, isTTY: tt.isTTY}
			cfg.ReadPassword = mockReadPassword.ReadPassword

			// Mock IsTerminal
			mockIsTTY := &MockIsTTY{isTTY: tt.isTTY}
			cfg.IsTerminal = mockIsTTY.IsTerminal

			// Mock MakeRaw and RestoreTerminal
			mockMakeRaw := &MockMakeRaw{returnError: tt.makeRawError}
			cfg.MakeRaw = mockMakeRaw.MakeRaw
			cfg.RestoreTerminal = (&MockRestoreTerminal{}).RestoreTerminal // Restore always succeeds for now

			cfg.Label = "Enter secret"
			cfg.Writer = buf
			cfg.Confirm = tt.confirm
			// For non-TTY, we need to use the MockReader
			if !tt.isTTY {
				cfg.Reader = &MockReader{lines: tt.input}
			}

			result, err := SecretWithConfig(cfg)

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
				t.Errorf("Expected result %q, got %q", tt.expected, result)
			}
			if buf.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, buf.String())
			}
		})
	}
}
