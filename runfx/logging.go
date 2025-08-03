package runfx

import (
	"fmt"
	"os"
)

var debugMode = false

// EnableDebug turns on verbose terminal state logging
func EnableDebug() {
	debugMode = true
}

// DebugLog prints debug info if debug mode is enabled
func DebugLog(format string, args ...interface{}) {
	if debugMode {
		fmt.Fprintf(os.Stderr, "[RunFX DEBUG] "+format+"\n", args...)
	}
}

// LogMount logs visual mounting
func LogMount(name string) {
	DebugLog("Mounted visual: %s", name)
}

// LogUnmount logs visual unmounting
func LogUnmount(name string) {
	DebugLog("Unmounted visual: %s", name)
}

// LogRender logs rendering events
func LogRender(name string) {
	DebugLog("Rendered visual: %s", name)
}
