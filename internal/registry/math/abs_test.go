package math

import (
	"testing"

	"github.com/Ferivision/formula-engine/internal/registry"
)

func TestAbs_Name(t *testing.T) {
	if (absFunction{}).Name() != "ABS" {
		t.Errorf("Name() = %q, want ABS", (absFunction{}).Name())
	}
}

func TestAbs_ArgBounds(t *testing.T) {
	f := absFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestAbs_Evaluate_NegativeInput(t *testing.T) {
	got, err := (absFunction{}).Evaluate([]registry.Value{-5.0})
	if err != nil || got != 5.0 {
		t.Errorf("Evaluate() = (%v, %v), want (5.0, nil)", got, err)
	}
}

func TestAbs_Evaluate_PositiveInput(t *testing.T) {
	got, err := (absFunction{}).Evaluate([]registry.Value{5.0})
	if err != nil || got != 5.0 {
		t.Errorf("Evaluate() = (%v, %v), want (5.0, nil)", got, err)
	}
}

func TestAbs_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (absFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
