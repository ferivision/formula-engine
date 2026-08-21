package formulaengine

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/evaluator"
)

// This file is black-box acceptance coverage for each of PRD (0002)
// §5's six use cases, exercised only through the public Evaluate()
// API -- no internal package imports beyond apperror's error codes.

// Use case 1: conditional sum -- sum qty only for orders where
// status == "shipped".
func TestPRDUseCase1_ConditionalSum(t *testing.T) {
	data := map[string]any{
		"orders": []map[string]any{
			{"status": "shipped", "qty": 5.0},
			{"status": "pending", "qty": 2.0},
			{"status": "shipped", "qty": 3.0},
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "shippedTotal", Expression: `SUMIF(orders, "status", "shipped", "qty")`},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["shippedTotal"].Err != nil {
		t.Fatalf("shippedTotal.Err = %v, want nil", results["shippedTotal"].Err)
	}
	if results["shippedTotal"].Value != 8.0 {
		t.Errorf("shippedTotal.Value = %v, want 8", results["shippedTotal"].Value)
	}
}

// Use case 2: filter then count -- filter records down to those
// matching a condition, then count how many remain. Also proves a
// formula can consume another formula's Array result by name.
func TestPRDUseCase2_FilterThenCount(t *testing.T) {
	data := map[string]any{
		"orders": []map[string]any{
			{"status": "shipped", "qty": 5.0},
			{"status": "pending", "qty": 2.0},
			{"status": "shipped", "qty": 3.0},
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "shippedCount", Expression: `COUNTIF(shippedOrders, "status", "shipped")`},
		{Name: "shippedOrders", Expression: `FILTER(orders, "status", "shipped")`},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["shippedOrders"].Err != nil {
		t.Fatalf("shippedOrders.Err = %v, want nil", results["shippedOrders"].Err)
	}
	if results["shippedCount"].Err != nil {
		t.Fatalf("shippedCount.Err = %v, want nil", results["shippedCount"].Err)
	}
	if results["shippedCount"].Value != 2.0 {
		t.Errorf("shippedCount.Value = %v, want 2", results["shippedCount"].Value)
	}
}

// Use case 3: deduplicate -- return only unique values, preserving
// first-seen order.
func TestPRDUseCase3_Deduplicate(t *testing.T) {
	data := map[string]any{
		"skus": []any{"A1", "B2", "A1", "C3", "B2"},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "uniqueSkus", Expression: "UNIQUE(skus)"},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["uniqueSkus"].Err != nil {
		t.Fatalf("uniqueSkus.Err = %v, want nil", results["uniqueSkus"].Err)
	}
	arr, ok := results["uniqueSkus"].Value.(evaluator.Array)
	if !ok {
		t.Fatalf("uniqueSkus.Value = %T, want evaluator.Array", results["uniqueSkus"].Value)
	}
	want := []evaluator.Value{"A1", "B2", "C3"}
	if len(arr.Elements) != len(want) {
		t.Fatalf("Elements = %v, want %v", arr.Elements, want)
	}
	for i := range want {
		if arr.Elements[i] != want[i] {
			t.Errorf("Elements[%d] = %v, want %v", i, arr.Elements[i], want[i])
		}
	}
}

// Use case 4: reference lookup -- given a key and a separate Array
// acting as a lookup table, return the matching row's value for a
// named column.
func TestPRDUseCase4_ReferenceLookup(t *testing.T) {
	data := map[string]any{
		"key": "B2",
		"prices": []map[string]any{
			{"sku": "A1", "price": 10.0},
			{"sku": "B2", "price": 20.0},
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "price", Expression: `VLOOKUP(key, prices, "sku", "price")`},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["price"].Err != nil {
		t.Fatalf("price.Err = %v, want nil", results["price"].Err)
	}
	if results["price"].Value != 20.0 {
		t.Errorf("price.Value = %v, want 20", results["price"].Value)
	}
}

// Use case 5: missing key in lookup -- a clear, specific error rather
// than a crash or a silently wrong result, and it must not block an
// unrelated formula in the same call.
func TestPRDUseCase5_MissingKeyInLookup(t *testing.T) {
	data := map[string]any{
		"prices": []map[string]any{
			{"sku": "A1", "price": 10.0},
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "price", Expression: `VLOOKUP("ZZ", prices, "sku", "price")`},
		{Name: "unrelated", Expression: "1 + 1"},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (formula-level error)", err)
	}
	if results["price"].Err == nil {
		t.Fatal("price.Err = nil, want a not-found error")
	}
	var fe *apperror.FormulaError
	if !errors.As(results["price"].Err, &fe) || fe.Code != apperror.ErrLookupNotFound {
		t.Errorf("price.Err = %v, want ErrLookupNotFound", results["price"].Err)
	}
	if results["unrelated"].Err != nil || results["unrelated"].Value != 2.0 {
		t.Errorf("unrelated = %+v, want {Value: 2, Err: nil}", results["unrelated"])
	}
}

// Use case 6: empty Array input -- behavior is defined and
// consistent, not undefined. SUMIF of nothing is 0; AVERAGEIF of
// nothing is a formula-level error (average of nothing is
// undefined), and that error doesn't block an unrelated formula.
func TestPRDUseCase6_EmptyArrayInput(t *testing.T) {
	data := map[string]any{
		"orders": []map[string]any{},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "sum", Expression: `SUMIF(orders, "status", "shipped", "qty")`},
		{Name: "avg", Expression: `AVERAGEIF(orders, "status", "shipped", "qty")`},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (formula-level error)", err)
	}
	if results["sum"].Err != nil {
		t.Fatalf("sum.Err = %v, want nil", results["sum"].Err)
	}
	if results["sum"].Value != 0.0 {
		t.Errorf("sum.Value = %v, want 0", results["sum"].Value)
	}
	if results["avg"].Err == nil {
		t.Fatal("avg.Err = nil, want a formula-level error (average of nothing is undefined)")
	}
}
