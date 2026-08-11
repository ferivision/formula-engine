package formulaengine

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

// PRD use case 3: chained formulas -- all computed in the correct
// order and returned together, regardless of input order.
func TestEvaluate_ChainedFormulas(t *testing.T) {
	data := map[string]any{"price": 100.0, "quantity": 2.0}

	for _, order := range [][]FormulaInput{
		{
			{Name: "subtotal", Expression: "price * quantity"},
			{Name: "discount", Expression: "subtotal * 0.1"},
			{Name: "total", Expression: "subtotal - discount"},
		},
		{
			{Name: "total", Expression: "subtotal - discount"},
			{Name: "discount", Expression: "subtotal * 0.1"},
			{Name: "subtotal", Expression: "price * quantity"},
		},
	} {
		results, err := Evaluate(order, data)
		if err != nil {
			t.Fatalf("Evaluate() error = %v", err)
		}
		if results["subtotal"].Value != 200.0 {
			t.Errorf("subtotal = %v, want 200", results["subtotal"].Value)
		}
		if results["discount"].Value != 20.0 {
			t.Errorf("discount = %v, want 20", results["discount"].Value)
		}
		if results["total"].Value != 180.0 {
			t.Errorf("total = %v, want 180", results["total"].Value)
		}
	}
}

// PRD use case 5: a circular reference is a call-level error, with no
// partial results, regardless of cycle length.
func TestEvaluate_CircularReference(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "a", Expression: "b + 1"},
		{Name: "b", Expression: "a + 1"},
	}, nil)

	if err == nil {
		t.Fatal("Evaluate() error = nil, want a circular-reference error")
	}
	if results != nil {
		t.Errorf("Evaluate() results = %v, want nil (no partial results on a call-level error)", results)
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrCircularReference {
		t.Errorf("error = %v, want ErrCircularReference", err)
	}
}

// PRD use case 6: a reference to a field/formula not included in the
// call is a call-level error, with no partial results.
func TestEvaluate_MissingReference(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "total", Expression: "subtotal + 1"},
	}, nil)

	if err == nil {
		t.Fatal("Evaluate() error = nil, want an undefined-reference error")
	}
	if results != nil {
		t.Errorf("Evaluate() results = %v, want nil (no partial results on a call-level error)", results)
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrUndefinedReference {
		t.Errorf("error = %v, want ErrUndefinedReference", err)
	}
}

// PRD use case 1: a single, non-chained formula still works.
func TestEvaluate_SimpleCalculation(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "total", Expression: "price * quantity"},
	}, map[string]any{"price": 10.0, "quantity": 3.0})

	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["total"].Value != 30.0 {
		t.Errorf("total = %v, want 30", results["total"].Value)
	}
	if results["total"].Err != nil {
		t.Errorf("total.Err = %v, want nil", results["total"].Err)
	}
}
