package aggregate

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestMaxif_Name(t *testing.T) {
	if (maxifFunction{}).Name() != "MAXIF" {
		t.Errorf("Name() = %q, want MAXIF", (maxifFunction{}).Name())
	}
}

func TestMaxif_ArgBounds(t *testing.T) {
	f := maxifFunction{}
	if f.MinArgs() != 4 || f.MaxArgs() != 4 {
		t.Errorf("bounds = (%d, %d), want (4, 4)", f.MinArgs(), f.MaxArgs())
	}
}

func TestMaxif_FindsMaximumOfMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped", "qty": 4.0},
		{"status": "pending", "qty": 100.0},
		{"status": "shipped", "qty": 9.0},
	})

	got, err := (maxifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 9.0 {
		t.Errorf("Evaluate() = %v, want 9", got)
	}
}

func TestMaxif_EmptyArrayErrors(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	if _, err := (maxifFunction{}).Evaluate([]registry.Value{empty, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for empty array")
	}
}

func TestMaxif_ZeroMatchesErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "pending", "qty": 1.0}})
	if _, err := (maxifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for zero matches")
	}
}

func TestMaxif_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped", "qty": 4.0},
		nil,
		map[string]any{"status": "shipped", "qty": 9.0},
	})

	got, err := (maxifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 9.0 {
		t.Errorf("Evaluate() = %v, want 9", got)
	}
}

func TestMaxif_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (maxifFunction{}).Evaluate([]registry.Value{nums, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestMaxif_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (maxifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (MAXIF needs 4)")
	}
}
