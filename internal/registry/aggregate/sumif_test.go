package aggregate

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestSumif_Name(t *testing.T) {
	if (sumifFunction{}).Name() != "SUMIF" {
		t.Errorf("Name() = %q, want SUMIF", (sumifFunction{}).Name())
	}
}

func TestSumif_ArgBounds(t *testing.T) {
	f := sumifFunction{}
	if f.MinArgs() != 4 || f.MaxArgs() != 4 {
		t.Errorf("bounds = (%d, %d), want (4, 4)", f.MinArgs(), f.MaxArgs())
	}
}

func mustArray(t *testing.T, v any) evaluator.Array {
	t.Helper()
	arr, err := evaluator.NewArray(v)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	return arr
}

func TestSumif_SumsMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped", "qty": 5.0},
		{"status": "pending", "qty": 2.0},
		{"status": "shipped", "qty": 3.0},
	})

	got, err := (sumifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 8.0 {
		t.Errorf("Evaluate() = %v, want 8", got)
	}
}

func TestSumif_EmptyArrayReturnsZero(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	got, err := (sumifFunction{}).Evaluate([]registry.Value{empty, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 0.0 {
		t.Errorf("Evaluate() = %v, want 0", got)
	}
}

func TestSumif_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped", "qty": 5.0},
		nil,
		map[string]any{"status": "shipped", "qty": 3.0},
	})

	got, err := (sumifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 8.0 {
		t.Errorf("Evaluate() = %v, want 8 (null element skipped, not erroring)", got)
	}
}

func TestSumif_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0, 3.0})
	if _, err := (sumifFunction{}).Evaluate([]registry.Value{nums, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestSumif_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (sumifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (SUMIF needs 4)")
	}
}
