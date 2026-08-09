package math

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(sumFunction{})
}

type sumFunction struct{}

func (sumFunction) Name() string { return "SUM" }
func (sumFunction) MinArgs() int { return 1 }
func (sumFunction) MaxArgs() int { return -1 }

func (f sumFunction) ValidateArgTypes(args []registry.Value) error {
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

func (f sumFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	var total float64
	for _, a := range args {
		v, _ := toFloat64(f.Name(), a)
		total += v
	}
	return total, nil
}
