package math

import "github.com/Ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(minFunction{})
}

type minFunction struct{}

func (minFunction) Name() string { return "MIN" }
func (minFunction) MinArgs() int { return 1 }
func (minFunction) MaxArgs() int { return -1 }

func (f minFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	for _, a := range args {
		if _, err := toFloat64(f.Name(), a); err != nil {
			return err
		}
	}
	return nil
}

func (f minFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	min, _ := toFloat64(f.Name(), args[0])
	for _, a := range args[1:] {
		v, _ := toFloat64(f.Name(), a)
		if v < min {
			min = v
		}
	}
	return min, nil
}
