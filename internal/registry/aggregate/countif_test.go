package aggregate

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestCountif_Name(t *testing.T) {
	if (countifFunction{}).Name() != "COUNTIF" {
		t.Errorf("Name() = %q, want COUNTIF", (countifFunction{}).Name())
	}
}

func TestCountif_ArgBounds(t *testing.T) {
	f := countifFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestCountif_CountsMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped"},
		{"status": "pending"},
		{"status": "shipped"},
	})

	got, err := (countifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2", got)
	}
}

func TestCountif_EmptyArrayReturnsZero(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	got, err := (countifFunction{}).Evaluate([]registry.Value{empty, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 0.0 {
		t.Errorf("Evaluate() = %v, want 0", got)
	}
}

func TestCountif_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped"},
		nil,
		map[string]any{"status": "shipped"},
	})

	got, err := (countifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2 (null element skipped)", got)
	}
}

func TestCountif_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (countifFunction{}).Evaluate([]registry.Value{nums, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestCountif_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (countifFunction{}).Evaluate([]registry.Value{orders, "status"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (COUNTIF needs 3)")
	}
}
