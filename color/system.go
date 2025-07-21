package color

// Clean Color Encoding System
// This provides the elegant API: color.ANSI.Blue, color.Material.Blue, color.Blue

// EncodingSystem provides colors in a specific encoding mode
type EncodingSystem struct {
	// Basic colors
	Black   Color
	Red     Color
	Green   Color
	Yellow  Color
	Blue    Color
	Magenta Color
	Cyan    Color
	White   Color

	// Bright colors
	BrightBlack   Color
	BrightRed     Color
	BrightGreen   Color
	BrightYellow  Color
	BrightBlue    Color
	BrightMagenta Color
	BrightCyan    Color
	BrightWhite   Color
}

// ThemeSystem provides colors in a specific theme
type ThemeSystem struct {
	// Basic colors
	Black   Color
	Red     Color
	Green   Color
	Yellow  Color
	Blue    Color
	Magenta Color
	Cyan    Color
	White   Color

	// Extended theme colors
	Purple Color
	Orange Color
	Pink   Color
	Teal   Color
	Lime   Color
	Indigo Color
}

// Global encoding systems - clean names!
var (
	// Encoding-specific access: color.ANSI.Blue
	ANSI      *EncodingSystem
	TrueColor *EncodingSystem
	Color256  *EncodingSystem

	// Theme-specific access: color.Material.Blue
	Material *ThemeSystem
	Dracula  *ThemeSystem
	Nord     *ThemeSystem
	GitHub   *ThemeSystem
)

// Default colors - clean names that follow active theme
var (
	// Basic colors - these follow the active theme and encoding
	Black   Color
	Red     Color
	Green   Color
	Yellow  Color
	Blue    Color
	Magenta Color
	Cyan    Color
	White   Color

	// Extended colors
	Purple Color
	Orange Color
	Pink   Color
	Teal   Color
	Lime   Color
	Indigo Color
)

// Global configuration
var (
	currentEncoding = ModeANSI
	currentTheme    = "material"
)

func init() {
	initializeCleanColorSystems()
}

// initializeCleanColorSystems sets up all encoding and theme systems
func initializeCleanColorSystems() {
	createCleanEncodingSystems()
	createCleanThemeSystems()
	updateCleanDefaultColors()
}

