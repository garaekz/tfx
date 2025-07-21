// progress/render.go - REDESIGNED
package progress

import (
	"fmt"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/terminal"
)

// Enhanced progress bar rendering with effects and smart color detection
func RenderBar(p *Progress) string {
	percent := float64(p.current) / float64(p.total)

	// Create terminal detector if not exists
	detector := terminal.NewDetector(p.writer)

	// Render label with theme color
	labelColor := p.theme.RenderColor(p.theme.LabelColor, detector)
	label := labelColor + p.label + color.Reset

	// Render progress bar with effects
	var bar string
	if p.theme.EffectEnabled && p.effect != EffectNone {
		bar = p.theme.RenderProgress(percent, p.width, p.effect, detector)
	} else {
		bar = p.theme.renderSolidProgress(int(percent*float64(p.width)), p.width, detector)
	}

	// Render percentage with theme color
	percentColor := p.theme.RenderColor(p.theme.PercentColor, detector)
	percentText := percentColor + fmt.Sprintf("%3d%%", int(percent*100)) + color.Reset

	// Render borders with theme color
	borderColor := p.theme.RenderColor(p.theme.BorderColor, detector)
	leftBorder := borderColor + "[" + color.Reset
	rightBorder := borderColor + "]" + color.Reset

	return fmt.Sprintf("\r%s %s%s%s %s", label, leftBorder, bar, rightBorder, percentText)
}

// Enhanced spinner rendering with color themes
func RenderSpinner(s *Spinner, frame string) string {
	// Create terminal detector if not exists
	detector := terminal.NewDetector(s.writer)

	// Apply theme colors to spinner frame and message
	frameColor := s.theme.RenderColor(s.theme.CompleteColor, detector)
	labelColor := s.theme.RenderColor(s.theme.LabelColor, detector)

	styledFrame := frameColor + frame + color.Reset
	styledMessage := labelColor + s.message + color.Reset

	return fmt.Sprintf("\r%s %s", styledFrame, styledMessage)
}

// Advanced rendering for completion messages
func RenderCompletion(theme ProgressTheme, message string, success bool, detector *terminal.Detector) string {
	var messageColor color.Color
	var icon string

	if success {
		messageColor = theme.CompleteColor
		icon = "✅"
	} else {
		messageColor = color.MaterialRed // Error color
		icon = "❌"
	}

	coloredMessage := theme.RenderColor(messageColor, detector) + message + color.Reset
	return fmt.Sprintf(" %s %s", icon, coloredMessage)
}

// Render progress with time estimation and rate
func RenderAdvancedBar(p *Progress) string {
	detector := terminal.NewDetector(p.writer)

	// Basic bar rendering
	basicBar := RenderBar(p)

	// Add time estimation if progress is started
	if p.started && p.current > 0 {
		elapsed := p.GetElapsed()
		rate := float64(p.current) / elapsed.Seconds()

		if rate > 0 {
			remaining := float64(p.total-p.current) / rate
			eta := fmt.Sprintf("ETA: %ds", int(remaining))

			// Style ETA with theme
			etaColor := p.theme.RenderColor(p.theme.PercentColor, detector)
			styledETA := etaColor + eta + color.Reset

			return basicBar + " " + styledETA
		}
	}

	return basicBar
}

// Multi-line progress display for complex operations
func RenderMultiProgress(bars []*Progress) string {
	if len(bars) == 0 {
		return ""
	}

	detector := terminal.NewDetector(bars[0].writer)
	result := ""

	for i, bar := range bars {
		if i > 0 {
			result += "\n"
		}

		// Add numbering for multiple bars
		numberColor := bar.theme.RenderColor(bar.theme.BorderColor, detector)
		number := numberColor + fmt.Sprintf("%d. ", i+1) + color.Reset

		barRender := RenderBar(bar)
		// Remove the \r from individual bars in multi-progress
		barRender = barRender[1:] // Remove leading \r

		result += number + barRender
	}

	return result
}

// ASCII art style rendering for special occasions
func RenderASCIIBar(p *Progress, style ASCIIStyle) string {
	percent := float64(p.current) / float64(p.total)
	detector := terminal.NewDetector(p.writer)

	var bar string
	filled := int(percent * float64(p.width))

	labelColor := p.theme.RenderColor(p.theme.LabelColor, detector)
	completeColor := p.theme.RenderColor(p.theme.CompleteColor, detector)
	incompleteColor := p.theme.RenderColor(p.theme.IncompleteColor, detector)

	switch style {
	case ASCIIClassic:
		// Classic [=====>    ] style
		for i := 0; i < p.width; i++ {
			if i < filled-1 {
				bar += completeColor + "=" + color.Reset
			} else if i == filled-1 {
				bar += completeColor + ">" + color.Reset
			} else {
				bar += incompleteColor + " " + color.Reset
			}
		}
	case ASCIIDots:
		// Dots style [●●●●○○○○]
		for i := 0; i < p.width; i++ {
			if i < filled {
				bar += completeColor + "●" + color.Reset
			} else {
				bar += incompleteColor + "○" + color.Reset
			}
		}
	case ASCIIBlocks:
		// Unicode blocks [████░░░░]
		for i := 0; i < p.width; i++ {
			if i < filled {
				bar += completeColor + "█" + color.Reset
			} else {
				bar += incompleteColor + "░" + color.Reset
			}
		}
	}

	label := labelColor + p.label + color.Reset
	percentColor := p.theme.RenderColor(p.theme.PercentColor, detector)
	percentText := percentColor + fmt.Sprintf(" %3d%%", int(percent*100)) + color.Reset

	return fmt.Sprintf("\r%s [%s]%s", label, bar, percentText)
}

// ASCII style types
type ASCIIStyle int

const (
	ASCIIClassic ASCIIStyle = iota
	ASCIIDots
	ASCIIBlocks
)

// Render with dynamic width based on terminal size
func RenderResponsiveBar(p *Progress, maxWidth int) string {
	// Adjust bar width based on available space
	labelLen := len(p.label)
	percentLen := 5 // " 100%"
	borderLen := 4  // " [" + "] "

	availableWidth := min(
		// Minimum width
		max(maxWidth-labelLen-percentLen-borderLen,

			10),
		// Maximum width
		60)

	// Temporarily adjust width for rendering
	originalWidth := p.width
	p.width = availableWidth

	result := RenderBar(p)

	// Restore original width
	p.width = originalWidth

	return result
}
