package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestConcat_Name(t *testing.T) {
	if (concatFunction{}).Name() != "CONCAT" {
		t.Errorf("Name() = %q, want CONCAT", (concatFunction{}).Name())
	}
}

func TestConcat_ArgBounds(t *testing.T) {
	f := concatFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != -1 {
		t.Errorf("bounds = (%d, %d), want (1, -1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestConcat_Evaluate(t *testing.T) {
	got, err := (concatFunction{}).Evaluate([]registry.Value{"foo", "bar", "baz"})
	if err != nil || got != "foobarbaz" {
		t.Errorf("Evaluate() = (%v, %v), want (\"foobarbaz\", nil)", got, err)
	}
}

func TestConcat_Evaluate_NullBecomesEmptyString(t *testing.T) {
	got, err := (concatFunction{}).Evaluate([]registry.Value{"foo", nil, "bar"})
	if err != nil || got != "foobar" {
		t.Errorf("Evaluate() = (%v, %v), want (\"foobar\", nil)", got, err)
	}
}

func TestConcat_Evaluate_NonStringArgErrors(t *testing.T) {
	if _, err := (concatFunction{}).Evaluate([]registry.Value{"foo", 1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-string arg")
	}
}

func TestConcat_Evaluate_TooFewArgs(t *testing.T) {
	if _, err := (concatFunction{}).Evaluate(nil); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 0 args")
	}
}
