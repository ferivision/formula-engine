package math

import (
	"math"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(floorFunction{})
}

type floorFunction struct{}

func (floorFunction) Name() string { return "FLOOR" }
func (floorFunction) MinArgs() int { return 1 }
func (floorFunction) MaxArgs() int { return 1 }

func (f floorFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toFloat64(f.Name(), args[0])
	return err
}

func (f floorFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	value, _ := toFloat64(f.Name(), args[0])
	return math.Floor(value), nil
}
