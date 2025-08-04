package progress

import (
	"io"

	"github.com/garaekz/tfx/internal/share"
)

// --- 1. Configuración ---

// ProgressConfig contiene la configuración declarativa para un componente Progress.
type ProgressConfig struct {
	Total   int
	Label   string
	Width   int
	Theme   ProgressTheme
	Style   ProgressStyle
	Effect  ProgressEffect
	Writer  io.Writer // Usado para la detección de TTY, no para escritura directa.
	ShowETA bool
}

// DefaultProgressConfig devuelve la configuración por defecto.
func DefaultProgressConfig() ProgressConfig {
	return ProgressConfig{
		Total: 100,
		Label: "Progress",
		Width: 40,
		// Asigna valores por defecto para Theme, Style, Effect...
	}
}

// --- 2. API Multipath ---

// Start es la función de conveniencia de alto nivel (vías Express e Instantiated).
// Crea y devuelve un componente Progress listo para ser montado en un runfx.Loop.
func Start(opts ...any) *Progress {
	cfg := share.OverloadWithOptions[ProgressConfig](opts, DefaultProgressConfig())
	return newProgress(cfg)
}

// --- 3. Vía DSL (Builder) ---

// ProgressBuilder proporciona la vía DSL para una configuración fluida.
type ProgressBuilder struct {
	config ProgressConfig
}

// NewProgressBuilder es el punto de entrada para la vía DSL.
func NewProgressBuilder() *ProgressBuilder {
	return &ProgressBuilder{
		config: DefaultProgressConfig(),
	}
}

// Total establece el valor total de la barra de progreso.
func (b *ProgressBuilder) Total(total int) *ProgressBuilder {
	b.config.Total = total
	return b
}

// Label establece la etiqueta de la barra de progreso.
func (b *ProgressBuilder) Label(label string) *ProgressBuilder {
	b.config.Label = label
	return b
}

// Width establece el ancho de la barra de progreso.
func (b *ProgressBuilder) Width(width int) *ProgressBuilder {
	b.config.Width = width
	return b
}

// Build construye el componente Progress con la configuración proporcionada.
func (b *ProgressBuilder) Build() *Progress {
	return newProgress(b.config)
}
