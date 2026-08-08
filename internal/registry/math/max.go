package math

import "github.com/Ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(maxFunction{})
}

type maxFunction struct{}

func (maxFunction) Name() string { return "MAX" }
func (maxFunction) MinArgs() int { return 1 }
func (maxFunction) MaxArgs() int { return -1 }

func (f maxFunction) ValidateArgTypes(args []registry.Value) error {
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

func (f maxFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	max, _ := toFloat64(f.Name(), args[0])
	for _, a := range args[1:] {
		v, _ := toFloat64(f.Name(), a)
		if v > max {
			max = v
		}
	}
	return max, nil
}
