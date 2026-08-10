package comparison

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestEquals_Name(t *testing.T) {
	if (equalsFunction{}).Name() != "EQUALS" {
		t.Errorf("Name() = %q, want EQUALS", (equalsFunction{}).Name())
	}
}

func TestEquals_ArgBounds(t *testing.T) {
	f := equalsFunction{}
	if f.MinArgs() != 2 || f.MaxArgs() != 2 {
		t.Errorf("bounds = (%d, %d), want (2, 2)", f.MinArgs(), f.MaxArgs())
	}
}

func TestEquals_NumberNumber(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{5.0, 5.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}

	got, err = (equalsFunction{}).Evaluate([]registry.Value{5.0, 6.0})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestEquals_StringString(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{"foo", "foo"})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}

	got, err = (equalsFunction{}).Evaluate([]registry.Value{"foo", "bar"})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestEquals_BoolBool(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{true, true})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

// Mixed number vs. numeric string: resolved per rfc.md §10's
// number+string rule (attempt numeric parse).
func TestEquals_NumberVsNumericString(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{5.0, "5"})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestEquals_NumberVsNonNumericStringErrors(t *testing.T) {
	if _, err := (equalsFunction{}).Evaluate([]registry.Value{5.0, "abc"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-numeric string")
	}
}

// bool in numeric context: true -> 1.
func TestEquals_BoolVsNumber(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{true, 1.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

// null in arithmetic context -> 0.
func TestEquals_NilVsZero(t *testing.T) {
	got, err := (equalsFunction{}).Evaluate([]registry.Value{nil, 0.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestEquals_WrongArgCount(t *testing.T) {
	if _, err := (equalsFunction{}).Evaluate([]registry.Value{1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 1 arg")
	}
}
