package progress

import (
	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
)

// ProgressTheme now uses Color structs for smart terminal optimization
type ProgressTheme struct {
	Name            string
	CompleteColor   color.Color
	IncompleteColor color.Color
	LabelColor      color.Color
	PercentColor    color.Color
	BorderColor     color.Color // New: for frames/brackets
	EffectEnabled   bool        // New: enable gradient/rainbow effects
}

// Enhanced theme creation from existing color palettes
func NewThemeFromPalette(name string, palette color.Palette) ProgressTheme {
	// Smart color selection from palette
	complete, _ := palette.Get("green")
	if complete == (color.Color{}) {
		complete = color.MaterialGreen
	}

	incomplete, _ := palette.Get("gray")
	if incomplete == (color.Color{}) {
		incomplete = color.NewANSI(8) // Dark gray
	}

	label, _ := palette.Get("blue")
	if label == (color.Color{}) {
		label = color.MaterialBlue
	}

	return ProgressTheme{
		Name:            name,
		CompleteColor:   complete,
		IncompleteColor: incomplete,
		LabelColor:      label,
		PercentColor:    label,
		BorderColor:     color.NewANSI(7), // Light gray
		EffectEnabled:   false,
	}
}

// Professional themes based on existing color palettes
var (
	// Material Design Theme
	MaterialTheme = ProgressTheme{
		Name:            "material",
		CompleteColor:   color.MaterialGreen,
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.MaterialBlue,
		PercentColor:    color.MaterialCyan,
		BorderColor:     color.NewANSI(7),
		EffectEnabled:   false,
	}

	// Dracula Theme (enhanced from your original)
	DraculaTheme = ProgressTheme{
		Name:            "dracula",
		CompleteColor:   color.DraculaGreen,
		IncompleteColor: color.RGB(68, 71, 90), // Dracula background
		LabelColor:      color.DraculaPurple,
		PercentColor:    color.DraculaCyan,
		BorderColor:     color.DraculaPink,
		EffectEnabled:   true, // Dracula loves effects!
	}

	// Nord Theme
	NordTheme = ProgressTheme{
		Name:            "nord",
		CompleteColor:   color.NordGreen,
		IncompleteColor: color.RGB(59, 66, 82), // Nord dark
		LabelColor:      color.NordBlue,
		PercentColor:    color.NordCyan,
		BorderColor:     color.NordLightBlue,
		EffectEnabled:   false,
	}

	// GitHub Theme
	GitHubTheme = ProgressTheme{
		Name:            "github",
		CompleteColor:   color.GithubGreenLight,
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.GithubBlueLight,
		PercentColor:    color.GithubBlueLight,
		BorderColor:     color.NewANSI(7),
		EffectEnabled:   false,
	}

	// Tailwind Theme
	TailwindTheme = ProgressTheme{
		Name:            "tailwind",
		CompleteColor:   color.TailwindGreen,
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.TailwindBlue,
		PercentColor:    color.TailwindCyan,
		BorderColor:     color.TailwindIndigo,
		EffectEnabled:   false,
	}

	// VS Code Theme
	VSCodeTheme = ProgressTheme{
		Name:            "vscode",
		CompleteColor:   color.VSCodeGreen,
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.VSCodeBlue,
		PercentColor:    color.VSCodeCyan,
		BorderColor:     color.NewANSI(7),
		EffectEnabled:   false,
	}

	// Rainbow Theme (special effects)
	RainbowTheme = ProgressTheme{
		Name:            "rainbow",
		CompleteColor:   color.MaterialGreen, // Base color, effects override
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.MaterialPurple,
		PercentColor:    color.MaterialCyan,
		BorderColor:     color.NewANSI(7),
		EffectEnabled:   true, // Rainbow effects!
	}

	// Neon Theme (bright and vibrant)
	NeonTheme = ProgressTheme{
		Name:            "neon",
		CompleteColor:   color.RGB(0, 255, 0), // Bright green
		IncompleteColor: color.NewANSI(8),
		LabelColor:      color.RGB(0, 150, 255), // Bright blue
		PercentColor:    color.RGB(255, 0, 255), // Bright magenta
		BorderColor:     color.RGB(255, 255, 0), // Bright yellow
		EffectEnabled:   true,
	}
)

// AllProgressThemes for easy iteration
var AllProgressThemes = []ProgressTheme{
	MaterialTheme,
	DraculaTheme,
	NordTheme,
	GitHubTheme,
	TailwindTheme,
	VSCodeTheme,
	RainbowTheme,
	NeonTheme,
}

