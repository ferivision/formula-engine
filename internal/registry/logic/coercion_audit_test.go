package logic

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Audits AND/OR/NOT against rfc.md §10: real booleans work, and
// number/string/null -- deliberately, per support.go's toBool -- have
// no defined boolean-context coercion and stay a TypeError rather
// than inventing numeric truthiness. IF's condition coercion is
// audited alongside internal/evaluator (see evaluator_if_test.go),
// since IF's short-circuit logic lives there, not in this package.
func TestCoercionAudit_LogicFunctions_RealBoolWorks(t *testing.T) {
	fns := []registry.Function{andFunction{}, orFunction{}, notFunction{}}
	for _, fn := range fns {
		t.Run(fn.Name(), func(t *testing.T) {
			args := make([]registry.Value, fn.MinArgs())
			for i := range args {
				args[i] = true
			}
			if _, err := fn.Evaluate(args); err != nil {
				t.Errorf("Evaluate(true...) error = %v, want nil", err)
			}
		})
	}
}

func TestCoercionAudit_LogicFunctions_NumberArgErrors(t *testing.T) {
	fns := []registry.Function{andFunction{}, orFunction{}, notFunction{}}
	for _, fn := range fns {
		t.Run(fn.Name(), func(t *testing.T) {
			args := make([]registry.Value, fn.MinArgs())
			for i := range args {
				args[i] = 1.0
			}
			if _, err := fn.Evaluate(args); err == nil {
				t.Errorf("Evaluate(1.0...) error = nil, want type error -- rfc.md §10 defines no numeric-truthiness coercion")
			}
		})
	}
}
