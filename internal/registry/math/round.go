package math

import (
	"math"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(roundFunction{})
}

type roundFunction struct{}

func (roundFunction) Name() string { return "ROUND" }
func (roundFunction) MinArgs() int { return 1 }
func (roundFunction) MaxArgs() int { return 2 }

func (f roundFunction) ValidateArgTypes(args []registry.Value) error {
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

func (f roundFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	value, _ := toFloat64(f.Name(), args[0])
	digits := 0.0
	if len(args) == 2 {
		digits, _ = toFloat64(f.Name(), args[1])
	}
	mult := math.Pow(10, digits)
	return math.Round(value*mult) / mult, nil
}
