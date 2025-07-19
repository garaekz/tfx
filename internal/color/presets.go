package color

// Predefined colors using the Color struct system
var (
	// Material Design Colors
	MaterialRed        = NewHex("#F44336")
	MaterialPink       = NewHex("#E91E63")
	MaterialPurple     = NewHex("#9C27B0")
	MaterialDeepPurple = NewHex("#673AB7")
	MaterialIndigo     = NewHex("#3F51B5")
	MaterialBlue       = NewHex("#2196F3")
	MaterialLightBlue  = NewHex("#03A9F4")
	MaterialCyan       = NewHex("#00BCD4")
	MaterialTeal       = NewHex("#009688")
	MaterialGreen      = NewHex("#4CAF50")
	MaterialLightGreen = NewHex("#8BC34A")
	MaterialLime       = NewHex("#CDDC39")
	MaterialYellow     = NewHex("#FFEB3B")
	MaterialAmber      = NewHex("#FFC107")
	MaterialOrange     = NewHex("#FF9800")
	MaterialDeepOrange = NewHex("#FF5722")

	// Tailwind CSS Colors (popular selections)
	TailwindRed     = NewHex("#EF4444")
	TailwindOrange  = NewHex("#F97316")
	TailwindAmber   = NewHex("#F59E0B")
	TailwindYellow  = NewHex("#EAB308")
	TailwindLime    = NewHex("#84CC16")
	TailwindGreen   = NewHex("#22C55E")
	TailwindEmerald = NewHex("#10B981")
	TailwindTeal    = NewHex("#14B8A6")
	TailwindCyan    = NewHex("#06B6D4")
	TailwindSky     = NewHex("#0EA5E9")
	TailwindBlue    = NewHex("#3B82F6")
	TailwindIndigo  = NewHex("#6366F1")
	TailwindViolet  = NewHex("#8B5CF6")
	TailwindPurple  = NewHex("#A855F7")
	TailwindFuchsia = NewHex("#D946EF")
	TailwindPink    = NewHex("#EC4899")
	TailwindRose    = NewHex("#F43F5E")

	// GitHub Colors
	GithubGreenLight  = NewHex("#28A745")
	GithubGreenDark   = NewHex("#22863A")
	GithubRedLight    = NewHex("#D73A49")
	GithubRedDark     = NewHex("#B31D28")
	GithubBlueLight   = NewHex("#0366D6")
	GithubBlueDark    = NewHex("#005CC5")
	GithubOrangeLight = NewHex("#F66A0A")
	GithubOrangeDark  = NewHex("#E36209")
	GithubPurple      = NewHex("#6F42C1")
	GithubYellow      = NewHex("#FFD33D")

	// VS Code Theme Colors
	VSCodeBlue   = NewHex("#007ACC")
	VSCodeGreen  = NewHex("#608B4E")
	VSCodeRed    = NewHex("#F44747")
	VSCodeOrange = NewHex("#FF8C00")
	VSCodePurple = NewHex("#C586C0")
	VSCodeYellow = NewHex("#FFCD3C")
	VSCodeCyan   = NewHex("#4EC9B0")

	// Dracula Theme
	DraculaPurple = NewHex("#BD93F9")
	DraculaPink   = NewHex("#FF79C6")
	DraculaGreen  = NewHex("#50FA7B")
	DraculaOrange = NewHex("#FFB86C")
	DraculaRed    = NewHex("#FF5555")
	DraculaYellow = NewHex("#F1FA8C")
	DraculaCyan   = NewHex("#8BE9FD")

	// Nord Theme
	NordBlue      = NewHex("#5E81AC")
	NordLightBlue = NewHex("#81A1C1")
	NordCyan      = NewHex("#88C0D0")
	NordGreen     = NewHex("#A3BE8C")
	NordYellow    = NewHex("#EBCB8B")
	NordOrange    = NewHex("#D08770")
	NordRed       = NewHex("#BF616A")
	NordPurple    = NewHex("#B48EAD")

	// Semantic Colors with multiple shades
	SuccessLight = NewHex("#10B981") // Emerald 500
	SuccessDark  = NewHex("#047857") // Emerald 700
	ErrorLight   = NewHex("#EF4444") // Red 500
	ErrorDark    = NewHex("#B91C1C") // Red 700
	WarningLight = NewHex("#F59E0B") // Amber 500
	WarningDark  = NewHex("#D97706") // Amber 600
	InfoLight    = NewHex("#3B82F6") // Blue 500
	InfoDark     = NewHex("#1D4ED8") // Blue 700
)

