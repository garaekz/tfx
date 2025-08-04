package runfx

import (
	"bufio"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
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
func (kr *KeyReader) ReadKey(ctx context.Context) (Key, error) {
	keyCh := make(chan Key, 1)
	errCh := make(chan error, 1)

	go func() {
		key, err := kr.readKeyBlocking()
		if err != nil {
			errCh <- err
			return
		}
		keyCh <- key
	}()

	select {
	case key := <-keyCh:
		return key, nil
	case err := <-errCh:
		return Key{Code: KeyUnknown}, err
	case <-ctx.Done():
		return Key{Code: KeyUnknown}, ctx.Err()
	}
}

func (kr *KeyReader) readKeyBlocking() (Key, error) {
	b, err := kr.reader.ReadByte()
	if err != nil {
		return Key{Code: KeyUnknown}, err
	}

	if b == 27 {
		next, err := kr.reader.Peek(1)
		if err != nil || len(next) == 0 {
			return Key{Code: KeyEscape}, nil
		}

		if next[0] == '[' {
			kr.reader.ReadByte() // consume '['
			return kr.parseCSISequence()
		}
		return Key{Code: KeyEscape}, nil
	}

	return kr.parseRegularKey(b), nil
}

func (kr *KeyReader) parseCSISequence() (Key, error) {
	seq := []byte{}

	for {
		b, err := kr.reader.ReadByte()
		if err != nil {
			return Key{Code: KeyUnknown}, err
		}
		seq = append(seq, b)
		if (b >= 'A' && b <= 'Z') || b == '~' {
			break
		}
	}

	return kr.decodeCSI(seq)
}

func (kr *KeyReader) decodeCSI(seq []byte) (Key, error) {
	s := string(seq)

	switch s {
	case "A":
		return Key{Code: KeyArrowUp}, nil
	case "B":
		return Key{Code: KeyArrowDown}, nil
	case "C":
		return Key{Code: KeyArrowRight}, nil
	case "D":
		return Key{Code: KeyArrowLeft}, nil
	case "3~":
		return Key{Code: KeyDelete}, nil
	}

	if strings.Contains(s, ";") {
		parts := strings.Split(s, ";")
		if len(parts) != 2 {
			return Key{Code: KeyUnknown}, nil
		}
		modNum, _ := strconv.Atoi(parts[1][:1])
		var mod Modifier
		switch modNum {
		case 2:
			mod = ModShift
		case 3:
			mod = ModAlt
		case 4:
			mod = ModShift | ModAlt
		case 5:
			mod = ModCtrl
		case 6:
			mod = ModCtrl | ModShift
		case 7:
			mod = ModCtrl | ModAlt
		case 8:
			mod = ModCtrl | ModAlt | ModShift
		default:
			mod = ModNone
		}
		last := parts[1][1:]
		switch last {
		case "A":
			return Key{Code: KeyArrowUp, Modifier: mod}, nil
		case "B":
			return Key{Code: KeyArrowDown, Modifier: mod}, nil
		case "C":
			return Key{Code: KeyArrowRight, Modifier: mod}, nil
		case "D":
			return Key{Code: KeyArrowLeft, Modifier: mod}, nil
		}
	}

	return Key{Code: KeyUnknown}, nil
}

func (kr *KeyReader) parseRegularKey(b byte) Key {
	switch b {
	case '\r', '\n':
		return Key{Code: KeyEnter}
	case '\t':
		return Key{Code: KeyTab}
	case ' ':
		return Key{Code: KeySpace}
	case 127, 8:
		return Key{Code: KeyBackspace}
	case 3:
		return Key{Code: KeyCtrlC, Modifier: ModCtrl}
	case 4:
		return Key{Code: KeyCtrlD, Modifier: ModCtrl}
	default:
		if b >= '0' && b <= '9' {
			return Key{Code: KeyCode(int(Key0) + int(b-'0')), Rune: rune(b)}
		}
		if b >= 'a' && b <= 'z' {
			return Key{Code: KeyCode(int(KeyA) + int(b-'a')), Rune: rune(b)}
		}
		if b >= 'A' && b <= 'Z' {
			return Key{Code: KeyCode(int(KeyA) + int(b-'A')), Rune: rune(b), Modifier: ModShift}
		}
		return Key{Code: KeyUnknown, Rune: rune(b)}
	}
}
