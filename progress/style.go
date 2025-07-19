package progress

type ProgressStyle int

const (
	ProgressStyleBar ProgressStyle = iota
	ProgressStyleDots
	ProgressStyleArrows
	ProgressStyleAscii
)

// StyleChars contiene los símbolos de cada style
var styleChars = map[ProgressStyle]struct {
	Filled, Empty string
}{
	ProgressStyleBar:    {"█", "░"},
	ProgressStyleDots:   {"●", "○"},
	ProgressStyleArrows: {">", "-"},
	ProgressStyleAscii:  {"=", "-"},
}

func (s ProgressStyle) FilledChar() string { return styleChars[s].Filled }
func (s ProgressStyle) EmptyChar() string  { return styleChars[s].Empty }
