package color

import "github.com/garaekz/tfx/internal/share"

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
func NewPalette(cfg PaletteConfig, opts ...share.Option[PaletteConfig]) Palette {
	share.ApplyOptions(&cfg, opts...)
	return Palette(cfg.Colors)
}

// 3. FLUENT: Functional options + DSL chaining support
func NewPaletteWith(opts ...share.Option[PaletteConfig]) Palette {
	cfg := DefaultPaletteConfig()
	return NewPalette(cfg, opts...)
}

// --- FUNCTIONAL OPTIONS ---

func WithPaletteName(name string) share.Option[PaletteConfig] {
	return func(cfg *PaletteConfig) {
		cfg.Name = name
	}
}

func WithColors(colors map[string]Color) share.Option[PaletteConfig] {
	return func(cfg *PaletteConfig) {
		if cfg.Colors == nil {
			cfg.Colors = make(map[string]Color)
		}
		for name, color := range colors {
			cfg.Colors[name] = color
		}
	}
}

func WithColor(name string, color Color) share.Option[PaletteConfig] {
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
	ColorBlack   = NewANSI(0).WithName("black")
	ColorRed     = NewANSI(1).WithName("red")
	ColorGreen   = NewANSI(2).WithName("green")
	ColorYellow  = NewANSI(3).WithName("yellow")
	ColorBlue    = NewANSI(4).WithName("blue")
	ColorMagenta = NewANSI(5).WithName("magenta")
	ColorCyan    = NewANSI(6).WithName("cyan")
	ColorWhite   = NewANSI(7).WithName("white")

	// Bright ANSI colors
	ColorBrightBlack   = NewANSI(8).WithName("bright_black")
	ColorBrightRed     = NewANSI(9).WithName("bright_red")
	ColorBrightGreen   = NewANSI(10).WithName("bright_green")
	ColorBrightYellow  = NewANSI(11).WithName("bright_yellow")
	ColorBrightBlue    = NewANSI(12).WithName("bright_blue")
	ColorBrightMagenta = NewANSI(13).WithName("bright_magenta")
	ColorBrightCyan    = NewANSI(14).WithName("bright_cyan")
	ColorBrightWhite   = NewANSI(15).WithName("bright_white")

	// Semantic colors
	ColorSuccess = ColorBrightGreen.WithName("success")
	ColorError   = ColorBrightRed.WithName("error")
	ColorWarning = ColorBrightYellow.WithName("warning")
	ColorInfo    = ColorBrightCyan.WithName("info")
	ColorDebug   = ColorBrightMagenta.WithName("debug")
)

// Material Design Colors
var (
	MaterialRed        = NewHex("#F44336").WithName("material_red")
	MaterialPink       = NewHex("#E91E63").WithName("material_pink")
	MaterialPurple     = NewHex("#9C27B0").WithName("material_purple")
	MaterialDeepPurple = NewHex("#673AB7").WithName("material_deep_purple")
	MaterialIndigo     = NewHex("#3F51B5").WithName("material_indigo")
	MaterialBlue       = NewHex("#2196F3").WithName("material_blue")
	MaterialLightBlue  = NewHex("#03A9F4").WithName("material_light_blue")
	MaterialCyan       = NewHex("#00BCD4").WithName("material_cyan")
	MaterialTeal       = NewHex("#009688").WithName("material_teal")
	MaterialGreen      = NewHex("#4CAF50").WithName("material_green")
	MaterialLightGreen = NewHex("#8BC34A").WithName("material_light_green")
	MaterialLime       = NewHex("#CDDC39").WithName("material_lime")
	MaterialYellow     = NewHex("#FFEB3B").WithName("material_yellow")
	MaterialAmber      = NewHex("#FFC107").WithName("material_amber")
	MaterialOrange     = NewHex("#FF9800").WithName("material_orange")
	MaterialDeepOrange = NewHex("#FF5722").WithName("material_deep_orange")
)

