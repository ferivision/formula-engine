package comparison

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(equalsFunction{})
}

type equalsFunction struct{}

func (equalsFunction) Name() string { return "EQUALS" }
func (equalsFunction) MinArgs() int { return 2 }
func (equalsFunction) MaxArgs() int { return 2 }

func (f equalsFunction) ValidateArgTypes(args []registry.Value) error {
	return checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args)
}

// Evaluate compares args[0] and args[1] for equality. Same-typed
// strings or bools compare directly; anything else is coerced
// numerically per rfc.md §10 (numbers, bools, nil, and numeric
// strings all participate in that coercion) -- a non-numeric string
// compared against a number is a formula-level TypeError, the same
// rule arithmetic already applies to "number + string".
func (f equalsFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	a, b := args[0], args[1]

	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok {
			return as == bs, nil
		}
	}
	if ab, ok := a.(bool); ok {
		if bb, ok := b.(bool); ok {
			return ab == bb, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return nil, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return nil, err
	}
	return af == bf, nil
}
