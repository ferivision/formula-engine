package lookup

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(matchFunction{})
}

type matchFunction struct{}

func (matchFunction) Name() string { return "MATCH" }
func (matchFunction) MinArgs() int { return 3 }
func (matchFunction) MaxArgs() int { return 3 }

func (f matchFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[1]); err != nil {
		return err
	}
	_, err := toFieldName(f.Name(), args[2])
	return err
}

// Evaluate returns the 1-based position of the first record in array
// where field equals key. Positions are based on the array's actual
// indices (a skipped null element still occupies its own position),
// not a count of records considered. A key not found is the same
// formula-level not-found treatment as VLOOKUP.
func (f matchFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	key := args[0]
	arr, _ := toRecordArray(f.Name(), args[1])
	field, _ := toFieldName(f.Name(), args[2])

	for i, elem := range arr.Elements {
		if isNullElement(elem) {
			continue
		}
		rec, ok := elem.(evaluator.Record)
		if !ok {
			return nil, newTypeError(f.Name() + ": array element is not a record")
		}
		matches, err := valuesEqual(fieldValue(rec, field), key)
		if err != nil {
			return nil, err
		}
		if matches {
			return float64(i + 1), nil
		}
	}
	return nil, newNotFoundError(f.Name(), key)
}
