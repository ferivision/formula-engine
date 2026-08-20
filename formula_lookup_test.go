package formulaengine

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

// PRD (0002) use case 5: a missing lookup key is a formula-level
// error -- it must not block an unrelated formula in the same call,
// same guarantee feature-0001 already proved for other formula-level
// errors (see formula_partial_success_test.go). rfc.md §10: the error
// is the dedicated ErrLookupNotFound (Phase 7), not the interim
// ErrRuntime used before F2-FASE-7.1.
func TestEvaluate_VlookupNotFoundDoesNotBlockSiblings(t *testing.T) {
	data := map[string]any{
		"prices": []map[string]any{
			{"sku": "A1", "price": 10.0},
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "missing", Expression: `VLOOKUP("ZZ", prices, "sku", "price")`},
		{Name: "found", Expression: `VLOOKUP("A1", prices, "sku", "price")`},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (formula-level errors don't abort the call)", err)
	}
	if results["missing"].Err == nil {
		t.Fatal("missing.Err = nil, want a not-found error")
	}
	var fe *apperror.FormulaError
	if !errors.As(results["missing"].Err, &fe) || fe.Code != apperror.ErrLookupNotFound {
		t.Errorf("missing.Err = %v, want ErrLookupNotFound", results["missing"].Err)
	}
	if results["found"].Err != nil {
		t.Errorf("found.Err = %v, want nil", results["found"].Err)
	}
	if results["found"].Value != 10.0 {
		t.Errorf("found.Value = %v, want 10", results["found"].Value)
	}
}

// rfc.md §10: a malformed (mixed-type) Array is ErrArrayTypeMismatch,
// and stays formula-level -- an unrelated sibling formula in the same
// call still succeeds.
func TestEvaluate_ArrayTypeMismatchDoesNotBlockSiblings(t *testing.T) {
	data := map[string]any{
		"mixed": []any{
			map[string]any{"sku": "A1"},
			5.0,
		},
	}

	results, err := Evaluate([]FormulaInput{
		{Name: "broken", Expression: "UNIQUE(mixed)"},
		{Name: "ok", Expression: "5 + 5"},
	}, data)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (formula-level errors don't abort the call)", err)
	}
	if results["broken"].Err == nil {
		t.Fatal("broken.Err = nil, want a mixed-type array error")
	}
	var fe *apperror.FormulaError
	if !errors.As(results["broken"].Err, &fe) || fe.Code != apperror.ErrArrayTypeMismatch {
		t.Errorf("broken.Err = %v, want ErrArrayTypeMismatch", results["broken"].Err)
	}
	if results["ok"].Err != nil {
		t.Errorf("ok.Err = %v, want nil", results["ok"].Err)
	}
	if results["ok"].Value != 10.0 {
		t.Errorf("ok.Value = %v, want 10", results["ok"].Value)
	}
}
