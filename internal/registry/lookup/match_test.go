package lookup

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestMatch_Name(t *testing.T) {
	if (matchFunction{}).Name() != "MATCH" {
		t.Errorf("Name() = %q, want MATCH", (matchFunction{}).Name())
	}
}

func TestMatch_ArgBounds(t *testing.T) {
	f := matchFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestMatch_ReturnsOneBasedIndex(t *testing.T) {
	arr := mustArray(t, []map[string]any{
		{"sku": "A1"},
		{"sku": "B2"},
		{"sku": "C3"},
	})

	got, err := (matchFunction{}).Evaluate([]registry.Value{"B2", arr, "sku"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2 (1-based index)", got)
	}
}

func TestMatch_KeyNotFoundErrors(t *testing.T) {
	arr := mustArray(t, []map[string]any{{"sku": "A1"}})
	if _, err := (matchFunction{}).Evaluate([]registry.Value{"ZZ", arr, "sku"}); err == nil {
		t.Fatal("Evaluate() error = nil, want not-found error")
	}
}

func TestMatch_NullElementSkipped(t *testing.T) {
	arr := mustArray(t, []any{
		nil,
		map[string]any{"sku": "A1"},
	})

	got, err := (matchFunction{}).Evaluate([]registry.Value{"A1", arr, "sku"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2 (null element still occupies position 1)", got)
	}
}

func TestMatch_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (matchFunction{}).Evaluate([]registry.Value{1.0, nums, "sku"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestMatch_EmptyArrayErrors(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	if _, err := (matchFunction{}).Evaluate([]registry.Value{"A1", empty, "sku"}); err == nil {
		t.Fatal("Evaluate() error = nil, want not-found error for an empty array")
	}
}

func TestMatch_WrongArgCount(t *testing.T) {
	arr := mustArray(t, []map[string]any{{"sku": "A1"}})
	if _, err := (matchFunction{}).Evaluate([]registry.Value{"A1", arr}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (MATCH needs 3)")
	}
}
