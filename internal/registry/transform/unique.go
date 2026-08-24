package transform

import (
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(uniqueFunction{})
}

type uniqueFunction struct{}

func (uniqueFunction) Name() string { return "UNIQUE" }
func (uniqueFunction) MinArgs() int { return 1 }
func (uniqueFunction) MaxArgs() int { return 2 }

func (f uniqueFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toArray(f.Name(), args[0]); err != nil {
		return err
	}
	if len(args) == 2 {
		_, err := toFieldName(f.Name(), args[1])
		return err
	}
	return nil
}

// Evaluate deduplicates array, preserving first-seen order (PRD use
// case 3). UNIQUE(array) dedupes a scalar array by direct value
// equality (no numeric coercion -- 1 and "1" are different values
// here, unlike FILTER's condition matching). UNIQUE(array, field)
// dedupes a record array by a field's value instead. The input Array
// is never mutated -- a new Elements slice is built from scratch.
func (f uniqueFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toArray(f.Name(), args[0])

	if len(args) == 2 {
		field, _ := toFieldName(f.Name(), args[1])
		return uniqueByField(f.Name(), arr, field)
	}
	if len(arr.Elements) > 0 && arr.IsRecord {
		return nil, newTypeError(f.Name() + ": array of records requires a field argument")
	}
	return uniqueScalars(arr), nil
}

func uniqueScalars(arr evaluator.Array) evaluator.Value {
	seen := make(map[any]bool, len(arr.Elements))
	var result []evaluator.Value
	for _, elem := range arr.Elements {
		if isNullElement(elem) || seen[elem] {
			continue
		}
		seen[elem] = true
		result = append(result, elem)
	}
	return evaluator.Array{IsRecord: false, Elements: result}
}

func uniqueByField(name string, arr evaluator.Array, field string) (evaluator.Value, error) {
	if len(arr.Elements) > 0 && !arr.IsRecord {
		return nil, newTypeError(name + ": expected an array of records, got a scalar array")
	}
	seen := make(map[any]bool, len(arr.Elements))
	var result []evaluator.Value
	for _, elem := range arr.Elements {
		if isNullElement(elem) {
			continue
		}
		rec, ok := elem.(evaluator.Record)
		if !ok {
			return nil, newTypeError(name + ": array element is not a record")
		}
		key := fieldValue(rec, field)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, rec)
	}
	return evaluator.Array{IsRecord: true, Elements: result}, nil
}
