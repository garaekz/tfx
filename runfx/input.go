package runfx

import (
	"bufio"
	"context"
	"io"
	"os"
)

// KeyReader handles reading keyboard input from a terminal and converting it to Key events.
type KeyReader struct {
	reader *bufio.Reader
	input  io.Reader
}

// NewKeyReader creates a new keyboard input reader.
func NewKeyReader(input io.Reader) *KeyReader {
	if input == nil {
		input = os.Stdin
	}
	return &KeyReader{
		reader: bufio.NewReader(input),
		input:  input,
	}
}

// ReadKey reads the next keyboard input and returns the corresponding Key.
// This function blocks until a key is pressed or the context is cancelled.
func (kr *KeyReader) ReadKey(ctx context.Context) (Key, error) {
	// Create a channel to receive the key
	keyCh := make(chan Key, 1)
	errCh := make(chan error, 1)

	// Read input in a goroutine
	go func() {
		key, err := kr.readKeyBlocking()
		if err != nil {
			errCh <- err
			return
		}
		keyCh <- key
	}()

	// Wait for either a key or context cancellation
	select {
	case key := <-keyCh:
		return key, nil
	case err := <-errCh:
		return KeyUnknown, err
	case <-ctx.Done():
		return KeyUnknown, ctx.Err()
	}
}

// readKeyBlocking performs the actual blocking read and key parsing.
func (kr *KeyReader) readKeyBlocking() (Key, error) {
	b, err := kr.reader.ReadByte()
	if err != nil {
		return KeyUnknown, err
	}

	// Handle escape sequences (arrow keys, etc.)
	if b == 27 { // ESC
		// Peek ahead to see if this is an escape sequence
		next, err := kr.reader.Peek(1)
		if err != nil || len(next) == 0 {
			// Just ESC key
			return KeyEscape, nil
		}

		if next[0] == '[' {
			// ANSI escape sequence
			kr.reader.ReadByte() // consume '['
			return kr.parseANSISequence()
		}

		// Just ESC key
		return KeyEscape, nil
	}

	// Handle regular characters
	return kr.parseRegularKey(b), nil
}

// parseANSISequence parses ANSI escape sequences like arrow keys.
func (kr *KeyReader) parseANSISequence() (Key, error) {
	b, err := kr.reader.ReadByte()
	if err != nil {
		return KeyUnknown, err
	}

	switch b {
	case 'A':
		return KeyArrowUp, nil
	case 'B':
		return KeyArrowDown, nil
	case 'C':
		return KeyArrowRight, nil
	case 'D':
		return KeyArrowLeft, nil
	case '3':
		// Delete key (ESC[3~)
		kr.reader.ReadByte() // consume '~'
		return KeyDelete, nil
	default:
		// Unknown sequence
		return KeyUnknown, nil
	}
}

// parseRegularKey converts regular byte input to Key constants.
func (kr *KeyReader) parseRegularKey(b byte) Key {
	switch b {
	case '\r', '\n':
		return KeyEnter
	case '\t':
		return KeyTab
	case ' ':
		return KeySpace
	case 127, 8: // DEL or Backspace
		return KeyBackspace
	case 3: // Ctrl+C
		return KeyCtrlC
	case 4: // Ctrl+D
		return KeyCtrlD

	// Letters (convert to uppercase)
	case 'a', 'A':
		return KeyA
	case 's', 'S':
		return KeyS
	case 'd', 'D':
		return KeyD
	case 'w', 'W':
		return KeyW
	case 'q', 'Q':
		return KeyQ
	case 'y', 'Y':
		return KeyY
	case 'n', 'N':
		return KeyN

	// Numbers
	case '0':
		return Key0
	case '1':
		return Key1
	case '2':
		return Key2
	case '3':
		return Key3
	case '4':
		return Key4
	case '5':
		return Key5
	case '6':
		return Key6
	case '7':
		return Key7
	case '8':
		return Key8
	case '9':
		return Key9

	default:
		return KeyUnknown
	}
}
