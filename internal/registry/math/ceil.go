package math

import (
	"math"

	"github.com/Ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(ceilFunction{})
}

type ceilFunction struct{}

func (ceilFunction) Name() string { return "CEIL" }
func (ceilFunction) MinArgs() int { return 1 }
func (ceilFunction) MaxArgs() int { return 1 }

func (f ceilFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toFloat64(f.Name(), args[0])
	return err
}

func (f ceilFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	value, _ := toFloat64(f.Name(), args[0])
	return math.Ceil(value), nil
}
