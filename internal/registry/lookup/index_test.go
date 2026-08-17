package lookup

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestIndex_Name(t *testing.T) {
	if (indexFunction{}).Name() != "INDEX" {
		t.Errorf("Name() = %q, want INDEX", (indexFunction{}).Name())
	}
}

func TestIndex_ArgBounds(t *testing.T) {
	f := indexFunction{}
	if f.MinArgs() != 2 || f.MaxArgs() != 2 {
		t.Errorf("bounds = (%d, %d), want (2, 2)", f.MinArgs(), f.MaxArgs())
	}
}

func TestIndex_ReturnsElementAtOneBasedPosition(t *testing.T) {
	nums := mustArray(t, []any{10.0, 20.0, 30.0})
	got, err := (indexFunction{}).Evaluate([]registry.Value{nums, 2.0})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 20.0 {
		t.Errorf("Evaluate() = %v, want 20 (1-based position 2)", got)
	}
}

func TestIndex_WorksOnRecordArrays(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"sku": "A1"},
		{"sku": "B2"},
	})
	got, err := (indexFunction{}).Evaluate([]registry.Value{orders, 2.0})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	rec := got.(evaluator.Record)
	if rec["sku"] != "B2" {
		t.Errorf("Evaluate() sku = %v, want B2", rec["sku"])
	}
}

func TestIndex_PositionTooHighErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (indexFunction{}).Evaluate([]registry.Value{nums, 5.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want out-of-range error")
	}
}

func TestIndex_PositionZeroOrNegativeErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (indexFunction{}).Evaluate([]registry.Value{nums, 0.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want out-of-range error for position 0")
	}
	if _, err := (indexFunction{}).Evaluate([]registry.Value{nums, -1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want out-of-range error for negative position")
	}
}

func TestIndex_NonArrayFirstArgErrors(t *testing.T) {
	if _, err := (indexFunction{}).Evaluate([]registry.Value{5.0, 1.0}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-array first arg")
	}
}

func TestIndex_WrongArgCount(t *testing.T) {
	nums := mustArray(t, []any{1.0})
	if _, err := (indexFunction{}).Evaluate([]registry.Value{nums}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 1 arg (INDEX needs 2)")
	}
}
