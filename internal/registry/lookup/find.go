package lookup

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(findFunction{})
}

type findFunction struct{}

func (findFunction) Name() string { return "FIND" }
func (findFunction) MinArgs() int { return 3 }
func (findFunction) MaxArgs() int { return 3 }

func (f findFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[0]); err != nil {
		return err
	}
	_, err := toFieldName(f.Name(), args[1])
	return err
}

// Evaluate returns the first record where conditionField equals
// conditionValue. A key not found is the same formula-level
// not-found treatment as VLOOKUP/MATCH.
func (f findFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]

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
			return rec, nil
		}
	}
	return nil, newNotFoundError(f.Name(), conditionValue)
}
