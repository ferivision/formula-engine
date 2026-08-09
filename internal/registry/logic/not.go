package logic

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(notFunction{})
}

type notFunction struct{}

func (notFunction) Name() string { return "NOT" }
func (notFunction) MinArgs() int { return 1 }
func (notFunction) MaxArgs() int { return 1 }

func (f notFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toBool(f.Name(), args[0])
	return err
}

func (f notFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	b, _ := toBool(f.Name(), args[0])
	return !b, nil
}
