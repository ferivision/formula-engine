package lookup

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func mustArray(t *testing.T, v any) evaluator.Array {
	t.Helper()
	arr, err := evaluator.NewArray(v)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	return arr
}

func TestVlookup_Name(t *testing.T) {
	if (vlookupFunction{}).Name() != "VLOOKUP" {
		t.Errorf("Name() = %q, want VLOOKUP", (vlookupFunction{}).Name())
	}
}

func TestVlookup_ArgBounds(t *testing.T) {
	f := vlookupFunction{}
	if f.MinArgs() != 4 || f.MaxArgs() != 4 {
		t.Errorf("bounds = (%d, %d), want (4, 4)", f.MinArgs(), f.MaxArgs())
	}
}

func TestVlookup_ReturnsMatchingRowValue(t *testing.T) {
	table := mustArray(t, []map[string]any{
		{"sku": "A1", "price": 10.0},
		{"sku": "B2", "price": 20.0},
	})

	got, err := (vlookupFunction{}).Evaluate([]registry.Value{"B2", table, "sku", "price"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 20.0 {
		t.Errorf("Evaluate() = %v, want 20", got)
	}
}

// PRD use case 5: a missing key is a formula-level error, not a
// call-level one.
func TestVlookup_KeyNotFoundErrors(t *testing.T) {
	table := mustArray(t, []map[string]any{{"sku": "A1", "price": 10.0}})
	if _, err := (vlookupFunction{}).Evaluate([]registry.Value{"ZZ", table, "sku", "price"}); err == nil {
		t.Fatal("Evaluate() error = nil, want not-found error")
	}
}

func TestVlookup_NullElementSkipped(t *testing.T) {
	table := mustArray(t, []any{
		nil,
		map[string]any{"sku": "A1", "price": 10.0},
	})

	got, err := (vlookupFunction{}).Evaluate([]registry.Value{"A1", table, "sku", "price"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 10.0 {
		t.Errorf("Evaluate() = %v, want 10", got)
	}
}

func TestVlookup_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (vlookupFunction{}).Evaluate([]registry.Value{1.0, nums, "sku", "price"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestVlookup_WrongArgCount(t *testing.T) {
	table := mustArray(t, []map[string]any{{"sku": "A1"}})
	if _, err := (vlookupFunction{}).Evaluate([]registry.Value{"A1", table, "sku"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (VLOOKUP needs 4)")
	}
}
