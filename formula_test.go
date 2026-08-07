package formulaengine

import "testing"

func TestEvaluate_StubReturnsNotImplementedError(t *testing.T) {
	results, err := Evaluate([]FormulaInput{{Name: "total", Expression: "1 + 1"}}, map[string]any{})

	if err == nil {
		t.Fatal("Evaluate() error = nil, want a not-implemented error")
	}
	if results != nil {
		t.Errorf("Evaluate() results = %v, want nil", results)
	}
}
