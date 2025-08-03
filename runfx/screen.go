package runfx

import "fmt"

// ScreenManager handles region allocation and optimized clearing
type ScreenManager struct {
	regions map[string][2]int // region name -> [start, end]
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{regions: make(map[string][2]int)}
}

func (s *ScreenManager) AllocateRegion(name string, start, end int) {
	s.regions[name] = [2]int{start, end}
}

// ClearRegion clears the specified region using ANSI codes
func (s *ScreenManager) ClearRegion(name string) {
	region, ok := s.regions[name]
	if !ok {
		return
	}
	start, end := region[0], region[1]
	for line := start; line <= end; line++ {
		// Move to beginning of line and clear
		// ANSI: \033[{line};1H\033[2K
		fmt.Printf("\033[%d;1H\033[2K", line+1)
	}
}

func (s *ScreenManager) Reallocate(name string, newStart, newEnd int) {
	s.regions[name] = [2]int{newStart, newEnd}
}

// GetRegion returns the region bounds for a given name
func (s *ScreenManager) GetRegion(name string) ([2]int, bool) {
	region, ok := s.regions[name]
	return region, ok
}
