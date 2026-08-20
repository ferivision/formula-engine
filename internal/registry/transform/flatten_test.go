package transform

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestFlatten_Name(t *testing.T) {
	if (flattenFunction{}).Name() != "FLATTEN" {
		t.Errorf("Name() = %q, want FLATTEN", (flattenFunction{}).Name())
	}
}

func TestFlatten_ArgBounds(t *testing.T) {
	f := flattenFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 1 {
		t.Errorf("bounds = (%d, %d), want (1, 1)", f.MinArgs(), f.MaxArgs())
	}
}

func TestFlatten_ConcatenatesScalarSubArrays(t *testing.T) {
	groups := mustArray(t, []any{
		[]any{1.0, 2.0},
		[]any{3.0, 4.0},
	})

	got, err := (flattenFunction{}).Evaluate([]registry.Value{groups})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if arr.IsRecord {
		t.Error("IsRecord = true, want false")
	}
	want := []evaluator.Value{1.0, 2.0, 3.0, 4.0}
	if len(arr.Elements) != len(want) {
		t.Fatalf("Elements = %v, want %v", arr.Elements, want)
	}
	for i := range want {
		if arr.Elements[i] != want[i] {
			t.Errorf("Elements[%d] = %v, want %v", i, arr.Elements[i], want[i])
		}
	}
}

func TestFlatten_ConcatenatesRecordSubArrays(t *testing.T) {
	groups := mustArray(t, []any{
		[]map[string]any{{"sku": "A1"}},
		[]map[string]any{{"sku": "B2"}, {"sku": "C3"}},
	})

	got, err := (flattenFunction{}).Evaluate([]registry.Value{groups})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true")
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("len(Elements) = %d, want 3", len(arr.Elements))
	}
	if arr.Elements[0].(evaluator.Record)["sku"] != "A1" {
		t.Errorf("Elements[0] sku = %v, want A1", arr.Elements[0].(evaluator.Record)["sku"])
	}
}

func TestFlatten_EmptyOuterArrayReturnsEmptyArray(t *testing.T) {
	empty := mustArray(t, []any{})
	got, err := (flattenFunction{}).Evaluate([]registry.Value{empty})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if len(got.(evaluator.Array).Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(got.(evaluator.Array).Elements))
	}
}

func TestFlatten_MixingRecordAndScalarSubArraysErrors(t *testing.T) {
	groups := mustArray(t, []any{
		[]any{1.0, 2.0},
		[]map[string]any{{"sku": "A1"}},
	})
	if _, err := (flattenFunction{}).Evaluate([]registry.Value{groups}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for mixed record/scalar sub-arrays")
	}
}

// rfc.md §9: a null slot in the outer array contributes nothing to
// the flattened result, rather than erroring as "not an array".
func TestFlatten_NullElementInOuterArraySkipped(t *testing.T) {
	groups := mustArray(t, []any{
		[]any{1.0, 2.0},
		nil,
		[]any{3.0},
	})

	got, err := (flattenFunction{}).Evaluate([]registry.Value{groups})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	want := []evaluator.Value{1.0, 2.0, 3.0}
	if len(arr.Elements) != len(want) {
		t.Fatalf("Elements = %v, want %v (null sub-array slot contributes nothing)", arr.Elements, want)
	}
	for i := range want {
		if arr.Elements[i] != want[i] {
			t.Errorf("Elements[%d] = %v, want %v", i, arr.Elements[i], want[i])
		}
	}
}

func TestFlatten_NonArrayElementErrors(t *testing.T) {
	groups := mustArray(t, []any{5.0})
	if _, err := (flattenFunction{}).Evaluate([]registry.Value{groups}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for a non-array element")
	}
}

func TestFlatten_WrongArgCount(t *testing.T) {
	groups := mustArray(t, []any{[]any{1.0}})
	if _, err := (flattenFunction{}).Evaluate([]registry.Value{groups, "extra"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (FLATTEN takes 1)")
	}
}
