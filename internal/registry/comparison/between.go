package comparison

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(betweenFunction{})
}

type betweenFunction struct{}

func (betweenFunction) Name() string { return "BETWEEN" }
func (betweenFunction) MinArgs() int { return 3 }
func (betweenFunction) MaxArgs() int { return 3 }

func (f betweenFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	for _, a := range args {
		if _, err := toFloat64(a); err != nil {
			return err
		}
	}
	return nil
}

// Evaluate reports whether args[0] falls between args[1] (low) and
// args[2] (high), inclusive of both bounds.
func (f betweenFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	value, _ := toFloat64(args[0])
	low, _ := toFloat64(args[1])
	high, _ := toFloat64(args[2])
	return value >= low && value <= high, nil
}
