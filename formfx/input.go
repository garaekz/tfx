package formfx

import "github.com/garaekz/tfx/runfx"

// InputPrompt stores the state for a text input primitive.
type InputPrompt struct {
	Value      []rune
	CursorPos  int
	keyHandler TextKeyHandlerFunc

	Done     chan string
	Canceled chan struct{}
}

type TextKeyHandlerFunc func(p *InputPrompt, key runfx.Key) bool

// NewInputPrompt creates a new InputPrompt with a default value.
func NewInputPrompt(defaultValue string) *InputPrompt {
	return &InputPrompt{
		Value:      []rune(defaultValue),
		CursorPos:  len(defaultValue),
		Done:       make(chan string, 1),
		Canceled:   make(chan struct{}),
		keyHandler: TextInputKeyHandler,
	}
}

// SetKeyHandler sets a custom keyboard handler.
func (p *InputPrompt) SetKeyHandler(h TextKeyHandlerFunc) {
	p.keyHandler = h
}

// TextInputKeyHandler provides basic editing for text input.
func TextInputKeyHandler(p *InputPrompt, key runfx.Key) bool {
	if !key.IsPrintable() {
		return false
	}

	switch key.Code {
	case runfx.KeyEnter:
		p.Done <- string(p.Value)
		return true
	case runfx.KeyEscape, runfx.KeyCtrlC:
		close(p.Canceled)
		return true
	case runfx.KeyBackspace:
		if p.CursorPos > 0 {
			p.Value = append(p.Value[:p.CursorPos-1], p.Value[p.CursorPos:]...)
			p.CursorPos--
		}
	case runfx.KeyArrowLeft:
		if p.CursorPos > 0 {
			p.CursorPos--
		}
	case runfx.KeyArrowRight:
		if p.CursorPos < len(p.Value) {
			p.CursorPos++
		}
	default:
		if key.Rune != 0 {
			p.Value = append(p.Value[:p.CursorPos], append([]rune{key.Rune}, p.Value[p.CursorPos:]...)...)
			p.CursorPos++
		}
	}
	return false
}

// OnKey delegates to the configured key handler.
func (p *InputPrompt) OnKey(key runfx.Key) bool {
	if p.keyHandler != nil {
		return p.keyHandler(p, key)
	}
	return TextInputKeyHandler(p, key)
}

func (p *InputPrompt) Render() []byte          { return nil }
func (p *InputPrompt) OnResize(cols, rows int) {}
