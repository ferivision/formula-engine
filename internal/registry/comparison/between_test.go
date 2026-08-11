package comparison

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestBetween_Name(t *testing.T) {
	if (betweenFunction{}).Name() != "BETWEEN" {
		t.Errorf("Name() = %q, want BETWEEN", (betweenFunction{}).Name())
	}
}

func TestBetween_ArgBounds(t *testing.T) {
	f := betweenFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestBetween_ValueWithinRange(t *testing.T) {
	got, err := (betweenFunction{}).Evaluate([]registry.Value{5.0, 1.0, 10.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

// BETWEEN is inclusive of both bounds.
func TestBetween_ValueAtLowerBoundIsInclusive(t *testing.T) {
	got, err := (betweenFunction{}).Evaluate([]registry.Value{1.0, 1.0, 10.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestBetween_ValueAtUpperBoundIsInclusive(t *testing.T) {
	got, err := (betweenFunction{}).Evaluate([]registry.Value{10.0, 1.0, 10.0})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestBetween_ValueOutsideRange(t *testing.T) {
	got, err := (betweenFunction{}).Evaluate([]registry.Value{11.0, 1.0, 10.0})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestBetween_NonNumericArgErrors(t *testing.T) {
	if _, err := (betweenFunction{}).Evaluate([]registry.Value{"abc", 1.0, 10.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-numeric value")
	}
}

func TestBetween_WrongArgCount(t *testing.T) {
	if _, err := (betweenFunction{}).Evaluate([]registry.Value{1.0, 2.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args")
	}
}
