package formulaengine

import "testing"

// PRD (0002) use case 5: a missing lookup key is a formula-level
// error -- it must not block an unrelated formula in the same call,
// same guarantee feature-0001 already proved for other formula-level
// errors (see formula_partial_success_test.go).
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
		t.Error("missing.Err = nil, want a not-found error")
	}
	if results["found"].Err != nil {
		t.Errorf("found.Err = %v, want nil", results["found"].Err)
	}
	if results["found"].Value != 10.0 {
		t.Errorf("found.Value = %v, want 10", results["found"].Value)
	}
}
