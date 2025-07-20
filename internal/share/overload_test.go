package share

import (
	"testing"
)

type Config struct {
	Name string
}

func TestOverload(t *testing.T) {
	def := Config{Name: "default"}

	t.Run("no args returns default", func(t *testing.T) {
		got := Overload([]any{}, def)
		if got != def {
			t.Errorf("expected default, got %+v", got)
		}
	})

	t.Run("direct value overrides", func(t *testing.T) {
		custom := Config{Name: "custom"}
		got := Overload([]any{custom}, def)
		if got != custom {
			t.Errorf("expected custom, got %+v", got)
		}
	})

	t.Run("pointer value overrides", func(t *testing.T) {
		custom := Config{Name: "ptr"}
		got := Overload([]any{&custom}, def)
		if got != custom {
			t.Errorf("expected ptr, got %+v", got)
		}
	})
}
func TestOverloadWithOptions(t *testing.T) {
	def := Config{Name: "default"}

	t.Run("no args returns default", func(t *testing.T) {
		got := OverloadWithOptions([]any{}, def)
		if got != def {
			t.Errorf("expected default, got %+v", got)
		}
	})

	t.Run("direct value overrides with options", func(t *testing.T) {
		custom := Config{Name: "custom"}
		got := OverloadWithOptions([]any{custom}, def, func(c *Config) {
			c.Name = "overridden"
		})
		if got.Name != "overridden" {
			t.Errorf("expected overridden, got %+v", got)
		}
	})

	t.Run("pointer value overrides with options", func(t *testing.T) {
		custom := Config{Name: "ptr"}
		got := OverloadWithOptions([]any{&custom}, def, func(c *Config) {
			c.Name = "ptr-overridden"
		})
		if got.Name != "ptr-overridden" {
			t.Errorf("expected ptr-overridden, got %+v", got)
		}
	})
}
