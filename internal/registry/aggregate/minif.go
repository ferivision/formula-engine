package aggregate

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(minifFunction{})
}

type minifFunction struct{}

func (minifFunction) Name() string { return "MINIF" }
func (minifFunction) MinArgs() int { return 4 }
func (minifFunction) MaxArgs() int { return 4 }

func (f minifFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[0]); err != nil {
		return err
	}
	if _, err := toFieldName(f.Name(), args[1]); err != nil {
		return err
	}
	if _, err := toFieldName(f.Name(), args[3]); err != nil {
		return err
	}
	return nil
}

// Evaluate finds the minimum of minField across records where
// conditionField equals conditionValue. Unlike SUMIF/COUNTIF, an
// empty array or zero matching records is a formula-level TypeError,
// per rfc.md §9 -- min of nothing is undefined.
func (f minifFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]
	minField, _ := toFieldName(f.Name(), args[3])

	var min float64
	found := false
	for _, elem := range arr.Elements {
		if isNullElement(elem) {
			continue
		}
		rec, ok := elem.(evaluator.Record)
		if !ok {
			return nil, newTypeError(f.Name() + ": array element is not a record")
		}
		matches, err := valuesEqual(fieldValue(rec, conditionField), conditionValue)
		if err != nil {
			return nil, err
		}
		if !matches {
			continue
		}
		val, err := toFloat64(fieldValue(rec, minField))
		if err != nil {
			return nil, err
		}
		if !found || val < min {
			min = val
			found = true
		}
	}
	if !found {
		return nil, newTypeError(f.Name() + ": no matching records")
	}
	return min, nil
}