// createCleanEncodingSystems initializes encoding-specific color systems
func createCleanEncodingSystems() {
	// ANSI encoding system (16 colors)
	ANSI = &EncodingSystem{
		Black:   NewColor(ColorConfig{ANSI: 0, Name: "ansi_black"}),
		Red:     NewColor(ColorConfig{ANSI: 1, Name: "ansi_red"}),
		Green:   NewColor(ColorConfig{ANSI: 2, Name: "ansi_green"}),
		Yellow:  NewColor(ColorConfig{ANSI: 3, Name: "ansi_yellow"}),
		Blue:    NewColor(ColorConfig{ANSI: 4, Name: "ansi_blue"}),
		Magenta: NewColor(ColorConfig{ANSI: 5, Name: "ansi_magenta"}),
		Cyan:    NewColor(ColorConfig{ANSI: 6, Name: "ansi_cyan"}),
		White:   NewColor(ColorConfig{ANSI: 7, Name: "ansi_white"}),

		BrightBlack:   NewColor(ColorConfig{ANSI: 8, Name: "ansi_bright_black"}),
		BrightRed:     NewColor(ColorConfig{ANSI: 9, Name: "ansi_bright_red"}),
		BrightGreen:   NewColor(ColorConfig{ANSI: 10, Name: "ansi_bright_green"}),
		BrightYellow:  NewColor(ColorConfig{ANSI: 11, Name: "ansi_bright_yellow"}),
		BrightBlue:    NewColor(ColorConfig{ANSI: 12, Name: "ansi_bright_blue"}),
		BrightMagenta: NewColor(ColorConfig{ANSI: 13, Name: "ansi_bright_magenta"}),
		BrightCyan:    NewColor(ColorConfig{ANSI: 14, Name: "ansi_bright_cyan"}),
		BrightWhite:   NewColor(ColorConfig{ANSI: 15, Name: "ansi_bright_white"}),
	}

	// TrueColor encoding system (24-bit RGB)
	TrueColor = &EncodingSystem{
		Black:   NewColor(ColorConfig{R: 0, G: 0, B: 0, Name: "true_black"}),
		Red:     NewColor(ColorConfig{R: 255, G: 0, B: 0, Name: "true_red"}),
		Green:   NewColor(ColorConfig{R: 0, G: 255, B: 0, Name: "true_green"}),
		Yellow:  NewColor(ColorConfig{R: 255, G: 255, B: 0, Name: "true_yellow"}),
		Blue:    NewColor(ColorConfig{R: 0, G: 0, B: 255, Name: "true_blue"}),
		Magenta: NewColor(ColorConfig{R: 255, G: 0, B: 255, Name: "true_magenta"}),
		Cyan:    NewColor(ColorConfig{R: 0, G: 255, B: 255, Name: "true_cyan"}),
		White:   NewColor(ColorConfig{R: 255, G: 255, B: 255, Name: "true_white"}),

		BrightBlack:   NewColor(ColorConfig{R: 128, G: 128, B: 128, Name: "true_bright_black"}),
		BrightRed:     NewColor(ColorConfig{R: 255, G: 128, B: 128, Name: "true_bright_red"}),
		BrightGreen:   NewColor(ColorConfig{R: 128, G: 255, B: 128, Name: "true_bright_green"}),
		BrightYellow:  NewColor(ColorConfig{R: 255, G: 255, B: 128, Name: "true_bright_yellow"}),
		BrightBlue:    NewColor(ColorConfig{R: 128, G: 128, B: 255, Name: "true_bright_blue"}),
		BrightMagenta: NewColor(ColorConfig{R: 255, G: 128, B: 255, Name: "true_bright_magenta"}),
		BrightCyan:    NewColor(ColorConfig{R: 128, G: 255, B: 255, Name: "true_bright_cyan"}),
		BrightWhite:   NewColor(ColorConfig{R: 255, G: 255, B: 255, Name: "true_bright_white"}),
	}

	// 256-color encoding system
	Color256 = &EncodingSystem{
		Black:   NewColor(ColorConfig{Color256: 0, Name: "256_black"}),
		Red:     NewColor(ColorConfig{Color256: 1, Name: "256_red"}),
		Green:   NewColor(ColorConfig{Color256: 2, Name: "256_green"}),
		Yellow:  NewColor(ColorConfig{Color256: 3, Name: "256_yellow"}),
		Blue:    NewColor(ColorConfig{Color256: 4, Name: "256_blue"}),
		Magenta: NewColor(ColorConfig{Color256: 5, Name: "256_magenta"}),
		Cyan:    NewColor(ColorConfig{Color256: 6, Name: "256_cyan"}),
		White:   NewColor(ColorConfig{Color256: 7, Name: "256_white"}),

		BrightBlack:   NewColor(ColorConfig{Color256: 8, Name: "256_bright_black"}),
		BrightRed:     NewColor(ColorConfig{Color256: 9, Name: "256_bright_red"}),
		BrightGreen:   NewColor(ColorConfig{Color256: 10, Name: "256_bright_green"}),
		BrightYellow:  NewColor(ColorConfig{Color256: 11, Name: "256_bright_yellow"}),
		BrightBlue:    NewColor(ColorConfig{Color256: 12, Name: "256_bright_blue"}),
		BrightMagenta: NewColor(ColorConfig{Color256: 13, Name: "256_bright_magenta"}),
		BrightCyan:    NewColor(ColorConfig{Color256: 14, Name: "256_bright_cyan"}),
		BrightWhite:   NewColor(ColorConfig{Color256: 15, Name: "256_bright_white"}),
	}
}

