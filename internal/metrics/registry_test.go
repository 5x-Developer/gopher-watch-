package metrics

import (
	"testing"
)

func TestRegistry_IncidentLifecycle(t *testing.T) {
	reg := GetInstance()
	target := "test-service"

	// 1. Simulate 3 Failures
	reg.UpdateStreak(target, false)
	reg.UpdateStreak(target, false)
	streak := reg.UpdateStreak(target, false)

	if streak != 3 {
		t.Errorf("Expected streak of 3 after 3 failures, got %d", streak)
	}

	// 2. Check Recovery Logic
	oldStreak := reg.GetFailureStreak(target)
	if oldStreak != 3 {
		t.Errorf("Expected old streak to be 3, got %d", oldStreak)
	}

	// 3. Simulate Recovery
	reg.UpdateStreak(target, true)
	if reg.GetFailureStreak(target) != 0 {
		t.Error("Streak should reset to 0 after success")
	}
}
