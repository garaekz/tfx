package runfx

import "sync"

// Multiplexer manages thread-safe visual mounting/unmounting
type Multiplexer struct {
	visuals map[string]Visual
	mu      sync.Mutex
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{visuals: make(map[string]Visual)}
}

// Mount adds a visual, reusing gaps if possible
func (m *Multiplexer) Mount(name string, v Visual) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Reuse gap if available
	if _, exists := m.visuals[name]; !exists {
		m.visuals[name] = v
	}
}

// Unmount removes a visual and leaves a gap for reuse
func (m *Multiplexer) Unmount(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Set to nil to create a reusable gap instead of deleting
	if _, exists := m.visuals[name]; exists {
		m.visuals[name] = nil
	}
}

// ListVisuals returns the names of all mounted visuals
func (m *Multiplexer) ListVisuals() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]string, 0, len(m.visuals))
	for k, v := range m.visuals {
		if v != nil { // Only include non-nil visuals
			keys = append(keys, k)
		}
	}
	return keys
}

// GetVisual returns a visual by name
func (m *Multiplexer) GetVisual(name string) (Visual, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.visuals[name]
	// Return false if the visual is nil (gap)
	if v == nil {
		return nil, false
	}
	return v, ok
}

// ReuseGap finds and reuses a gap for a new visual
func (m *Multiplexer) ReuseGap(v Visual) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, vis := range m.visuals {
		if vis == nil {
			m.visuals[k] = v
			return k, true
		}
	}
	return "", false
}
