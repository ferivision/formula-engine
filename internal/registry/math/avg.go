package math

import "github.com/Ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(avgFunction{})
}

type avgFunction struct{}

func (avgFunction) Name() string { return "AVG" }
func (avgFunction) MinArgs() int { return 1 }
func (avgFunction) MaxArgs() int { return -1 }

func (f avgFunction) ValidateArgTypes(args []registry.Value) error {
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

func (f avgFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	var total float64
	for _, a := range args {
		v, _ := toFloat64(f.Name(), a)
		total += v
	}
	return total / float64(len(args)), nil
}
