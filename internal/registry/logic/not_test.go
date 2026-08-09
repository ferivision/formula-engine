package logic

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestNot_Name(t *testing.T) {
	if (notFunction{}).Name() != "NOT" {
		t.Errorf("Name() = %q, want NOT", (notFunction{}).Name())
	}
}

func TestNot_ArgBounds(t *testing.T) {
	f := notFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestNot_Evaluate_NegatesTrue(t *testing.T) {
	got, err := (notFunction{}).Evaluate([]registry.Value{true})
	if err != nil || got != false {
		t.Errorf("Evaluate() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestNot_Evaluate_NegatesFalse(t *testing.T) {
	got, err := (notFunction{}).Evaluate([]registry.Value{false})
	if err != nil || got != true {
		t.Errorf("Evaluate() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestNot_Evaluate_WrongArgCount(t *testing.T) {
	if _, err := (notFunction{}).Evaluate([]registry.Value{true, false}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args")
	}
}