// Tailwind CSS Colors
var (
	TailwindRed     = NewHex("#EF4444").WithName("tailwind_red")
	TailwindOrange  = NewHex("#F97316").WithName("tailwind_orange")
	TailwindAmber   = NewHex("#F59E0B").WithName("tailwind_amber")
	TailwindYellow  = NewHex("#EAB308").WithName("tailwind_yellow")
	TailwindLime    = NewHex("#84CC16").WithName("tailwind_lime")
	TailwindGreen   = NewHex("#22C55E").WithName("tailwind_green")
	TailwindEmerald = NewHex("#10B981").WithName("tailwind_emerald")
	TailwindTeal    = NewHex("#14B8A6").WithName("tailwind_teal")
	TailwindCyan    = NewHex("#06B6D4").WithName("tailwind_cyan")
	TailwindSky     = NewHex("#0EA5E9").WithName("tailwind_sky")
	TailwindBlue    = NewHex("#3B82F6").WithName("tailwind_blue")
	TailwindIndigo  = NewHex("#6366F1").WithName("tailwind_indigo")
	TailwindViolet  = NewHex("#8B5CF6").WithName("tailwind_violet")
	TailwindPurple  = NewHex("#A855F7").WithName("tailwind_purple")
	TailwindFuchsia = NewHex("#D946EF").WithName("tailwind_fuchsia")
	TailwindPink    = NewHex("#EC4899").WithName("tailwind_pink")
	TailwindRose    = NewHex("#F43F5E").WithName("tailwind_rose")
)

// Dracula Theme Colors
var (
	DraculaPurple = NewHex("#BD93F9").WithName("dracula_purple")
	DraculaPink   = NewHex("#FF79C6").WithName("dracula_pink")
	DraculaGreen  = NewHex("#50FA7B").WithName("dracula_green")
	DraculaOrange = NewHex("#FFB86C").WithName("dracula_orange")
	DraculaRed    = NewHex("#FF5555").WithName("dracula_red")
	DraculaYellow = NewHex("#F1FA8C").WithName("dracula_yellow")
	DraculaCyan   = NewHex("#8BE9FD").WithName("dracula_cyan")
)

// Nord Theme Colors
var (
	NordBlue      = NewHex("#5E81AC").WithName("nord_blue")
	NordLightBlue = NewHex("#81A1C1").WithName("nord_light_blue")
	NordCyan      = NewHex("#88C0D0").WithName("nord_cyan")
	NordGreen     = NewHex("#A3BE8C").WithName("nord_green")
	NordYellow    = NewHex("#EBCB8B").WithName("nord_yellow")
	NordOrange    = NewHex("#D08770").WithName("nord_orange")
	NordRed       = NewHex("#BF616A").WithName("nord_red")
	NordPurple    = NewHex("#B48EAD").WithName("nord_purple")
)

// GitHub Colors
var (
	GithubGreenLight  = NewHex("#28A745").WithName("github_green_light")
	GithubGreenDark   = NewHex("#22863A").WithName("github_green_dark")
	GithubRedLight    = NewHex("#D73A49").WithName("github_red_light")
	GithubRedDark     = NewHex("#B31D28").WithName("github_red_dark")
	GithubBlueLight   = NewHex("#0366D6").WithName("github_blue_light")
	GithubBlueDark    = NewHex("#005CC5").WithName("github_blue_dark")
	GithubOrangeLight = NewHex("#F66A0A").WithName("github_orange_light")
	GithubOrangeDark  = NewHex("#E36209").WithName("github_orange_dark")
	GithubPurple      = NewHex("#6F42C1").WithName("github_purple")
	GithubYellow      = NewHex("#FFD33D").WithName("github_yellow")
)

// VS Code Colors
var (
	VSCodeBlue   = NewHex("#007ACC").WithName("vscode_blue")
	VSCodeGreen  = NewHex("#608B4E").WithName("vscode_green")
	VSCodeRed    = NewHex("#F44747").WithName("vscode_red")
	VSCodeOrange = NewHex("#FF8C00").WithName("vscode_orange")
	VSCodePurple = NewHex("#C586C0").WithName("vscode_purple")
	VSCodeYellow = NewHex("#FFCD3C").WithName("vscode_yellow")
	VSCodeCyan   = NewHex("#4EC9B0").WithName("vscode_cyan")
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
