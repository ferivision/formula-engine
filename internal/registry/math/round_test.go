package math

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestRound_Name(t *testing.T) {
	if (roundFunction{}).Name() != "ROUND" {
		t.Errorf("Name() = %q, want ROUND", (roundFunction{}).Name())
	}
}

func TestRound_ArgBounds(t *testing.T) {
	f := roundFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 2 {
		t.Errorf("bounds = (%d, %d), want (1, 2)", f.MinArgs(), f.MaxArgs())
	}
}

func TestRound_Evaluate_DefaultsToZeroDecimals(t *testing.T) {
	got, err := (roundFunction{}).Evaluate([]registry.Value{3.6})
	if err != nil || got != 4.0 {
		t.Errorf("Evaluate() = (%v, %v), want (4.0, nil)", got, err)
	}
}

func TestRound_Evaluate_WithDecimalPlaces(t *testing.T) {
	got, err := (roundFunction{}).Evaluate([]registry.Value{3.14159, 2.0})
	if err != nil || got != 3.14 {
		t.Errorf("Evaluate() = (%v, %v), want (3.14, nil)", got, err)
	}
}

func TestRound_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (roundFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}

func TestRound_Evaluate_TooManyArgs(t *testing.T) {
	if _, err := (roundFunction{}).Evaluate([]registry.Value{1.0, 2.0, 3.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args")
	}
}
