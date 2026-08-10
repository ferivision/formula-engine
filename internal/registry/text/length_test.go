package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestLength_Name(t *testing.T) {
	if (lengthFunction{}).Name() != "LENGTH" {
		t.Errorf("Name() = %q, want LENGTH", (lengthFunction{}).Name())
	}
}

func TestLength_ArgBounds(t *testing.T) {
	f := lengthFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestLength_Evaluate_ASCII(t *testing.T) {
	got, err := (lengthFunction{}).Evaluate([]registry.Value{"hello"})
	if err != nil || got != 5.0 {
		t.Errorf("Evaluate() = (%v, %v), want (5, nil)", got, err)
	}
}

// "héllo" has 5 runes but 6 bytes (é is 2 bytes in UTF-8) -- LENGTH
// must count runes, not bytes.
func TestLength_Evaluate_UnicodeCountsRunesNotBytes(t *testing.T) {
	got, err := (lengthFunction{}).Evaluate([]registry.Value{"héllo"})
	if err != nil || got != 5.0 {
		t.Errorf("Evaluate() = (%v, %v), want (5, nil)", got, err)
	}
}

func TestLength_Evaluate_EmptyString(t *testing.T) {
	got, err := (lengthFunction{}).Evaluate([]registry.Value{""})
	if err != nil || got != 0.0 {
		t.Errorf("Evaluate() = (%v, %v), want (0, nil)", got, err)
	}
}

func TestLength_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (lengthFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
