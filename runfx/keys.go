package runfx

import (
	"fmt"
)

// Key represents a keyboard input event.
type Key int

// Key constants for common keyboard inputs.
const (
	KeyUnknown Key = iota

	// Special keys
	KeyEnter
	KeyEscape
	KeyBackspace
	KeyTab
	KeySpace
	KeyDelete

	// Arrow keys
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight

	// Letter keys (WASD for navigation)
	KeyA
	KeyS
	KeyD
	KeyW

	// Letter keys (Y/N for confirmation)
	KeyY
	KeyN

	// Number keys
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9

	// Common shortcuts
	KeyCtrlC
	KeyCtrlD
	KeyQ
)

// String returns a human-readable representation of the key.
func (k Key) String() string {
	switch k {
	case KeyEnter:
		return "Enter"
	case KeyEscape:
		return "Escape"
	case KeyBackspace:
		return "Backspace"
	case KeyTab:
		return "Tab"
	case KeySpace:
		return "Space"
	case KeyDelete:
		return "Delete"
	case KeyArrowUp:
		return "ArrowUp"
	case KeyArrowDown:
		return "ArrowDown"
	case KeyArrowLeft:
		return "ArrowLeft"
	case KeyArrowRight:
		return "ArrowRight"
	case KeyA:
		return "A"
	case KeyS:
		return "S"
	case KeyD:
		return "D"
	case KeyW:
		return "W"
	case KeyY:
		return "Y"
	case KeyN:
		return "N"
	case KeyQ:
		return "Q"
	case KeyCtrlC:
		return "Ctrl+C"
	case KeyCtrlD:
		return "Ctrl+D"
	case Key0, Key1, Key2, Key3, Key4, Key5, Key6, Key7, Key8, Key9:
		return fmt.Sprintf("%d", int(k-Key0))
	default:
		return "Unknown"
	}
}

// IsArrow returns true if the key is an arrow key.
func (k Key) IsArrow() bool {
	return k >= KeyArrowUp && k <= KeyArrowRight
}

// IsWASD returns true if the key is a WASD navigation key.
func (k Key) IsWASD() bool {
	return k == KeyW || k == KeyA || k == KeyS || k == KeyD
}

// IsNavigation returns true if the key is used for navigation (arrows or WASD).
func (k Key) IsNavigation() bool {
	return k.IsArrow() || k.IsWASD()
}

// IsNumber returns true if the key is a number key.
func (k Key) IsNumber() bool {
	return k >= Key0 && k <= Key9
}

// ToNumber converts a number key to its integer value.
// Returns -1 if the key is not a number.
func (k Key) ToNumber() int {
	if k.IsNumber() {
		return int(k - Key0)
	}
	return -1
}
