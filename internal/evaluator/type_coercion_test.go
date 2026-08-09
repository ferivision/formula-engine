package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/parser"
)

func parseHelper(t *testing.T, input string) *parser.Node {
	t.Helper()
	tokens, err := parser.Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", input, err)
	}
	return node
}

// Each test below corresponds to one row of rfc.md §10's coercion table.

func TestCoercion_NumberPlusNumber(t *testing.T) {
	got := mustEvaluate(t, "2 + 3")
	if got != 5.0 {
		t.Errorf("Evaluate() = %v, want 5", got)
	}
}

func TestCoercion_StringPlusStringConcatenates(t *testing.T) {
	got := mustEvaluate(t, `"foo" + "bar"`)
	if got != "foobar" {
		t.Errorf("Evaluate() = %v, want \"foobar\"", got)
	}
}

func TestCoercion_NumberPlusString_ParsesNumericString(t *testing.T) {
	got := mustEvaluate(t, `1 + "2"`)
	if got != 3.0 {
		t.Errorf("Evaluate() = %v, want 3", got)
	}
}

func TestCoercion_NumberPlusString_ErrorsIfNotParseable(t *testing.T) {
	node := parseHelper(t, `1 + "abc"`)
	_, err := Evaluate(node, NewContext(nil))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-numeric string")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("error = %v, want ErrTypeMismatch", err)
	}
}

func TestCoercion_BoolInNumericContext(t *testing.T) {
	got := mustEvaluateWithData(t, "active + 1", map[string]any{"active": true})
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2 (true -> 1)", got)
	}

	got = mustEvaluateWithData(t, "active + 1", map[string]any{"active": false})
	if got != 1.0 {
		t.Errorf("Evaluate() = %v, want 1 (false -> 0)", got)
	}
}

func TestCoercion_NullInArithmetic_TreatedAsZeroForPlusMinus(t *testing.T) {
	got := mustEvaluateWithData(t, "x + 1", map[string]any{"x": nil})
	if got != 1.0 {
		t.Errorf("Evaluate() = %v, want 1 (null -> 0)", got)
	}

	got = mustEvaluateWithData(t, "5 - x", map[string]any{"x": nil})
	if got != 5.0 {
		t.Errorf("Evaluate() = %v, want 5 (null -> 0)", got)
	}
}

func TestCoercion_DivisionByNull_Errors(t *testing.T) {
	node := parseHelper(t, "5 / x")
	_, err := Evaluate(node, NewContext(map[string]any{"x": nil}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want division error for null divisor")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("error = %v, want ErrRuntime", err)
	}
}

func TestCoercion_NullInTextContext_TreatedAsEmptyString(t *testing.T) {
	got, err := coerceForText(nil)
	if err != nil || got != "" {
		t.Errorf("coerceForText(nil) = (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestCoercion_MismatchedTypeWithNoDefinedRule(t *testing.T) {
	_, err := coerceForArithmetic([]int{1, 2})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("error = %v, want ErrTypeMismatch", err)
	}
}
