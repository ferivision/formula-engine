package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/parser"
)

func TestEvaluate_FieldReference(t *testing.T) {
	got := mustEvaluateWithData(t, "price * quantity", map[string]any{
		"price":    10.0,
		"quantity": 2.0,
	})
	if got != 20.0 {
		t.Errorf("Evaluate() = %v, want 20", got)
	}
}

func TestEvaluate_UndefinedFieldReference(t *testing.T) {
	tokens, err := parser.Tokenize("price * quantity")
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	_, err = Evaluate(node, NewContext(map[string]any{"price": 10.0}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want undefined-reference error for missing \"quantity\"")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrUndefinedReference {
		t.Errorf("error = %v, want ErrUndefinedReference", err)
	}
}
