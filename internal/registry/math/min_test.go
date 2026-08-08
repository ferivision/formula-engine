package math

import (
	"testing"

	"github.com/Ferivision/formula-engine/internal/registry"
)

func TestMin_Name(t *testing.T) {
	if (minFunction{}).Name() != "MIN" {
		t.Errorf("Name() = %q, want MIN", (minFunction{}).Name())
	}
}

func TestMin_ArgBounds(t *testing.T) {
	f := minFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestMin_Evaluate(t *testing.T) {
	got, err := (minFunction{}).Evaluate([]registry.Value{4.0, 1.0, 3.0})
	if err != nil || got != 1.0 {
		t.Errorf("Evaluate() = (%v, %v), want (1.0, nil)", got, err)
	}
}

func TestMin_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (minFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
