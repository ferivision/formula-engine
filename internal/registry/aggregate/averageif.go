package aggregate

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(averageifFunction{})
}

type averageifFunction struct{}

func (averageifFunction) Name() string { return "AVERAGEIF" }
func (averageifFunction) MinArgs() int { return 4 }
func (averageifFunction) MaxArgs() int { return 4 }

func (f averageifFunction) ValidateArgTypes(args []registry.Value) error {
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

// Evaluate averages avgField across records where conditionField
// equals conditionValue. Unlike SUMIF/COUNTIF, an empty array or zero
// matching records is a formula-level TypeError, per rfc.md §9 --
// average of nothing is undefined, not 0.
func (f averageifFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	conditionField, _ := toFieldName(f.Name(), args[1])
	conditionValue := args[2]
	avgField, _ := toFieldName(f.Name(), args[3])

	var total float64
	var count int
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
		val, err := toFloat64(fieldValue(rec, avgField))
		if err != nil {
			return nil, err
		}
		total += val
		count++
	}
	if count == 0 {
		return nil, newTypeError(f.Name() + ": no matching records to average")
	}
	return total / float64(count), nil
}
