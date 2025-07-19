package color

// Semantic Color Palette - Meaningful names for better UX
var (
	// Status Colors
	Success = BrightGreen
	Error   = BrightRed
	Warning = BrightYellow
	Info    = BrightCyan
	Debug   = BrightMagenta

	// Semantic Colors
	Primary   = BrightBlue
	Secondary = BrightMagenta
	Accent    = BrightCyan
	Muted     = BrightBlack // Dark gray
	Highlight = BrightWhite

	// UI Element Colors
	Border   = BrightBlack
	Text     = White
	Subtitle = BrightBlack
	Link     = BrightBlue
	Code     = BrightYellow

	// Alert Colors
	Critical = BrightRed
	Caution  = BrightYellow
	Notice   = BrightCyan
	Tip      = BrightGreen

	// Progress Colors
	ProgressComplete   = BrightGreen
	ProgressIncomplete = BrightBlack
	ProgressBar        = BrightBlue

	// Background Variants
	BgError   = BgRed
	BgSuccess = BgGreen
	BgWarning = BgYellow
	BgInfo    = BgBlue
	BgDebug   = BgMagenta
)

// Brand/Theme Colors - Inspired by popular color schemes
var (
	// GitHub-inspired
	GithubGreen  = BrightGreen
	GithubRed    = BrightRed
	GithubBlue   = BrightBlue
	GithubOrange = BrightYellow

	// Terminal-inspired
	TerminalGreen = Green
	TerminalRed   = Red
	TerminalBlue  = Blue
	TerminalGray  = BrightBlack

	// Retro Colors
	RetroAmber = BrightYellow
	RetroGreen = Green
	RetroCyan  = Cyan

	// Neon Colors (using bright variants)
	NeonPink   = BrightMagenta
	NeonGreen  = BrightGreen
	NeonBlue   = BrightCyan
	NeonYellow = BrightYellow
)

// Log Level Color Mapping
var LogLevelColors = map[string]string{
	"TRACE": BrightBlack,
	"DEBUG": BrightMagenta,
	"INFO":  BrightCyan,
	"WARN":  BrightYellow,
	"ERROR": BrightRed,
	"FATAL": Bold + BrightRed,
	"PANIC": Bold + BrightRed + BgRed,
}

// Badge Theme Colors
var BadgeThemes = map[string]BadgeTheme{
	"default": {
		Success: Success,
		Error:   Error,
		Warning: Warning,
		Info:    Info,
		Debug:   Debug,
	},
	"github": {
		Success: GithubGreen,
		Error:   GithubRed,
		Warning: GithubOrange,
		Info:    GithubBlue,
		Debug:   BrightBlack,
	},
	"terminal": {
		Success: TerminalGreen,
		Error:   TerminalRed,
		Warning: BrightYellow,
		Info:    TerminalBlue,
		Debug:   TerminalGray,
	},
	"neon": {
		Success: NeonGreen,
		Error:   NeonPink,
		Warning: NeonYellow,
		Info:    NeonBlue,
		Debug:   BrightBlack,
	},
}

type BadgeTheme struct {
	Success string
	Error   string
	Warning string
	Info    string
	Debug   string
}
