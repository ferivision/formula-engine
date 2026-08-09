package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	_ "github.com/ferivision/formula-engine/internal/registry/logic"
)

func TestEvaluate_If_TrueBranchShortCircuits(t *testing.T) {
	// The false-branch (1/0) must never be evaluated -- if it were,
	// this would error instead of returning 1.
	got := mustEvaluateWithData(t, "IF(cond, 1, 1/0)", map[string]any{"cond": true})
	if got != 1.0 {
		t.Errorf("Evaluate() = %v, want 1", got)
	}
}

func TestEvaluate_If_FalseBranchShortCircuits(t *testing.T) {
	// The true-branch (1/0) must never be evaluated -- if it were,
	// this would error instead of returning 2.
	got := mustEvaluateWithData(t, "IF(cond, 1/0, 2)", map[string]any{"cond": false})
	if got != 2.0 {
		t.Errorf("Evaluate() = %v, want 2", got)
	}
}

func TestEvaluate_If_NonBoolConditionErrors(t *testing.T) {
	node := parseHelper(t, "IF(cond, 1, 2)")
	_, err := Evaluate(node, NewContext(map[string]any{"cond": 5.0}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want type error for non-bool condition")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("error = %v, want ErrTypeMismatch", err)
	}
}

func TestEvaluate_If_WrongArgCount(t *testing.T) {
	node := parseHelper(t, "IF(cond, 1)")
	_, err := Evaluate(node, NewContext(map[string]any{"cond": true}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want error for 2 args (IF needs 3)")
	}
}

func TestEvaluate_If_NestedInArithmetic(t *testing.T) {
	got := mustEvaluateWithData(t, "IF(cond, 1, 2) + 10", map[string]any{"cond": true})
	if got != 11.0 {
		t.Errorf("Evaluate() = %v, want 11", got)
	}
}
