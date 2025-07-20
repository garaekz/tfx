package color

import "github.com/garaekz/tfx/internal/share"

// Re-export commonly used ANSI sequences for backward compatibility
var (
	Reset     = share.Reset
	Bold      = share.Bold
	Dim       = share.Dim
	Italic    = share.Italic
	Underline = share.Underline
	Blink     = share.Blink
	Reverse   = share.Reverse
	Strike    = share.Strike
)

// ANSISeq provides access to raw ANSI escape sequences
// Use this when you need direct terminal control
var ANSISeq = share.ANSISeq

// Backward compatibility - these will be deprecated
// Use color.ANSI.Blue instead of these constants
// These old ANSI escape sequence constants have been moved to share.ANSISeq
// Example: share.ANSISeq.Blue() instead of Blue
//
// For the new clean color API, use:
// - color.ANSI.Blue (for ANSI encoding)
// - color.Material.Blue (for Material theme)  
// - color.Blue (for default colors following active theme)