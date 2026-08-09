package logic

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(andFunction{})
}

type andFunction struct{}

func (andFunction) Name() string { return "AND" }
func (andFunction) MinArgs() int { return 1 }
func (andFunction) MaxArgs() int { return -1 }

func (f andFunction) ValidateArgTypes(args []registry.Value) error {
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

func (f andFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	for _, a := range args {
		b, _ := toBool(f.Name(), a)
		if !b {
			return false, nil
		}
	}
	return true, nil
}
