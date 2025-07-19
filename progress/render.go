package progress

import (
	"fmt"
)

func RenderBar(p *Progress) string {
	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(p.width))
	bar := ""
	for i := 0; i < p.width; i++ {
		if i < filled {
			bar += p.theme.CompleteColor + p.style.FilledChar() + "\033[0m"
		} else {
			bar += p.theme.IncompleteColor + p.style.EmptyChar() + "\033[0m"
		}
	}
	label := p.theme.LabelColor + p.label + "\033[0m"
	return fmt.Sprintf("\r%s [%s] %3d%%", label, bar, int(percent*100))
}
