package runfx

import (
	"testing"
)

func TestCursorManagerBasic(t *testing.T) {
	cm := &CursorManager{}

	// Test initial state
	if cm.visible {
		t.Error("Cursor should not be visible initially")
	}

	// Test hide/show - these mainly test that methods don't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Cursor operations panicked: %v", r)
		}
	}()

	cm.Hide()
	if cm.visible {
		t.Error("Cursor should be hidden after Hide()")
	}

	cm.Show()
	if !cm.visible {
		t.Error("Cursor should be visible after Show()")
	}
}

func TestCursorManagerMovement(t *testing.T) {
	cm := &CursorManager{}

	// Test movement
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Cursor movement panicked: %v", r)
		}
	}()

	cm.MoveTo(10, 20)
	if cm.posX != 10 || cm.posY != 20 {
		t.Errorf("Cursor position not updated correctly: got (%d, %d), want (10, 20)", cm.posX, cm.posY)
	}

	// Previous position should be stored
	if cm.prevX != 0 || cm.prevY != 0 {
		t.Errorf("Previous position not stored correctly: got (%d, %d), want (0, 0)", cm.prevX, cm.prevY)
	}
}

func TestCursorManagerRestore(t *testing.T) {
	cm := &CursorManager{}

	// Move to initial position
	cm.MoveTo(5, 10)

	// Move to new position
	cm.MoveTo(15, 25)

	if cm.posX != 15 || cm.posY != 25 {
		t.Errorf("Current position incorrect: got (%d, %d), want (15, 25)", cm.posX, cm.posY)
	}

	if cm.prevX != 5 || cm.prevY != 10 {
		t.Errorf("Previous position incorrect: got (%d, %d), want (5, 10)", cm.prevX, cm.prevY)
	}

	// Test restore
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Cursor restore panicked: %v", r)
		}
	}()

	cm.Restore()

	// Position should be swapped back
	if cm.posX != 5 || cm.posY != 10 {
		t.Errorf("Position after restore incorrect: got (%d, %d), want (5, 10)", cm.posX, cm.posY)
	}
}

func TestCursorManagerMultipleOperations(t *testing.T) {
	cm := &CursorManager{}

	// Sequence of operations
	cm.Hide()
	cm.MoveTo(1, 1)
	cm.MoveTo(2, 2)
	cm.MoveTo(3, 3)
	cm.Restore() // Should go back to (2, 2)
	cm.Show()

	if cm.posX != 2 || cm.posY != 2 {
		t.Errorf("Final position incorrect: got (%d, %d), want (2, 2)", cm.posX, cm.posY)
	}

	if !cm.visible {
		t.Error("Cursor should be visible at end")
	}
}
