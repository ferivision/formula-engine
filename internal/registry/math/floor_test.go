package math

import (
	"testing"

	"github.com/Ferivision/formula-engine/internal/registry"
)

func TestFloor_Name(t *testing.T) {
	if (floorFunction{}).Name() != "FLOOR" {
		t.Errorf("Name() = %q, want FLOOR", (floorFunction{}).Name())
	}
}

func TestFloor_ArgBounds(t *testing.T) {
	f := floorFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestFloor_Evaluate(t *testing.T) {
	got, err := (floorFunction{}).Evaluate([]registry.Value{3.7})
	if err != nil || got != 3.0 {
		t.Errorf("Evaluate() = (%v, %v), want (3.0, nil)", got, err)
	}
}

func TestFloor_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (floorFunction{}).Evaluate([]registry.Value{1.0, 2.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args")
	}
}
