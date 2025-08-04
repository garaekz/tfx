package progress

import (
	"sync"
	"time"

	"github.com/garaekz/tfx/terminal"
)

// --- 1. El Componente de Alto Nivel: Spinner ---

// Spinner es un componente de UI que muestra una barra de progreso.
// Es un componente pasivo diseñado para ser montado en un runfx.Loop.
type Spinner struct {
	// Estado
	total     int
	current   int
	label     string
	startTime time.Time
	isStarted bool

	// Configuración de Apariencia
	width    int
	theme    SpinnerTheme
	style    SpinnerStyle
	effect   SpinnerEffect
	detector *terminal.Detector
	ShowETA  bool

	mu sync.Mutex
}

// newSpinner es el constructor interno que ensambla el componente.
// Es llamado por la API pública (Start, NewSpinnerBuilder).
func newSpinner(cfg SpinnerConfig) *Spinner {
	return &Spinner{
		total:    cfg.Total,
		label:    cfg.Label,
		width:    cfg.Width,
		theme:    cfg.Theme,
		style:    cfg.Style,
		effect:   cfg.Effect,
		detector: terminal.NewDetector(cfg.Writer), // Asume que el detector se basa en el writer.
		ShowETA:  cfg.ShowETA,
	}
}

// --- 2. Implementación de Interfaces de RunFX ---

// Render implementa la interfaz runfx.Visual.
// Es llamado por el runfx.Loop en cada ciclo de renderizado.
func (p *Spinner) Render() []byte {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		return nil // No renderiza nada si no ha comenzado.
	}

	// Llama a la lógica de renderizado existente para generar el string.
	rendered := RenderBar(p, p.detector)
	return []byte(rendered)
}

// Tick implementa la interfaz runfx.Updatable.
// Es llamado por el runfx.Loop en cada tick, útil para efectos animados.
func (p *Spinner) Tick(now time.Time) {
	// La lógica para efectos de animación (pulsos, spinners) iría aquí.
	// Por ahora, no hace nada, pero el contrato está cumplido.
}

// OnResize implementa la interfaz runfx.Visual.
func (p *Spinner) OnResize(cols, rows int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// Lógica opcional para ajustar el ancho de la barra al tamaño del terminal.
	if cols > 0 && cols < p.width+20 {
		p.width = max(cols-20, 10)
	}
}

// --- 3. Métodos para Manipular el Estado ---

// Set actualiza el valor de progreso actual.
func (p *Spinner) Set(current int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		p.isStarted = true
		p.startTime = time.Now()
	}
	p.current = min(current, p.total)
}

// Add incrementa el progreso en una cantidad dada.
func (p *Spinner) Add(amount int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isStarted {
		p.isStarted = true
		p.startTime = time.Now()
	}
	p.current = min(p.current+amount, p.total)
}

// SetLabel actualiza la etiqueta de la barra de progreso.
func (p *Spinner) SetLabel(label string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.label = label
}

// Finish establece el progreso al 100%.
// En esta arquitectura, no detiene el loop ni se desmonta a sí mismo.
func (p *Spinner) Finish() {
	p.Set(p.total)
}
