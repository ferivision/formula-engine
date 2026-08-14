package transform

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(filterFunction{})
}

type filterFunction struct{}

func (filterFunction) Name() string { return "FILTER" }
func (filterFunction) MinArgs() int { return 3 }
func (filterFunction) MaxArgs() int { return 3 }

func (f filterFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[0]); err != nil {
		return err
	}
	_, err := toFieldName(f.Name(), args[1])
	return err
}

// Evaluate returns a new Array containing only the records where
// conditionField equals conditionValue. A null element can't be
// meaningfully matched against a condition, so it's excluded from
// the result rather than erroring. The input Array is never mutated
// -- a new Elements slice is built from scratch.
func (f filterFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]

	var matched []evaluator.Value
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
		if matches {
			matched = append(matched, rec)
		}
	}
	return evaluator.Array{IsRecord: true, Elements: matched}, nil
}
