package logic

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(orFunction{})
}

type orFunction struct{}

func (orFunction) Name() string { return "OR" }
func (orFunction) MinArgs() int { return 1 }
func (orFunction) MaxArgs() int { return -1 }

func (f orFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	for _, a := range args {
		if _, err := toBool(f.Name(), a); err != nil {
			return err
		}
	}
	return nil
}

func (f orFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	for _, a := range args {
		b, _ := toBool(f.Name(), a)
		if b {
			return true, nil
		}
	}
	return false, nil
}
