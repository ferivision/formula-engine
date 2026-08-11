package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestUpper_Name(t *testing.T) {
	if (upperFunction{}).Name() != "UPPER" {
		t.Errorf("Name() = %q, want UPPER", (upperFunction{}).Name())
	}
}

func TestUpper_ArgBounds(t *testing.T) {
	f := upperFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestUpper_Evaluate(t *testing.T) {
	got, err := (upperFunction{}).Evaluate([]registry.Value{"hello"})
	if err != nil || got != "HELLO" {
		t.Errorf("Evaluate() = (%v, %v), want (\"HELLO\", nil)", got, err)
	}
}

func TestUpper_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (upperFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
