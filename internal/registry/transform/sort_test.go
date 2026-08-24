package transform

import (
	"os"
	"strings"
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestSort_Name(t *testing.T) {
	if (sortFunction{}).Name() != "SORT" {
		t.Errorf("Name() = %q, want SORT", (sortFunction{}).Name())
	}
}

func TestSort_ArgBounds(t *testing.T) {
	f := sortFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func skus(t *testing.T, arr evaluator.Array) []string {
	t.Helper()
	out := make([]string, len(arr.Elements))
	for i, e := range arr.Elements {
		out[i] = e.(evaluator.Record)["sku"].(string)
	}
	return out
}

func TestSort_Ascending(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "C", "qty": 3.0},
		{"sku": "A", "qty": 1.0},
		{"sku": "B", "qty": 2.0},
	})

	got, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "asc"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	order := skus(t, got.(evaluator.Array))
	want := []string{"A", "B", "C"}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order = %v, want %v", order, want)
		}
	}
}

func TestSort_Descending(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "A", "qty": 1.0},
		{"sku": "B", "qty": 2.0},
		{"sku": "C", "qty": 3.0},
	})

	got, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "desc"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	order := skus(t, got.(evaluator.Array))
	want := []string{"C", "B", "A"}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order = %v, want %v", order, want)
		}
	}
}

// PRD NFR-3 / rfc.md §4: SORT must use sort.SliceStable, never
// sort.Slice -- elements sharing an equal sort key must keep their
// original relative order.
func TestSort_StableForEqualKeys(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "first", "qty": 5.0},
		{"sku": "second", "qty": 5.0},
		{"sku": "third", "qty": 5.0},
	})

	got, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "asc"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	order := skus(t, got.(evaluator.Array))
	want := []string{"first", "second", "third"}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order = %v, want %v (original relative order must be preserved for equal keys)", order, want)
		}
	}
}

func TestSort_InvalidDirectionErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"qty": 1.0}})
	if _, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "sideways"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for invalid direction")
	}
}

func TestSort_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{3.0, 1.0, 2.0})
	if _, err := (sortFunction{}).Evaluate([]registry.Value{nums, "qty", "asc"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestSort_DoesNotMutateInputArray(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "C", "qty": 3.0},
		{"sku": "A", "qty": 1.0},
	})

	if _, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "asc"}); err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if orders.Elements[0].(evaluator.Record)["sku"] != "C" || orders.Elements[1].(evaluator.Record)["sku"] != "A" {
		t.Error("input array order was mutated")
	}
}

func TestSort_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"qty": 1.0}})
	if _, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (SORT needs 3)")
	}
}

func TestSort_EmptyArrayReturnsEmptyArray(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	got, err := (sortFunction{}).Evaluate([]registry.Value{empty, "qty", "asc"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if len(got.(evaluator.Array).Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(got.(evaluator.Array).Elements))
	}
}

// rfc.md §9: SORT reorders, it doesn't filter -- a null element stays
// in the result (unlike FILTER/aggregates, which exclude it), sorting
// as if its key were the arithmetic-context zero value.
func TestSort_NullElementStaysInResultSortedAsZero(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"sku": "positive", "qty": 5.0},
		nil,
		map[string]any{"sku": "negative", "qty": -5.0},
	})

	got, err := (sortFunction{}).Evaluate([]registry.Value{orders, "qty", "asc"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	arr := got.(evaluator.Array)
	if len(arr.Elements) != 3 {
		t.Fatalf("len(Elements) = %d, want 3 (null element preserved, not dropped)", len(arr.Elements))
	}
	if arr.Elements[0].(evaluator.Record)["sku"] != "negative" {
		t.Errorf("Elements[0] = %v, want the qty=-5 record first", arr.Elements[0])
	}
	if arr.Elements[1] != nil {
		t.Errorf("Elements[1] = %v, want nil sorted between -5 and 5 (as key 0)", arr.Elements[1])
	}
	if arr.Elements[2].(evaluator.Record)["sku"] != "positive" {
		t.Errorf("Elements[2] = %v, want the qty=5 record last", arr.Elements[2])
	}
}

// Static guard, in addition to the behavioral stability test above:
// this file must never call the non-stable sort.Slice.
func TestSort_SourceNeverUsesUnstableSort(t *testing.T) {
	src, err := os.ReadFile("sort.go")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	body := string(src)
	if strings.Contains(body, "sort.Slice(") {
		t.Error("sort.go calls sort.Slice() -- must use sort.SliceStable() only, per PRD NFR-3")
	}
	if !strings.Contains(body, "sort.SliceStable(") {
		t.Error("sort.go doesn't call sort.SliceStable() at all")
	}
}
