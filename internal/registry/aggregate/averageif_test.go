package aggregate

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

func TestAverageif_Name(t *testing.T) {
	if (averageifFunction{}).Name() != "AVERAGEIF" {
		t.Errorf("Name() = %q, want AVERAGEIF", (averageifFunction{}).Name())
	}
}

func TestAverageif_ArgBounds(t *testing.T) {
	f := averageifFunction{}
	if f.MinArgs() != 4 || f.MaxArgs() != 4 {
		t.Errorf("bounds = (%d, %d), want (4, 4)", f.MinArgs(), f.MaxArgs())
	}
}

func TestAverageif_AveragesMatchingRecords(t *testing.T) {
	orders := mustArray(t, []map[string]any{
		{"status": "shipped", "qty": 4.0},
		{"status": "pending", "qty": 100.0},
		{"status": "shipped", "qty": 6.0},
	})

	got, err := (averageifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 5.0 {
		t.Errorf("Evaluate() = %v, want 5", got)
	}
}

// Deliberately different from SUMIF/COUNTIF: average of nothing is
// undefined, per rfc.md §9.
func TestAverageif_EmptyArrayErrors(t *testing.T) {
	empty := mustArray(t, []map[string]any{})
	if _, err := (averageifFunction{}).Evaluate([]registry.Value{empty, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for empty array")
	}
}

func TestAverageif_ZeroMatchesErrors(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "pending", "qty": 1.0}})
	if _, err := (averageifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for zero matches")
	}
}

func TestAverageif_NullElementSkipped(t *testing.T) {
	orders := mustArray(t, []any{
		map[string]any{"status": "shipped", "qty": 4.0},
		nil,
		map[string]any{"status": "shipped", "qty": 6.0},
	})

	got, err := (averageifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped", "qty"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got != 5.0 {
		t.Errorf("Evaluate() = %v, want 5 (null element excluded from average)", got)
	}
}

func TestAverageif_ScalarArrayErrors(t *testing.T) {
	nums := mustArray(t, []any{1.0, 2.0})
	if _, err := (averageifFunction{}).Evaluate([]registry.Value{nums, "status", "shipped", "qty"}); err == nil {
		t.Fatal("Evaluate() error = nil, want type error for scalar array")
	}
}

func TestAverageif_WrongArgCount(t *testing.T) {
	orders := mustArray(t, []map[string]any{{"status": "shipped"}})
	if _, err := (averageifFunction{}).Evaluate([]registry.Value{orders, "status", "shipped"}); err == nil {
		t.Fatal("Evaluate() error = nil, want error for 3 args (AVERAGEIF needs 4)")
	}
}
