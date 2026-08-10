package formulaengine

import (
	"errors"
	"strings"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

// rfc.md §7: when a formula depends on one that failed at runtime,
// it must inherit that failure with a clear "dependency failed"
// message -- not be silently re-evaluated with a missing value, and
// not be misreported as ErrUndefinedReference (a is defined, it just
// failed at runtime; that's a different situation than a not
// existing at all).
func TestEvaluate_DependentFormulaInheritsFailure(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "a", Expression: "1 / 0"},
		{Name: "b", Expression: "a + 1"},
	}, nil)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil (this is formula-level, not call-level)", err)
	}
	if results["a"].Err == nil {
		t.Fatal("a.Err = nil, want a division-by-zero error")
	}
	if results["b"].Err == nil {
		t.Fatal("b.Err = nil, want an inherited-failure error")
	}

	var fe *apperror.FormulaError
	if !errors.As(results["b"].Err, &fe) {
		t.Fatalf("b.Err = %v, want *apperror.FormulaError", results["b"].Err)
	}
	if fe.Code == apperror.ErrUndefinedReference {
		t.Errorf("b.Err code = %v, want NOT ErrUndefinedReference -- %q is defined, it just failed at runtime", fe.Code, "a")
	}

	// The message must distinguish "this formula failed" (a's own
	// error text) from "a dependency of this formula failed" (b's
	// message wrapping a's).
	if !strings.Contains(fe.Message, "a") || !strings.Contains(fe.Message, "depends on") {
		t.Errorf("b.Err = %q, want it to say it depends on the failed formula %q", fe.Message, "a")
	}
	if !strings.Contains(fe.Message, results["a"].Err.Error()) {
		t.Errorf("b.Err = %q, want it to include a's own error (%q)", fe.Message, results["a"].Err.Error())
	}
}

// The inherited failure must cascade transitively: c depends on b,
// which depends on the formula that actually failed.
func TestEvaluate_TransitiveInheritedFailure(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "a", Expression: "1 / 0"},
		{Name: "b", Expression: "a + 1"},
		{Name: "c", Expression: "b + 1"},
	}, nil)

	if err != nil {
		t.Fatalf("Evaluate() error = %v, want nil", err)
	}
	if results["c"].Err == nil {
		t.Fatal("c.Err = nil, want a cascading inherited-failure error")
	}
}
