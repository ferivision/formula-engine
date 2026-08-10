package date

import (
	"testing"
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestNow_Name(t *testing.T) {
	if (nowFunction{}).Name() != "NOW" {
		t.Errorf("Name() = %q, want NOW", (nowFunction{}).Name())
	}
}

func TestNow_ArgBounds(t *testing.T) {
	f := nowFunction{}
	if f.MinArgs() != 0 || f.MaxArgs() != 0 {
		t.Errorf("bounds = (%d, %d), want (0, 0)", f.MinArgs(), f.MaxArgs())
	}
}

// NOW must go through the injectable Now var, not time.Now directly,
// so callers (and tests) can make results deterministic per NFR-2.
func TestNow_Evaluate_UsesInjectedClock(t *testing.T) {
	original := Now
	fixed := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	Now = func() time.Time { return fixed }
	defer func() { Now = original }()

	got, err := (nowFunction{}).Evaluate(nil)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	gotTime, ok := got.(time.Time)
	if !ok || !gotTime.Equal(fixed) {
		t.Errorf("Evaluate() = %v, want %v", got, fixed)
	}
}

func TestNow_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (nowFunction{}).Evaluate([]registry.Value{1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 1 arg (NOW takes none)")
	}
}
