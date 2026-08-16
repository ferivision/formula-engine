package lookup

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(vlookupFunction{})
}

type vlookupFunction struct{}

func (vlookupFunction) Name() string { return "VLOOKUP" }
func (vlookupFunction) MinArgs() int { return 4 }
func (vlookupFunction) MaxArgs() int { return 4 }

func (f vlookupFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[1]); err != nil {
		return err
	}
	if _, err := toFieldName(f.Name(), args[2]); err != nil {
		return err
	}
	_, err := toFieldName(f.Name(), args[3])
	return err
}

// Evaluate finds the first record in table where keyField equals key,
// and returns that record's returnField value. A key not found is a
// formula-level error (interim ErrRuntime, see support.go's
// newNotFoundError) -- per PRD use case 5, it must not block
// unrelated formulas in the same Evaluate call.
func (f vlookupFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	key := args[0]
	table, _ := toRecordArray(f.Name(), args[1])
	keyField, _ := toFieldName(f.Name(), args[2])
	returnField, _ := toFieldName(f.Name(), args[3])

	for _, elem := range table.Elements {
		if isNullElement(elem) {
			continue
		}
		rec, ok := elem.(evaluator.Record)
		if !ok {
			return nil, newTypeError(f.Name() + ": array element is not a record")
		}
		matches, err := valuesEqual(fieldValue(rec, keyField), key)
		if err != nil {
			return nil, err
		}
		if matches {
			return fieldValue(rec, returnField), nil
		}
	}
	return nil, newNotFoundError(f.Name(), key)
}
