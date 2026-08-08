package math

import (
	"testing"

	"github.com/Ferivision/formula-engine/internal/registry"
)

func TestAvg_Name(t *testing.T) {
	if (avgFunction{}).Name() != "AVG" {
		t.Errorf("Name() = %q, want AVG", (avgFunction{}).Name())
	}
}

func TestAvg_ArgBounds(t *testing.T) {
	f := avgFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestAvg_Evaluate(t *testing.T) {
	got, err := (avgFunction{}).Evaluate([]registry.Value{2.0, 4.0, 6.0})
	if err != nil || got != 4.0 {
		t.Errorf("Evaluate() = (%v, %v), want (4.0, nil)", got, err)
	}
}

func TestAvg_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (avgFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
