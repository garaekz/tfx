package progress

import (
	"io"

	"github.com/garaekz/tfx/internal/share"
)

// --- 1. Configuración ---

// SpinnerConfig contiene la configuración declarativa para un componente Spinner.
type SpinnerConfig struct {
	Total   int
	Label   string
	Width   int
	Theme   ProgressTheme
	Style   ProgressStyle
	Effect  ProgressEffect
	Writer  io.Writer // Usado para la detección de TTY, no para escritura directa.
	ShowETA bool
}

// DefaultSpinnerConfig devuelve la configuración por defecto.
func DefaultSpinnerConfig() SpinnerConfig {
	return SpinnerConfig{
		Total: 100,
		Label: "Spinner",
		Width: 40,
		// Asigna valores por defecto para Theme, Style, Effect...
	}
}

// --- 2. API Multipath ---

// Start es la función de conveniencia de alto nivel (vías Express e Instantiated).
// Crea y devuelve un componente Spinner listo para ser montado en un runfx.Loop.
func StartSpinner(opts ...any) *Spinner {
	cfg := share.OverloadWithOptions[SpinnerConfig](opts, DefaultSpinnerConfig())
	return newSpinner(cfg)
}

// --- 3. Vía DSL (Builder) ---

// SpinnerBuilder proporciona la vía DSL para una configuración fluida.
type SpinnerBuilder struct {
	config SpinnerConfig
}

// NewSpinnerBuilder es el punto de entrada para la vía DSL.
func NewSpinnerBuilder() *SpinnerBuilder {
	return &SpinnerBuilder{
		config: DefaultSpinnerConfig(),
	}
}

// Total establece el valor total de la barra de progreso.
func (b *SpinnerBuilder) Total(total int) *SpinnerBuilder {
	b.config.Total = total
	return b
}

// Label establece la etiqueta de la barra de progreso.
func (b *SpinnerBuilder) Label(label string) *SpinnerBuilder {
	b.config.Label = label
	return b
}

// Width establece el ancho de la barra de progreso.
func (b *SpinnerBuilder) Width(width int) *SpinnerBuilder {
	b.config.Width = width
	return b
}

// Build construye el componente Spinner con la configuración proporcionada.
func (b *SpinnerBuilder) Build() *Spinner {
	return newSpinner(b.config)
}
