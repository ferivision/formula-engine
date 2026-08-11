package math

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestCeil_Name(t *testing.T) {
	if (ceilFunction{}).Name() != "CEIL" {
		t.Errorf("Name() = %q, want CEIL", (ceilFunction{}).Name())
	}
}

func TestCeil_ArgBounds(t *testing.T) {
	f := ceilFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestCeil_Evaluate(t *testing.T) {
	got, err := (ceilFunction{}).Evaluate([]registry.Value{3.2})
	if err != nil || got != 4.0 {
		t.Errorf("Evaluate() = (%v, %v), want (4.0, nil)", got, err)
	}
}

func TestCeil_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (ceilFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
