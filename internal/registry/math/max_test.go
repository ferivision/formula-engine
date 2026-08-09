package math

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestMax_Name(t *testing.T) {
	if (maxFunction{}).Name() != "MAX" {
		t.Errorf("Name() = %q, want MAX", (maxFunction{}).Name())
	}
}

func TestMax_ArgBounds(t *testing.T) {
	f := maxFunction{}
	if f.MinArgs() != 1 {
		t.Errorf("MinArgs() = %d, want 1", f.MinArgs())
	}
	if f.MaxArgs() != -1 {
		t.Errorf("MaxArgs() = %d, want -1 (unbounded)", f.MaxArgs())
	}
}

func TestMax_Evaluate(t *testing.T) {
	got, err := (maxFunction{}).Evaluate([]registry.Value{1.0, 5.0, 3.0})
	if err != nil || got != 5.0 {
		t.Errorf("Evaluate() = (%v, %v), want (5.0, nil)", got, err)
	}
}

func TestMax_Evaluate_TooFewArgs(t *testing.T) {
	_, err := (maxFunction{}).Evaluate(nil)
	if err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}

func TestMax_Evaluate_NonNumericArg(t *testing.T) {
	_, err := (maxFunction{}).Evaluate([]registry.Value{1.0, "x"})
	if err == nil {
		t.Fatal("Evaluate() error = nil, want type error")
	}
}
