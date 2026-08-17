package lookup

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestFind_Name(t *testing.T) {
	if (findFunction{}).Name() != "FIND" {
		t.Errorf("Name() = %q, want FIND", (findFunction{}).Name())
	}
}

func TestFind_ArgBounds(t *testing.T) {
	f := findFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestFind_ReturnsFirstMatchingRecord(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "pending", "sku": "A1"},
		{"status": "shipped", "sku": "B2"},
		{"status": "shipped", "sku": "C3"},
	})

	got, err := (findFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	rec := got.(evaluator.Record)
	if rec["sku"] != "B2" {
		t.Errorf("Evaluate() sku = %v, want B2 (first match)", rec["sku"])
	}
}

func TestFind_NotFoundErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "pending"}})
	if _, err := (findFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want not-found error")
	}
}

func TestFind_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		nil,
		map[string]any{"status": "shipped", "sku": "B2"},
	})

	got, err := (findFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.(evaluator.Record)["sku"] != "B2" {
		t.Errorf("Evaluate() sku = %v, want B2", got.(evaluator.Record)["sku"])
	}
}

func TestFind_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (findFunction{}).Evaluate([]registry.Value{nums, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestFind_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (findFunction{}).Evaluate([]registry.Value{orders, "status"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (FIND needs 3)")
	}
}
