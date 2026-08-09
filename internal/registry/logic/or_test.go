package logic

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestOr_Name(t *testing.T) {
	if (orFunction{}).Name() != "OR" {
		t.Errorf("Name() = %q, want OR", (orFunction{}).Name())
	}
}

func TestOr_ArgBounds(t *testing.T) {
	f := orFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestOr_Evaluate_OneTrue(t *testing.T) {
	got, err := (orFunction{}).Evaluate([]registry.Value{false, true, false})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestOr_Evaluate_AllFalse(t *testing.T) {
	got, err := (orFunction{}).Evaluate([]registry.Value{false, false})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestOr_Evaluate_NonBoolArgErrors(t *testing.T) {
	if _, err := (orFunction{}).Evaluate([]registry.Value{false, "x"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-bool arg")
	}
}

func TestOr_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (orFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
