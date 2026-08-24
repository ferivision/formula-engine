package aggregate

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestMinif_Name(t *testing.T) {
	if (minifFunction{}).Name() != "MINIF" {
		t.Errorf("Name() = %q, want MINIF", (minifFunction{}).Name())
	}
}

func TestMinif_ArgBounds(t *testing.T) {
	f := minifFunction{}
	if f.MinArgs() != 4 || f.MaxArgs() != 4 {
		t.Errorf("bounds = (%d, %d), want (4, 4)", f.MinArgs(), f.MaxArgs())
	}
}

func TestMinif_FindsMinimumOfMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped", "qty": 9.0},
		{"status": "pending", "qty": 1.0},
		{"status": "shipped", "qty": 4.0},
	})

	got, err := (minifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 4.0 {
		t.Errorf("Evaluate() = %v, want 4", got)
	}
}

func TestMinif_EmptyArrayErrors(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	if _, err := (minifFunction{}).Evaluate([]registry.Value{empty, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for empty array")
	}
}

func TestMinif_ZeroMatchesErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "pending", "qty": 1.0}})
	if _, err := (minifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for zero matches")
	}
}

func TestMinif_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped", "qty": 9.0},
		nil,
		map[string]any{"status": "shipped", "qty": 4.0},
	})

	got, err := (minifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 4.0 {
		t.Errorf("Evaluate() = %v, want 4", got)
	}
}

func TestMinif_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (minifFunction{}).Evaluate([]registry.Value{nums, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestMinif_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (minifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (MINIF needs 4)")
	}
}
