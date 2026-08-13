package aggregate

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(countifFunction{})
}

type countifFunction struct{}

func (countifFunction) Name() string { return "COUNTIF" }
func (countifFunction) MinArgs() int { return 3 }
func (countifFunction) MaxArgs() int { return 3 }

func (f countifFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[0]); err != nil {
		return err
	}
	_, err := toFieldName(f.Name(), args[1])
	return err
}

// Evaluate counts records where conditionField equals conditionValue.
// An empty array counts to 0 (rfc.md §9); a null element is skipped,
// not counted.
func (f countifFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]

	var count float64
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
			count++
		}
	}
	return count, nil
}
