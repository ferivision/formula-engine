package text

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Audits every text function against rfc.md §10's text-context row:
// null becomes "", a string passes through, and -- deliberately, per
// support.go's toText -- number/bool have no defined text-context
// coercion and stay a TypeError rather than silently formatting.
func TestCoercionAudit_TextFunctions_NullBecomesEmptyString(t *testing.T) {
	fns := []registry.Function{
		concatFunction{}, upperFunction{}, lowerFunction{}, trimFunction{}, lengthFunction{},
	}
	for _, fn := range fns {
		t.Run(fn.Name(), func(t *testing.T) {
			args := make([]registry.Value, fn.MinArgs())
			for i := range args {
				args[i] = nil
			}
			if _, err := fn.Evaluate(args); err != nil {
				t.Errorf("Evaluate(nil...) error = %v, want nil (null -> \"\" per rfc.md §10)", err)
			}
		})
	}
}

func TestCoercionAudit_TextFunctions_NumberArgErrors(t *testing.T) {
	fns := []registry.Function{
		concatFunction{}, upperFunction{}, lowerFunction{}, trimFunction{}, lengthFunction{},
	}
	for _, fn := range fns {
		t.Run(fn.Name(), func(t *testing.T) {
			args := make([]registry.Value, fn.MinArgs())
			for i := range args {
				args[i] = 5.0
			}
			if _, err := fn.Evaluate(args); err == nil {
				t.Errorf("Evaluate(5.0...) error = nil, want type error -- rfc.md §10 defines no number-to-text coercion")
			}
		})
	}
}
