package transform

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestUnique_Name(t *testing.T) {
	if (uniqueFunction{}).Name() != "UNIQUE" {
		t.Errorf("Name() = %q, want UNIQUE", (uniqueFunction{}).Name())
	}
}

func TestUnique_ArgBounds(t *testing.T) {
	f := uniqueFunction{}
	if f.MinArgs() != 1 || f.MaxArgs() != 2 {
		t.Errorf("bounds = (%d, %d), want (1, 2)", f.MinArgs(), f.MaxArgs())
	}
}

// PRD use case 3: deduplicate, preserving first-seen order.
func TestUnique_DeduplicatesScalarsPreservingFirstSeenOrder(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0, 1.0, 3.0, 2.0})

	got, err := (uniqueFunction{}).Evaluate([]registry.Value{nums})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	want := []evaluator.Value{1.0, 2.0, 3.0}
	if len(arr.Elements) != len(want) {
		t.Fatalf("Elements = %v, want %v", arr.Elements, want)
	}
	for i := range want {
		if arr.Elements[i] != want[i] {
			t.Errorf("Elements[%d] = %v, want %v", i, arr.Elements[i], want[i])
		}
	}
}

func TestUnique_DeduplicatesRecordsByField(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "A1", "qty": 1.0},
		{"sku": "B2", "qty": 2.0},
		{"sku": "A1", "qty": 3.0},
	})

	got, err := (uniqueFunction{}).Evaluate([]registry.Value{orders, "sku"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if len(arr.Elements) != 2 {
		t.Fatalf("len(Elements) = %d, want 2", len(arr.Elements))
	}
	if arr.Elements[0].(evaluator.Record)["sku"] != "A1" {
		t.Errorf("Elements[0] sku = %v, want A1 (first-seen)", arr.Elements[0].(evaluator.Record)["sku"])
	}
	if arr.Elements[1].(evaluator.Record)["sku"] != "B2" {
		t.Errorf("Elements[1] sku = %v, want B2", arr.Elements[1].(evaluator.Record)["sku"])
	}
}

func TestUnique_EmptyArrayReturnsEmptyArray(t *testing.T) {
	empty := mustArray(t, []any{})
	got, err := (uniqueFunction{}).Evaluate([]registry.Value{empty})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if len(got.(evaluator.Array).Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(got.(evaluator.Array).Elements))
	}
}

func TestUnique_RecordArrayWithoutFieldErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"sku": "A1"}})
	if _, err := (uniqueFunction{}).Evaluate([]registry.Value{orders}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error: record array needs a field argument")
	}
}

func TestUnique_DoesNotMutateInputArray(t *testing.T) {
	nums := mustArray(t, []any{1.0, 1.0, 2.0})
	original := len(nums.Elements)

	if _, err := (uniqueFunction{}).Evaluate([]registry.Value{nums}); err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if len(nums.Elements) != original {
		t.Errorf("input array was mutated: len(Elements) = %d, want unchanged %d", len(nums.Elements), original)
	}
}

func TestUnique_WrongArgCount(t *testing.T) {
	nums := mustArray(t, []any{1.0})
	if _, err := (uniqueFunction{}).Evaluate([]registry.Value{nums, "a", "b"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (UNIQUE takes at most 2)")
	}
}
