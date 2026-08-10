package date

import (
	"testing"
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestDateAdd_Name(t *testing.T) {
	if (dateAddFunction{}).Name() != "DATE_ADD" {
		t.Errorf("Name() = %q, want DATE_ADD", (dateAddFunction{}).Name())
	}
}

func TestDateAdd_ArgBounds(t *testing.T) {
	f := dateAddFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestDateAdd_Days(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := (dateAddFunction{}).Evaluate([]registry.Value{start, 10.0, "days"})
	want := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	if err != nil || !got.(time.Time).Equal(want) {
		t.Errorf("Evaluate() = (%v, %v), want (%v, nil)", got, err, want)
	}
}

// Adding 2 months to November rolls over both the month and the year.
func TestDateAdd_MonthAndYearRollover(t *testing.T) {
	start := time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC)
	got, err := (dateAddFunction{}).Evaluate([]registry.Value{start, 2.0, "months"})
	want := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if err != nil || !got.(time.Time).Equal(want) {
		t.Errorf("Evaluate() = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestDateAdd_Years(t *testing.T) {
	start := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	got, err := (dateAddFunction{}).Evaluate([]registry.Value{start, 1.0, "years"})
	want := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if err != nil || !got.(time.Time).Equal(want) {
		t.Errorf("Evaluate() = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestDateAdd_NegativeAmountSubtracts(t *testing.T) {
	start := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	got, err := (dateAddFunction{}).Evaluate([]registry.Value{start, -5.0, "days"})
	want := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	if err != nil || !got.(time.Time).Equal(want) {
		t.Errorf("Evaluate() = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestDateAdd_InvalidUnitErrors(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := (dateAddFunction{}).Evaluate([]registry.Value{start, 1.0, "fortnights"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for unsupported unit")
	}
}

func TestDateAdd_WrongArgCount(t *testing.T) {
	if _, err := (dateAddFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
