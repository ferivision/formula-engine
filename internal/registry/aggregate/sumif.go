package aggregate

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(sumifFunction{})
}

type sumifFunction struct{}

func (sumifFunction) Name() string { return "SUMIF" }
func (sumifFunction) MinArgs() int { return 4 }
func (sumifFunction) MaxArgs() int { return 4 }

func (f sumifFunction) ValidateArgTypes(args []registry.Value) error {
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

// Evaluate sums sumField across records where conditionField equals
// conditionValue. An empty array sums to 0 (rfc.md §9); a null
// element is skipped, not summed.
func (f sumifFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]
	sumField, _ := toFieldName(f.Name(), args[3])

	var total float64
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
		sumVal, err := toFloat64(fieldValue(rec, sumField))
		if err != nil {
			return nil, err
		}
		total += sumVal
	}
	return total, nil
}