// Smart color rendering based on terminal capabilities
func (pt ProgressTheme) RenderColor(c color.Color, detector *terminal.Detector) string {
	if detector == nil {
		return c.Render(color.ModeANSI) // Fallback
	}

	terminalMode := color.Mode(detector.GetMode())
	return c.Render(terminalMode)
}

// Effect types for enhanced progress bars
type ProgressEffect int

const (
	EffectNone ProgressEffect = iota
	EffectGradient
	EffectRainbow
	EffectPulse
	EffectGlow
)

// Enhanced progress rendering with effects
func (pt ProgressTheme) RenderProgress(percent float64, width int, effect ProgressEffect, detector *terminal.Detector) string {
	filled := int(percent * float64(width))

	switch effect {
	case EffectRainbow:
		return pt.renderRainbowProgress(filled, width, detector)
	case EffectGradient:
		return pt.renderGradientProgress(filled, width, detector)
	case EffectGlow:
		return pt.renderGlowProgress(filled, width, detector)
	default:
		return pt.renderSolidProgress(filled, width, detector)
	}
}

// Solid color progress (standard)
func (pt ProgressTheme) renderSolidProgress(filled, width int, detector *terminal.Detector) string {
	completeColor := pt.RenderColor(pt.CompleteColor, detector)
	incompleteColor := pt.RenderColor(pt.IncompleteColor, detector)

	bar := ""
	for i := range width {
		if i < filled {
			bar += completeColor + "█" + color.Reset
		} else {
			bar += incompleteColor + "░" + color.Reset
		}
	}
	return bar
}

// Rainbow progress effect
func (pt ProgressTheme) renderRainbowProgress(filled, width int, detector *terminal.Detector) string {
	rainbowColors := []color.Color{
		color.MaterialRed,
		color.MaterialOrange,
		color.MaterialYellow,
		color.MaterialGreen,
		color.MaterialBlue,
		color.MaterialPurple,
	}

	incompleteColor := pt.RenderColor(pt.IncompleteColor, detector)

	bar := ""
	for i := range width {
		if i < filled {
			colorIndex := i % len(rainbowColors)
			currentColor := pt.RenderColor(rainbowColors[colorIndex], detector)
			bar += currentColor + "█" + color.Reset
		} else {
			bar += incompleteColor + "░" + color.Reset
		}
	}
	return bar
}

// Gradient progress effect
func (pt ProgressTheme) renderGradientProgress(filled, width int, detector *terminal.Detector) string {
	startColor := pt.CompleteColor
	endColor := color.MaterialCyan // Could be configurable
	incompleteColor := pt.RenderColor(pt.IncompleteColor, detector)

	bar := ""
	for i := range width {
		if i < filled {
			// Calculate interpolation (simplified)
			ratio := float64(i) / float64(filled)
			if ratio < 0.5 {
				currentColor := pt.RenderColor(startColor, detector)
				bar += currentColor + "█" + color.Reset
			} else {
				currentColor := pt.RenderColor(endColor, detector)
				bar += currentColor + "█" + color.Reset
			}
		} else {
			bar += incompleteColor + "░" + color.Reset
		}
	}
	return bar
}

// Glow effect (brightness variation)
func (pt ProgressTheme) renderGlowProgress(filled, width int, detector *terminal.Detector) string {
	centerColor := pt.CompleteColor
	edgeColor := color.RGB(
		pt.CompleteColor.R/2,
		pt.CompleteColor.G/2,
		pt.CompleteColor.B/2,
	) // Dimmed version
	incompleteColor := pt.RenderColor(pt.IncompleteColor, detector)

	bar := ""
	for i := range width {
		if i < filled {
			// Glow effect: brighter in center, dimmer at edges
			distanceFromCenter := float64(abs(i-filled/2)) / float64(filled/2)
			if distanceFromCenter < 0.3 {
				currentColor := pt.RenderColor(centerColor, detector)
				bar += currentColor + "█" + color.Reset
			} else {
				currentColor := pt.RenderColor(edgeColor, detector)
				bar += currentColor + "█" + color.Reset
			}
		} else {
			bar += incompleteColor + "░" + color.Reset
		}
	}
	return bar
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetThemeByName returns a theme by name
func GetThemeByName(name string) (ProgressTheme, bool) {
	for _, theme := range AllProgressThemes {
		if theme.Name == name {
			return theme, true
		}
	}
	return MaterialTheme, false // Default fallback
}

// CreateCustomTheme allows users to create themes with specific colors
func CreateCustomTheme(name string, complete, incomplete, label color.Color) ProgressTheme {
	return ProgressTheme{
		Name:            name,
		CompleteColor:   complete,
		IncompleteColor: incomplete,
		LabelColor:      label,
		PercentColor:    label,
		BorderColor:     color.NewANSI(7),
		EffectEnabled:   false,
	}
}
