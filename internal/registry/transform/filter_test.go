package transform

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

func TestFilter_Name(t *testing.T) {
	if (filterFunction{}).Name() != "FILTER" {
		t.Errorf("Name() = %q, want FILTER", (filterFunction{}).Name())
	}
}

func TestFilter_ArgBounds(t *testing.T) {
	f := filterFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestFilter_ReturnsOnlyMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped", "qty": 5.0},
		{"status": "pending", "qty": 2.0},
		{"status": "shipped", "qty": 3.0},
	})

	got, err := (filterFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr, ok := got.(evaluator.Array)
	if !ok {
		t.Fatalf("Evaluate() = %T, want Array", got)
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("len(Elements) = %d, want 2", len(arr.Elements))
	}
	for _, elem := range arr.Elements {
		rec := elem.(evaluator.Record)
		if rec["status"] != "shipped" {
			t.Errorf("Elements contains non-matching record %v", rec)
		}
	}
}

func TestFilter_EmptyArrayReturnsEmptyArray(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	got, err := (filterFunction{}).Evaluate([]registry.Value{empty, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if len(arr.Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(arr.Elements))
	}
}

func TestFilter_DoesNotMutateInputArray(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped"},
		{"status": "pending"},
	})
	original := len(orders.Elements)

	if _, err := (filterFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if len(orders.Elements) != original {
		t.Errorf("input array was mutated: len(Elements) = %d, want unchanged %d", len(orders.Elements), original)
	}
	if orders.Elements[0].(evaluator.Record)["status"] != "shipped" || orders.Elements[1].(evaluator.Record)["status"] != "pending" {
		t.Error("input array elements were mutated")
	}
}

// rfc.md §9: a null element can't be meaningfully matched against a
// condition, so it's excluded from the result rather than erroring.
func TestFilter_NullElementExcludedFromResult(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped"},
		nil,
		map[string]any{"status": "pending"},
	})

	got, err := (filterFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if len(arr.Elements) != 1 {
		t.Fatalf("len(Elements) = %d, want 1 (null element excluded, doesn't error)", len(arr.Elements))
	}
	if arr.Elements[0].(evaluator.Record)["status"] != "shipped" {
		t.Errorf("Elements[0] = %v, want the shipped record", arr.Elements[0])
	}
}

func TestFilter_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (filterFunction{}).Evaluate([]registry.Value{nums, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestFilter_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (filterFunction{}).Evaluate([]registry.Value{orders, "status"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (FILTER needs 3)")
	}
}
