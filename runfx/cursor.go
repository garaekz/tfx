package runfx

// CursorManager provides advanced cursor control for visuals
import "fmt"

type CursorManager struct {
	visible bool
	posX    int
	posY    int
	prevX   int
	prevY   int
}

// Hide hides the cursor using ANSI escape code
func (c *CursorManager) Hide() {
	fmt.Print("\033[?25l")
	c.visible = false
}

// Show shows the cursor using ANSI escape code
func (c *CursorManager) Show() {
	fmt.Print("\033[?25h")
	c.visible = true
}

// MoveTo moves the cursor to (x, y) using ANSI escape code
func (c *CursorManager) MoveTo(x, y int) {
	c.prevX = c.posX
	c.prevY = c.posY
	c.posX = x
	c.posY = y
	// ANSI is 1-based: \033[{row};{col}H
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

// Restore moves the cursor back to previous position
func (c *CursorManager) Restore() {
	fmt.Printf("\033[%d;%dH", c.prevY+1, c.prevX+1)
	c.posX, c.posY = c.prevX, c.prevY
}
