package formfx

import (
	"fmt"
	"strings"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// SecretConfig holds configuration for a SecretPrompt.
type SecretConfig struct {
	Label    string
	Confirm  bool
	Mask     rune
	Renderer SecretRenderer
}

// DefaultSecretConfig returns default options for SecretPrompt.
func DefaultSecretConfig() SecretConfig {
	return SecretConfig{
		Label:    "Enter secret:",
		Confirm:  false,
		Mask:     '*',
		Renderer: &DefaultSecretRenderer{},
	}
}

// SecretRenderer defines the rendering interface for SecretPrompt.
type SecretRenderer interface {
	Render(s *SecretPrompt) []byte
}

// DefaultSecretRenderer renders the label and masked value.
type DefaultSecretRenderer struct{}

func (r *DefaultSecretRenderer) Render(s *SecretPrompt) []byte {
	maskedValue := strings.Repeat(string(s.Mask), len(s.prompt.Value))
	return fmt.Appendf(nil, "%s: %s", s.Label, maskedValue)
}

// SecretPrompt is a component for secret text input.
type SecretPrompt struct {
	prompt   *InputPrompt
	Label    string
	Confirm  bool
	Mask     rune
	renderer SecretRenderer
}

// NewSecretPrompt builds a SecretPrompt from configuration.
func NewSecretPrompt(cfg SecretConfig) (*SecretPrompt, error) {
	p := NewInputPrompt("")

	renderer := cfg.Renderer
	if renderer == nil {
		renderer = &DefaultSecretRenderer{}
	}

	return &SecretPrompt{
		Label:    cfg.Label,
		Confirm:  cfg.Confirm,
		Mask:     cfg.Mask,
		prompt:   p,
		renderer: renderer,
	}, nil
}

// Secret is a convenience function that creates a SecretPrompt.
func Secret(opts ...any) (*SecretPrompt, error) {
	cfg := share.OverloadWithOptions[SecretConfig](opts, DefaultSecretConfig())
	return NewSecretPrompt(cfg)
}

func (s *SecretPrompt) Done() <-chan string       { return s.prompt.Done }
func (s *SecretPrompt) Canceled() <-chan struct{} { return s.prompt.Canceled }

// Render uses the configured renderer.
func (s *SecretPrompt) Render() []byte { return s.renderer.Render(s) }

// OnKey delegates to the internal prompt. Confirmation logic can be added externally.
func (s *SecretPrompt) OnKey(key runfx.Key) bool {
	return s.prompt.OnKey(key)
}

func (s *SecretPrompt) OnResize(cols, rows int) {}
