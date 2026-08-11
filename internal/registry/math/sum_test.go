package math

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestSum_Name(t *testing.T) {
	if (sumFunction{}).Name() != "SUM" {
		t.Errorf("Name() = %q, want SUM", (sumFunction{}).Name())
	}
}

func TestSum_ArgBounds(t *testing.T) {
	f := sumFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestSum_Evaluate(t *testing.T) {
	got, err := (sumFunction{}).Evaluate([]registry.Value{1.0, 2.0, 3.5})
	if err != nil || got != 6.5 {
		t.Errorf("Evaluate() = (%v, %v), want (6.5, nil)", got, err)
	}
}

func TestSum_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (sumFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
