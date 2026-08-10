package formulaengine

import "testing"

// rfc.md §7: a formula-level runtime error must not block an
// independent sibling in the same call.
func TestEvaluate_IndependentFormulaFailureDoesNotBlockSiblings(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "a", Expression: "1 / 0"},
		{Name: "b", Expression: "5 + 5"},
	}, nil)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (formula-level errors don't abort the call)", err)
	}
	if results["a"].Err == nil {
		t.Error("a.Err = nil, want a division-by-zero error")
	}
	if results["b"].Err != nil {
		t.Errorf("b.Err = %v, want nil", results["b"].Err)
	}
	if results["b"].Value != 10.0 {
		t.Errorf("b.Value = %v, want 10", results["b"].Value)
	}
}