// createCleanThemeSystems initializes theme-specific color systems
func createCleanThemeSystems() {
	// Material Design theme
	Material = &ThemeSystem{
		Black:   NewHex("#424242").WithName("material_black"), // Material Dark Gray
		Red:     MaterialRed,
		Green:   MaterialGreen,
		Yellow:  MaterialYellow,
		Blue:    MaterialBlue,
		Magenta: MaterialPurple,
		Cyan:    MaterialCyan,
		White:   NewHex("#FFFFFF").WithName("material_white"),
		Purple:  MaterialPurple,
		Orange:  MaterialOrange,
		Pink:    MaterialPink,
		Teal:    MaterialTeal,
		Lime:    MaterialLime,
		Indigo:  MaterialIndigo,
	}

	// Dracula theme
	Dracula = &ThemeSystem{
		Black:   NewHex("#282A36").WithName("dracula_black"), // Dracula background
		Red:     DraculaRed,
		Green:   DraculaGreen,
		Yellow:  DraculaYellow,
		Blue:    NewHex("#6272A4").WithName("dracula_blue"), // Dracula comment
		Magenta: DraculaPurple,
		Cyan:    DraculaCyan,
		White:   NewHex("#F8F8F2").WithName("dracula_white"), // Dracula foreground
		Purple:  DraculaPurple,
		Orange:  DraculaOrange,
		Pink:    DraculaPink,
		Teal:    DraculaCyan,
		Lime:    DraculaGreen,
		Indigo:  DraculaPurple,
	}

	// Nord theme
	Nord = &ThemeSystem{
		Black:   NewHex("#2E3440").WithName("nord_black"), // Nord polar night
		Red:     NordRed,
		Green:   NordGreen,
		Yellow:  NordYellow,
		Blue:    NordBlue,
		Magenta: NordPurple,
		Cyan:    NordCyan,
		White:   NewHex("#ECEFF4").WithName("nord_white"), // Nord snow storm
		Purple:  NordPurple,
		Orange:  NordOrange,
		Pink:    NordPurple, // Nord doesn't have distinct pink
		Teal:    NordCyan,
		Lime:    NordGreen,
		Indigo:  NordBlue,
	}

	// GitHub theme
	GitHub = &ThemeSystem{
		Black:   NewHex("#24292e").WithName("github_black"),
		Red:     GithubRedLight,
		Green:   GithubGreenLight,
		Yellow:  NewHex("#FBBF40").WithName("github_yellow"),
		Blue:    GithubBlueLight,
		Magenta: NewHex("#B392F0").WithName("github_magenta"),
		Cyan:    GithubBlueLight,
		White:   NewHex("#FFFFFF").WithName("github_white"),
		Purple:  NewHex("#B392F0").WithName("github_purple"),
		Orange:  GithubOrangeLight,
		Pink:    NewHex("#F97583").WithName("github_pink"),
		Teal:    GithubBlueLight,
		Lime:    GithubGreenLight,
		Indigo:  GithubBlueLight,
	}
}

// updateCleanDefaultColors sets the global default colors based on current theme
func updateCleanDefaultColors() {
	var source *ThemeSystem

	switch currentTheme {
	case "material":
		source = Material
	case "dracula":
		source = Dracula
	case "nord":
		source = Nord
	case "github":
		source = GitHub
	default:
		source = Material
	}

	// Update default colors to follow the active theme
	Black = source.Black
	Red = source.Red
	Green = source.Green
	Yellow = source.Yellow
	Blue = source.Blue
	Magenta = source.Magenta
	Cyan = source.Cyan
	White = source.White
	Purple = source.Purple
	Orange = source.Orange
	Pink = source.Pink
	Teal = source.Teal
	Lime = source.Lime
	Indigo = source.Indigo
}

// SetDefaultEncoding changes the default encoding for color rendering
func SetDefaultEncoding(mode Mode) {
	currentEncoding = mode
}

// SetDefaultTheme changes the default theme for colors
func SetDefaultTheme(theme string) {
	currentTheme = theme
	updateCleanDefaultColors()
}

// GetDefaultEncoding returns the current default encoding
func GetDefaultEncoding() Mode {
	return currentEncoding
}

// GetDefaultTheme returns the current default theme
func GetDefaultTheme() string {
	return currentTheme
}

// Convenience functions to switch themes quickly
func UseMaterial() { SetDefaultTheme("material") }
func UseDracula()  { SetDefaultTheme("dracula") }
func UseNord()     { SetDefaultTheme("nord") }
func UseGitHub()   { SetDefaultTheme("github") }
