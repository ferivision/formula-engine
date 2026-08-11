package logic

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(ifFunction{})
}

// ifFunction is a registry marker for IF's name and arity only. IF
// needs its false/true branches evaluated lazily (short-circuit), but
// registry.Function.Evaluate receives already-evaluated args -- that
// signature can't support it. internal/evaluator special-cases IF at
// the AST level instead and never calls Evaluate below; this type
// exists purely so IF's arity is discoverable/validated the same way
// as every other function.
type ifFunction struct{}

func (ifFunction) Name() string { return "IF" }
func (ifFunction) MinArgs() int { return 3 }
func (ifFunction) MaxArgs() int { return 3 }

func (ifFunction) ValidateArgTypes(args []registry.Value) error { return nil }

func (ifFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	panic("logic.ifFunction.Evaluate should never be called -- IF is special-cased in internal/evaluator for short-circuit evaluation")
}
