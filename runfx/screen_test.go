package runfx

import (
	"testing"
)

func TestScreenManagerBasic(t *testing.T) {
	sm := NewScreenManager()

	// Test allocation
	sm.AllocateRegion("test1", 0, 5)
	sm.AllocateRegion("test2", 6, 10)

	// Test retrieval
	region1, ok1 := sm.GetRegion("test1")
	if !ok1 {
		t.Error("Region test1 not found")
	}
	if region1[0] != 0 || region1[1] != 5 {
		t.Errorf("Region test1 has wrong bounds: %v", region1)
	}

	region2, ok2 := sm.GetRegion("test2")
	if !ok2 {
		t.Error("Region test2 not found")
	}
	if region2[0] != 6 || region2[1] != 10 {
		t.Errorf("Region test2 has wrong bounds: %v", region2)
	}

	// Test non-existent region
	_, ok3 := sm.GetRegion("nonexistent")
	if ok3 {
		t.Error("Non-existent region should not be found")
	}
}

func TestScreenManagerReallocation(t *testing.T) {
	sm := NewScreenManager()

	// Initial allocation
	sm.AllocateRegion("test", 0, 5)

	region, ok := sm.GetRegion("test")
	if !ok {
		t.Fatal("Region not found after allocation")
	}
	if region[0] != 0 || region[1] != 5 {
		t.Errorf("Initial region has wrong bounds: %v", region)
	}

	// Reallocate
	sm.Reallocate("test", 10, 15)

	region, ok = sm.GetRegion("test")
	if !ok {
		t.Fatal("Region not found after reallocation")
	}
	if region[0] != 10 || region[1] != 15 {
		t.Errorf("Reallocated region has wrong bounds: %v", region)
	}
}

func TestScreenManagerClearRegion(t *testing.T) {
	sm := NewScreenManager()

	// Allocate region
	sm.AllocateRegion("test", 0, 5)

	// Clear should not panic (actual clearing depends on terminal)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ClearRegion panicked: %v", r)
		}
	}()

	sm.ClearRegion("test")

	// Clear non-existent region should not panic
	sm.ClearRegion("nonexistent")
}

func TestScreenManagerOverlappingRegions(t *testing.T) {
	sm := NewScreenManager()

	// Allocate overlapping regions (this is allowed)
	sm.AllocateRegion("region1", 0, 10)
	sm.AllocateRegion("region2", 5, 15)

	region1, ok1 := sm.GetRegion("region1")
	region2, ok2 := sm.GetRegion("region2")

	if !ok1 || !ok2 {
		t.Fatal("Regions not found")
	}

	if region1[0] != 0 || region1[1] != 10 {
		t.Errorf("Region1 has wrong bounds: %v", region1)
	}

	if region2[0] != 5 || region2[1] != 15 {
		t.Errorf("Region2 has wrong bounds: %v", region2)
	}
}

func TestScreenManagerReallocateNonExistent(t *testing.T) {
	sm := NewScreenManager()

	// Reallocate non-existent region (should create it)
	sm.Reallocate("newregion", 0, 5)

	region, ok := sm.GetRegion("newregion")
	if !ok {
		t.Error("Reallocated region not found")
	}
	if region[0] != 0 || region[1] != 5 {
		t.Errorf("Reallocated region has wrong bounds: %v", region)
	}
}
