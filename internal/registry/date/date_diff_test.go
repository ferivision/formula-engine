package date

import (
	"testing"
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestDateDiff_Name(t *testing.T) {
	if (dateDiffFunction{}).Name() != "DATE_DIFF" {
		t.Errorf("Name() = %q, want DATE_DIFF", (dateDiffFunction{}).Name())
	}
}

func TestDateDiff_ArgBounds(t *testing.T) {
	f := dateDiffFunction{}
	if f.MinArgs() != 2 || f.MaxArgs() != 2 {
		t.Errorf("bounds = (%d, %d), want (2, 2)", f.MinArgs(), f.MaxArgs())
	}
}

func TestDateDiff_SimpleDays(t *testing.T) {
	d1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	got, err := (dateDiffFunction{}).Evaluate([]registry.Value{d1, d2})
	if err != nil || got != 10.0 {
		t.Errorf("Evaluate() = (%v, %v), want (10, nil)", got, err)
	}
}

// 2024 is a leap year -- Feb 29 exists, so Feb 28 -> Mar 1 is 2 days,
// not 1.
func TestDateDiff_LeapYear(t *testing.T) {
	d1 := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	got, err := (dateDiffFunction{}).Evaluate([]registry.Value{d1, d2})
	if err != nil || got != 2.0 {
		t.Errorf("Evaluate() = (%v, %v), want (2, nil)", got, err)
	}
}

// US DST began at 2am on 2025-03-09, so noon-to-noon from Mar 8 to
// Mar 9 in America/New_York only spans 23 real hours. A naive
// duration/24h calculation would wrongly report ~0.958 days; the
// calendar-date-based result must still be exactly 1.
func TestDateDiff_AcrossDSTTransition(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("LoadLocation() error = %v", err)
	}
	d1 := time.Date(2025, 3, 8, 12, 0, 0, 0, loc)
	d2 := time.Date(2025, 3, 9, 12, 0, 0, 0, loc)

	if d2.Sub(d1) == 24*time.Hour {
		t.Fatal("test setup invalid: expected a non-24h gap across the DST transition")
	}

	got, err := (dateDiffFunction{}).Evaluate([]registry.Value{d1, d2})
	if err != nil || got != 1.0 {
		t.Errorf("Evaluate() = (%v, %v), want (1, nil) despite the 23-hour actual gap", got, err)
	}
}

func TestDateDiff_NegativeWhenReversed(t *testing.T) {
	d1 := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := (dateDiffFunction{}).Evaluate([]registry.Value{d1, d2})
	if err != nil || got != -10.0 {
		t.Errorf("Evaluate() = (%v, %v), want (-10, nil)", got, err)
	}
}

func TestDateDiff_WrongArgCount(t *testing.T) {
	if _, err := (dateDiffFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
