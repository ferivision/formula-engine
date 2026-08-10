package text

import (
	"strings"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(upperFunction{})
}

type upperFunction struct{}

func (upperFunction) Name() string { return "UPPER" }
func (upperFunction) MinArgs() int { return 1 }
func (upperFunction) MaxArgs() int { return 1 }

func (f upperFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toText(f.Name(), args[0])
	return err
}

func (f upperFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	s, _ := toText(f.Name(), args[0])
	return strings.ToUpper(s), nil
}
