package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/parser"
	_ "github.com/ferivision/formula-engine/internal/registry/math"
)

func TestEvaluate_FunctionCallWithArithmeticArgs(t *testing.T) {
	got := mustEvaluate(t, "MAX(1 + 2, 4 * 5)")
	if got != 20.0 {
		t.Errorf("Evaluate() = %v, want 20", got)
	}
}

func TestEvaluate_NestedFunctionCalls(t *testing.T) {
	got := mustEvaluate(t, "MAX(1, MIN(2, 3))")
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2", got)
	}
}

func TestEvaluate_FunctionCallWithinArithmetic(t *testing.T) {
	got := mustEvaluate(t, "MAX(1, 2) * 3")
	if got != 6.0 {
		t.Errorf("Evaluate() = %v, want 6", got)
	}
}

func TestEvaluate_UnknownFunction(t *testing.T) {
	tokens, err := parser.Tokenize("NOPE(1)")
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	_, err = Evaluate(node, NewContext(nil))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want error for unknown function")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) {
		t.Fatalf("error = %v, want *apperror.FormulaError", err)
	}
}
