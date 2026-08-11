package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestTrim_Name(t *testing.T) {
	if (trimFunction{}).Name() != "TRIM" {
		t.Errorf("Name() = %q, want TRIM", (trimFunction{}).Name())
	}
}

func TestTrim_ArgBounds(t *testing.T) {
	f := trimFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestTrim_Evaluate(t *testing.T) {
	got, err := (trimFunction{}).Evaluate([]registry.Value{"  hello  "})
	if err != nil || got != "hello" {
		t.Errorf("Evaluate() = (%v, %v), want (\"hello\", nil)", got, err)
	}
}

func TestTrim_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (trimFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
