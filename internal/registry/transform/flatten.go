package transform

import (
	"fmt"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(flattenFunction{})
}

type flattenFunction struct{}

func (flattenFunction) Name() string { return "FLATTEN" }
func (flattenFunction) MinArgs() int { return 1 }
func (flattenFunction) MaxArgs() int { return 1 }

func (f flattenFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toArray(f.Name(), args[0])
	return err
}

// Evaluate concatenates an array of nested arrays into one flat
// Array -- the one deliberate exception to "no nested arrays"
// (rfc.md §7), since it consumes pre-existing nesting rather than
// producing or navigating it. Each element of the outer array must
// itself be array-shaped ([]any, []map[string]any, or an
// already-converted Array); nesting sub-arrays that disagree on
// record-vs-scalar is a formula-level TypeError, same as any other
// mixed-type Array.
//
// Known gap: a concretely-typed [][]any input is not recognized here
// (or anywhere in NewArray) -- wrap it as []any{...} instead.
func (f flattenFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	outer, _ := toArray(f.Name(), args[0])

	var result []evaluator.Value
	isRecord := false
	seenAny := false

	for _, elem := range outer.Elements {
		if isNullElement(elem) {
			continue
		}
		inner, err := toInnerArray(f.Name(), elem)
		if err != nil {
			return nil, err
		}
		if seenAny && inner.IsRecord != isRecord {
			return nil, newTypeError(f.Name() + ": cannot flatten arrays mixing records and scalars")
		}
		isRecord = inner.IsRecord
		seenAny = true
		result = append(result, inner.Elements...)
	}
	return evaluator.Array{IsRecord: isRecord, Elements: result}, nil
}

func toInnerArray(name string, v evaluator.Value) (evaluator.Array, error) {
	if arr, ok := v.(evaluator.Array); ok {
		return arr, nil
	}
	switch v.(type) {
	case []any, []map[string]any:
		return evaluator.NewArray(v)
	default:
		return evaluator.Array{}, newTypeError(fmt.Sprintf("%s: expected a nested array, got %T", name, v))
	}
}
