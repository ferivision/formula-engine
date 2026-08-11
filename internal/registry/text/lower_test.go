package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestLower_Name(t *testing.T) {
	if (lowerFunction{}).Name() != "LOWER" {
		t.Errorf("Name() = %q, want LOWER", (lowerFunction{}).Name())
	}
}

func TestLower_ArgBounds(t *testing.T) {
	f := lowerFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestLower_Evaluate(t *testing.T) {
	got, err := (lowerFunction{}).Evaluate([]registry.Value{"HELLO"})
	if err != nil || got != "hello" {
		t.Errorf("Evaluate() = (%v, %v), want (\"hello\", nil)", got, err)
	}
}

func TestLower_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (lowerFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
