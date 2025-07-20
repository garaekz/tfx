package color

// ANSI Escape Codes and Control Sequences
const (
	// Reset code
	Reset = "\033[0m"

	// Text Attributes
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Strike    = "\033[9m"

	// Basic ANSI Foreground Colors (30-37)
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright ANSI Foreground Colors (90-97)
	BrightBlack   = "\033[90m" // Dark Gray
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Basic ANSI Background Colors (40-47)
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Bright ANSI Background Colors (100-107)
	BgBrightBlack   = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"
)

// Common terminal escape sequences for cursor control and screen manipulation
const (
	// Cursor Control
	CursorUp    = "\033[A"
	CursorDown  = "\033[B"
	CursorRight = "\033[C"
	CursorLeft  = "\033[D"
	CursorHome  = "\033[H"

	// Screen Control
	ClearScreen     = "\033[2J"
	ClearLine       = "\033[2K"
	ClearToEnd      = "\033[0K"
	ClearToStart    = "\033[1K"
	SaveCursor      = "\033[s"
	RestoreCursor   = "\033[u"
	HideCursor      = "\033[?25l"
	ShowCursor      = "\033[?25h"
)