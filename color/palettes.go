package color

import "github.com/garaekz/tfx/internal/shared"

// Palette represents a collection of named colors
type Palette map[string]Color

// PaletteConfig provides structured configuration for Palette creation
type PaletteConfig struct {
	Name   string
	Colors map[string]Color
}

// DefaultPaletteConfig returns default configuration for Palette
func DefaultPaletteConfig() PaletteConfig {
	return PaletteConfig{
		Name:   "default",
		Colors: make(map[string]Color),
	}
}

// --- MULTIPATH API: Three Entry Points ---

// 1. EXPRESS: Quick palette creation
func BasicPalette() Palette {
	return StatusPalette()
}

// 2. INSTANTIATED: Config struct
func NewPalette(cfg PaletteConfig, opts ...shared.Option[PaletteConfig]) Palette {
	shared.ApplyOptions(&cfg, opts...)
	return Palette(cfg.Colors)
}

// 3. FLUENT: Functional options
func NewPaletteWith(opts ...shared.Option[PaletteConfig]) Palette {
	cfg := DefaultPaletteConfig()
	return NewPalette(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

func WithPaletteName(name string) shared.Option[PaletteConfig] {
	return func(cfg *PaletteConfig) {
		cfg.Name = name
	}
}

func WithColors(colors map[string]Color) shared.Option[PaletteConfig] {
	return func(cfg *PaletteConfig) {
		if cfg.Colors == nil {
			cfg.Colors = make(map[string]Color)
		}
		for name, color := range colors {
			cfg.Colors[name] = color
		}
	}
}

func WithColor(name string, color Color) shared.Option[PaletteConfig] {
	return func(cfg *PaletteConfig) {
		if cfg.Colors == nil {
			cfg.Colors = make(map[string]Color)
		}
		cfg.Colors[name] = color
	}
}

// --- PALETTE METHODS ---

// Get returns a color by name from the palette
func (p Palette) Get(name string) (Color, bool) {
	color, exists := p[name]
	return color, exists
}

// Set adds or updates a color in the palette
func (p Palette) Set(name string, color Color) {
	p[name] = color
}

// Names returns all color names in the palette
func (p Palette) Names() []string {
	names := make([]string, 0, len(p))
	for name := range p {
		names = append(names, name)
	}
	return names
}

// Merge combines this palette with another
func (p Palette) Merge(other Palette) Palette {
	merged := make(Palette, len(p)+len(other))
	for name, color := range p {
		merged[name] = color
	}
	for name, color := range other {
		merged[name] = color
	}
	return merged
}

// --- PREDEFINED COLORS ---

// Basic Colors (using Color struct)
var (
	// Standard ANSI colors
	ColorBlack   = ANSI(0).WithName("black")
	ColorRed     = ANSI(1).WithName("red")
	ColorGreen   = ANSI(2).WithName("green")
	ColorYellow  = ANSI(3).WithName("yellow")
	ColorBlue    = ANSI(4).WithName("blue")
	ColorMagenta = ANSI(5).WithName("magenta")
	ColorCyan    = ANSI(6).WithName("cyan")
	ColorWhite   = ANSI(7).WithName("white")

	// Bright ANSI colors
	ColorBrightBlack   = ANSI(8).WithName("bright_black")
	ColorBrightRed     = ANSI(9).WithName("bright_red")
	ColorBrightGreen   = ANSI(10).WithName("bright_green")
	ColorBrightYellow  = ANSI(11).WithName("bright_yellow")
	ColorBrightBlue    = ANSI(12).WithName("bright_blue")
	ColorBrightMagenta = ANSI(13).WithName("bright_magenta")
	ColorBrightCyan    = ANSI(14).WithName("bright_cyan")
	ColorBrightWhite   = ANSI(15).WithName("bright_white")

	// Semantic colors
	ColorSuccess = ColorBrightGreen.WithName("success")
	ColorError   = ColorBrightRed.WithName("error")
	ColorWarning = ColorBrightYellow.WithName("warning")
	ColorInfo    = ColorBrightCyan.WithName("info")
	ColorDebug   = ColorBrightMagenta.WithName("debug")
)

// Material Design Colors
var (
	MaterialRed        = Hex("#F44336").WithName("material_red")
	MaterialPink       = Hex("#E91E63").WithName("material_pink")
	MaterialPurple     = Hex("#9C27B0").WithName("material_purple")
	MaterialDeepPurple = Hex("#673AB7").WithName("material_deep_purple")
	MaterialIndigo     = Hex("#3F51B5").WithName("material_indigo")
	MaterialBlue       = Hex("#2196F3").WithName("material_blue")
	MaterialLightBlue  = Hex("#03A9F4").WithName("material_light_blue")
	MaterialCyan       = Hex("#00BCD4").WithName("material_cyan")
	MaterialTeal       = Hex("#009688").WithName("material_teal")
	MaterialGreen      = Hex("#4CAF50").WithName("material_green")
	MaterialLightGreen = Hex("#8BC34A").WithName("material_light_green")
	MaterialLime       = Hex("#CDDC39").WithName("material_lime")
	MaterialYellow     = Hex("#FFEB3B").WithName("material_yellow")
	MaterialAmber      = Hex("#FFC107").WithName("material_amber")
	MaterialOrange     = Hex("#FF9800").WithName("material_orange")
	MaterialDeepOrange = Hex("#FF5722").WithName("material_deep_orange")
)

// Tailwind CSS Colors
var (
	TailwindRed     = Hex("#EF4444").WithName("tailwind_red")
	TailwindOrange  = Hex("#F97316").WithName("tailwind_orange")
	TailwindAmber   = Hex("#F59E0B").WithName("tailwind_amber")
	TailwindYellow  = Hex("#EAB308").WithName("tailwind_yellow")
	TailwindLime    = Hex("#84CC16").WithName("tailwind_lime")
	TailwindGreen   = Hex("#22C55E").WithName("tailwind_green")
	TailwindEmerald = Hex("#10B981").WithName("tailwind_emerald")
	TailwindTeal    = Hex("#14B8A6").WithName("tailwind_teal")
	TailwindCyan    = Hex("#06B6D4").WithName("tailwind_cyan")
	TailwindSky     = Hex("#0EA5E9").WithName("tailwind_sky")
	TailwindBlue    = Hex("#3B82F6").WithName("tailwind_blue")
	TailwindIndigo  = Hex("#6366F1").WithName("tailwind_indigo")
	TailwindViolet  = Hex("#8B5CF6").WithName("tailwind_violet")
	TailwindPurple  = Hex("#A855F7").WithName("tailwind_purple")
	TailwindFuchsia = Hex("#D946EF").WithName("tailwind_fuchsia")
	TailwindPink    = Hex("#EC4899").WithName("tailwind_pink")
	TailwindRose    = Hex("#F43F5E").WithName("tailwind_rose")
)

// Dracula Theme Colors
var (
	DraculaPurple = Hex("#BD93F9").WithName("dracula_purple")
	DraculaPink   = Hex("#FF79C6").WithName("dracula_pink")
	DraculaGreen  = Hex("#50FA7B").WithName("dracula_green")
	DraculaOrange = Hex("#FFB86C").WithName("dracula_orange")
	DraculaRed    = Hex("#FF5555").WithName("dracula_red")
	DraculaYellow = Hex("#F1FA8C").WithName("dracula_yellow")
	DraculaCyan   = Hex("#8BE9FD").WithName("dracula_cyan")
)

// Nord Theme Colors
var (
	NordBlue      = Hex("#5E81AC").WithName("nord_blue")
	NordLightBlue = Hex("#81A1C1").WithName("nord_light_blue")
	NordCyan      = Hex("#88C0D0").WithName("nord_cyan")
	NordGreen     = Hex("#A3BE8C").WithName("nord_green")
	NordYellow    = Hex("#EBCB8B").WithName("nord_yellow")
	NordOrange    = Hex("#D08770").WithName("nord_orange")
	NordRed       = Hex("#BF616A").WithName("nord_red")
	NordPurple    = Hex("#B48EAD").WithName("nord_purple")
)

// GitHub Colors
var (
	GithubGreenLight  = Hex("#28A745").WithName("github_green_light")
	GithubGreenDark   = Hex("#22863A").WithName("github_green_dark")
	GithubRedLight    = Hex("#D73A49").WithName("github_red_light")
	GithubRedDark     = Hex("#B31D28").WithName("github_red_dark")
	GithubBlueLight   = Hex("#0366D6").WithName("github_blue_light")
	GithubBlueDark    = Hex("#005CC5").WithName("github_blue_dark")
	GithubOrangeLight = Hex("#F66A0A").WithName("github_orange_light")
	GithubOrangeDark  = Hex("#E36209").WithName("github_orange_dark")
	GithubPurple      = Hex("#6F42C1").WithName("github_purple")
	GithubYellow      = Hex("#FFD33D").WithName("github_yellow")
)

// VS Code Colors
var (
	VSCodeBlue   = Hex("#007ACC").WithName("vscode_blue")
	VSCodeGreen  = Hex("#608B4E").WithName("vscode_green")
	VSCodeRed    = Hex("#F44747").WithName("vscode_red")
	VSCodeOrange = Hex("#FF8C00").WithName("vscode_orange")
	VSCodePurple = Hex("#C586C0").WithName("vscode_purple")
	VSCodeYellow = Hex("#FFCD3C").WithName("vscode_yellow")
	VSCodeCyan   = Hex("#4EC9B0").WithName("vscode_cyan")
)

// ColorTheme represents a theme with semantic colors for logging
type ColorTheme struct {
	Name      string
	Success   Color
	Error     Color
	Warning   Color
	Info      Color
	Debug     Color
	Primary   Color
	Secondary Color
	Accent    Color
}

// DefaultTheme is the default color theme
var DefaultTheme = ColorTheme{
	Name:      "default",
	Success:   ColorSuccess,
	Error:     ColorError,
	Warning:   ColorWarning,
	Info:      ColorInfo,
	Debug:     ColorDebug,
	Primary:   ColorBlue,
	Secondary: ColorCyan,
	Accent:    MaterialPink,
}

// DraculaTheme is the Dracula color theme
var DraculaTheme = ColorTheme{
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

// MaterialTheme is the Material Design color theme
var MaterialTheme = ColorTheme{
	Name:      "material",
	Success:   MaterialGreen,
	Error:     MaterialRed,
	Warning:   MaterialAmber,
	Info:      MaterialBlue,
	Debug:     MaterialPurple,
	Primary:   MaterialBlue,
	Secondary: MaterialCyan,
	Accent:    MaterialPink,
}

// --- PREDEFINED PALETTES ---

// StatusPalette contains common status colors
func StatusPalette() Palette {
	return Palette{
		"success":  ColorSuccess,
		"error":    ColorError,
		"warning":  ColorWarning,
		"info":     ColorInfo,
		"debug":    ColorDebug,
		"critical": MaterialDeepOrange.WithName("critical"),
		"notice":   MaterialCyan.WithName("notice"),
	}
}

// MaterialPalette contains Material Design colors
func MaterialPalette() Palette {
	return Palette{
		"red":         MaterialRed,
		"pink":        MaterialPink,
		"purple":      MaterialPurple,
		"deep_purple": MaterialDeepPurple,
		"indigo":      MaterialIndigo,
		"blue":        MaterialBlue,
		"light_blue":  MaterialLightBlue,
		"cyan":        MaterialCyan,
		"teal":        MaterialTeal,
		"green":       MaterialGreen,
		"light_green": MaterialLightGreen,
		"lime":        MaterialLime,
		"yellow":      MaterialYellow,
		"amber":       MaterialAmber,
		"orange":      MaterialOrange,
		"deep_orange": MaterialDeepOrange,
	}
}

// DraculaPalette contains Dracula theme colors
func DraculaPalette() Palette {
	return Palette{
		"purple": DraculaPurple,
		"pink":   DraculaPink,
		"green":  DraculaGreen,
		"orange": DraculaOrange,
		"red":    DraculaRed,
		"yellow": DraculaYellow,
		"cyan":   DraculaCyan,
	}
}

// NordPalette contains Nord theme colors
func NordPalette() Palette {
	return Palette{
		"blue":       NordBlue,
		"light_blue": NordLightBlue,
		"cyan":       NordCyan,
		"green":      NordGreen,
		"yellow":     NordYellow,
		"orange":     NordOrange,
		"red":        NordRed,
		"purple":     NordPurple,
	}
}

// TailwindPalette contains Tailwind CSS colors
func TailwindPalette() Palette {
	return Palette{
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
}

// GitHubPalette contains GitHub brand colors
func GitHubPalette() Palette {
	return Palette{
		"green_light":  GithubGreenLight,
		"green_dark":   GithubGreenDark,
		"red_light":    GithubRedLight,
		"red_dark":     GithubRedDark,
		"blue_light":   GithubBlueLight,
		"blue_dark":    GithubBlueDark,
		"orange_light": GithubOrangeLight,
		"orange_dark":  GithubOrangeDark,
		"purple":       GithubPurple,
		"yellow":       GithubYellow,
	}
}

// VSCodePalette contains VS Code theme colors
func VSCodePalette() Palette {
	return Palette{
		"blue":   VSCodeBlue,
		"green":  VSCodeGreen,
		"red":    VSCodeRed,
		"orange": VSCodeOrange,
		"purple": VSCodePurple,
		"yellow": VSCodeYellow,
		"cyan":   VSCodeCyan,
	}
}

// AllPalettes contains all predefined palettes
var AllPalettes = map[string]func() Palette{
	"status":   StatusPalette,
	"material": MaterialPalette,
	"dracula":  DraculaPalette,
	"nord":     NordPalette,
	"tailwind": TailwindPalette,
	"github":   GitHubPalette,
	"vscode":   VSCodePalette,
}

// GetPalette returns a palette by name
func GetPalette(name string) (Palette, bool) {
	paletteFunc, exists := AllPalettes[name]
	if !exists {
		return nil, false
	}
	return paletteFunc(), true
}

// ListPalettes returns all available palette names
func ListPalettes() []string {
	names := make([]string, 0, len(AllPalettes))
	for name := range AllPalettes {
		names = append(names, name)
	}
	return names
}
