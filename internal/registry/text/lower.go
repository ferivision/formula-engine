package text

import (
	"strings"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(lowerFunction{})
}

type lowerFunction struct{}

func (lowerFunction) Name() string { return "LOWER" }
func (lowerFunction) MinArgs() int { return 1 }
func (lowerFunction) MaxArgs() int { return 1 }

func (f lowerFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toText(f.Name(), args[0])
	return err
}

func (f lowerFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	s, _ := toText(f.Name(), args[0])
	return strings.ToLower(s), nil
}
