package math

import (
	"math"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(absFunction{})
}

type absFunction struct{}

func (absFunction) Name() string { return "ABS" }
func (absFunction) MinArgs() int { return 1 }
func (absFunction) MaxArgs() int { return 1 }

func (f absFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toFloat64(f.Name(), args[0])
	return err
}

func (f absFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	value, _ := toFloat64(f.Name(), args[0])
	return math.Abs(value), nil
}
