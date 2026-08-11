package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/parser"
)

func mustEvaluate(t *testing.T, input string) any {
	t.Helper()
	return mustEvaluateWithData(t, input, nil)
}

func mustEvaluateWithData(t *testing.T, input string, data map[string]any) any {
	t.Helper()
	tokens, err := parser.Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", input, err)
	}
	got, err := Evaluate(node, NewContext(data))
	if err != nil {
		t.Fatalf("Evaluate(%q) error = %v", input, err)
	}
	return got
}

func TestEvaluate_NumberLiteral(t *testing.T) {
	got := mustEvaluate(t, "42")
	if got != 42.0 {
		t.Errorf("Evaluate() = %v, want 42", got)
	}
}

func TestEvaluate_PrecedenceRespected(t *testing.T) {
	got := mustEvaluate(t, "1 + 2 * 3")
	if got != 7.0 {
		t.Errorf("Evaluate() = %v, want 7", got)
	}
}

func TestEvaluate_ParenthesesRespected(t *testing.T) {
	got := mustEvaluate(t, "(1 + 2) * 3")
	if got != 9.0 {
		t.Errorf("Evaluate() = %v, want 9", got)
	}
}

func TestEvaluate_Subtraction(t *testing.T) {
	got := mustEvaluate(t, "10 - 3 - 2")
	if got != 5.0 {
		t.Errorf("Evaluate() = %v, want 5", got)
	}
}

func TestEvaluate_Division(t *testing.T) {
	got := mustEvaluate(t, "10 / 4")
	if got != 2.5 {
		t.Errorf("Evaluate() = %v, want 2.5", got)
	}
}

func TestEvaluate_DivisionByZero(t *testing.T) {
	tokens, _ := parser.Tokenize("1 / 0")
	node, _ := parser.Parse(tokens)

	_, err := Evaluate(node, NewContext(nil))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want division-by-zero error")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("error = %v, want ErrRuntime", err)
	}
}