// Color palettes organized by purpose
var (
	// Status Colors
	StatusColors = ColorPalette{
		"success":  MaterialGreen,
		"error":    MaterialRed,
		"warning":  MaterialAmber,
		"info":     MaterialBlue,
		"debug":    MaterialPurple,
		"critical": MaterialDeepOrange,
		"notice":   MaterialCyan,
	}

	// Brand Colors
	BrandColors = ColorPalette{
		"github":   GithubBlueLight,
		"vscode":   VSCodeBlue,
		"google":   NewHex("#4285F4"),
		"facebook": NewHex("#1877F2"),
		"twitter":  NewHex("#1DA1F2"),
		"linkedin": NewHex("#0A66C2"),
		"discord":  NewHex("#5865F2"),
		"slack":    NewHex("#4A154B"),
	}

	// Theme Palettes
	DraculaPalette = ColorPalette{
		"purple": DraculaPurple,
		"pink":   DraculaPink,
		"green":  DraculaGreen,
		"orange": DraculaOrange,
		"red":    DraculaRed,
		"yellow": DraculaYellow,
		"cyan":   DraculaCyan,
	}

	NordPalette = ColorPalette{
		"blue":      NordBlue,
		"lightblue": NordLightBlue,
		"cyan":      NordCyan,
		"green":     NordGreen,
		"yellow":    NordYellow,
		"orange":    NordOrange,
		"red":       NordRed,
		"purple":    NordPurple,
	}

	TailwindPalette = ColorPalette{
		"red":     TailwindRed,
		"orange":  TailwindOrange,
		"amber":   TailwindAmber,
		"yellow":  TailwindYellow,
		"lime":    TailwindLime,
		"green":   TailwindGreen,
		"emerald": TailwindEmerald,
		"teal":    TailwindTeal,
		"cyan":    TailwindCyan,
		"sky":     TailwindSky,
		"blue":    TailwindBlue,
		"indigo":  TailwindIndigo,
		"violet":  TailwindViolet,
		"purple":  TailwindPurple,
		"fuchsia": TailwindFuchsia,
		"pink":    TailwindPink,
		"rose":    TailwindRose,
	}
)

// ColorPalette represents a collection of named colors
type ColorPalette map[string]Color

// Get returns a color by name from the palette
func (p ColorPalette) Get(name string) (Color, bool) {
	color, exists := p[name]
	return color, exists
}

// Names returns all color names in the palette
func (p ColorPalette) Names() []string {
	names := make([]string, 0, len(p))
	for name := range p {
		names = append(names, name)
	}
	return names
}

// ColorTheme represents a complete theme with specific purpose colors
type ColorTheme struct {
	Name       string
	Success    Color
	Error      Color
	Warning    Color
	Info       Color
	Debug      Color
	Primary    Color
	Secondary  Color
	Accent     Color
	Background Color
	Text       Color
}

// Predefined themes
var (
	DefaultTheme = ColorTheme{
		Name:      "default",
		Success:   MaterialGreen,
		Error:     MaterialRed,
		Warning:   MaterialAmber,
		Info:      MaterialBlue,
		Debug:     MaterialPurple,
		Primary:   MaterialBlue,
		Secondary: MaterialCyan,
		Accent:    MaterialPink,
	}

	DraculaTheme = ColorTheme{
		Name:      "dracula",
		Success:   DraculaGreen,
		Error:     DraculaRed,
		Warning:   DraculaYellow,
		Info:      DraculaCyan,
		Debug:     DraculaPurple,
		Primary:   DraculaPurple,
		Secondary: DraculaPink,
		Accent:    DraculaOrange,
	}

	NordTheme = ColorTheme{
		Name:      "nord",
		Success:   NordGreen,
		Error:     NordRed,
		Warning:   NordYellow,
		Info:      NordBlue,
		Debug:     NordPurple,
		Primary:   NordBlue,
		Secondary: NordCyan,
		Accent:    NordOrange,
	}

	GitHubTheme = ColorTheme{
		Name:      "github",
		Success:   GithubGreenLight,
		Error:     GithubRedLight,
		Warning:   GithubOrangeLight,
		Info:      GithubBlueLight,
		Debug:     GithubPurple,
		Primary:   GithubBlueLight,
		Secondary: GithubGreenLight,
		Accent:    GithubOrangeLight,
	}
)

// AllThemes returns all available themes
var AllThemes = []ColorTheme{
	DefaultTheme,
	DraculaTheme,
	NordTheme,
	GitHubTheme,
}

// GetThemeByName returns a theme by name
func GetThemeByName(name string) (ColorTheme, bool) {
	for _, theme := range AllThemes {
		if theme.Name == name {
			return theme, true
		}
	}
	return ColorTheme{}, false
}
