package logic

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestAnd_Name(t *testing.T) {
	if (andFunction{}).Name() != "AND" {
		t.Errorf("Name() = %q, want AND", (andFunction{}).Name())
	}
}

func TestAnd_ArgBounds(t *testing.T) {
	f := andFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestAnd_Evaluate_AllTrue(t *testing.T) {
	got, err := (andFunction{}).Evaluate([]registry.Value{true, true, true})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestAnd_Evaluate_OneFalse(t *testing.T) {
	got, err := (andFunction{}).Evaluate([]registry.Value{true, false, true})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestAnd_Evaluate_NonBoolArgErrors(t *testing.T) {
	if _, err := (andFunction{}).Evaluate([]registry.Value{true, 1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-bool arg")
	}
}

func TestAnd_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (andFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
