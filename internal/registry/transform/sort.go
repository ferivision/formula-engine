package transform

import (
	"fmt"
	"sort"

	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(sortFunction{})
}

type sortFunction struct{}

func (sortFunction) Name() string { return "SORT" }
func (sortFunction) MinArgs() int { return 3 }
func (sortFunction) MaxArgs() int { return 3 }

func (f sortFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toRecordArray(f.Name(), args[0]); err != nil {
		return err
	}
	if _, err := toFieldName(f.Name(), args[1]); err != nil {
		return err
	}
	_, err := toDirection(f.Name(), args[2])
	return err
}

// Evaluate returns a new Array with records ordered by sortField.
// Uses sort.SliceStable -- never sort.Slice -- so records sharing an
// equal sort key keep their original relative order (PRD NFR-3). The
// input Array is never mutated: a copy of Elements is sorted instead.
func (f sortFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toRecordArray(f.Name(), args[0])
	sortField, _ := toFieldName(f.Name(), args[1])
	direction, _ := toDirection(f.Name(), args[2])

	sorted := make([]evaluator.Value, len(arr.Elements))
	copy(sorted, arr.Elements)

	keyOf := func(elem evaluator.Value) (evaluator.Value, error) {
		if isNullElement(elem) {
			return nil, nil
		}
		rec, ok := elem.(evaluator.Record)
		if !ok {
			return nil, newTypeError(f.Name() + ": array element is not a record")
		}
		return fieldValue(rec, sortField), nil
	}

	var sortErr error
	sort.SliceStable(sorted, func(i, j int) bool {
		if sortErr != nil {
			return false
		}
		ki, err := keyOf(sorted[i])
		if err != nil {
			sortErr = err
			return false
		}
		kj, err := keyOf(sorted[j])
		if err != nil {
			sortErr = err
			return false
		}
		cmp, err := compareValues(ki, kj)
		if err != nil {
			sortErr = err
			return false
		}
		if direction == "desc" {
			return cmp > 0
		}
		return cmp < 0
	})
	if sortErr != nil {
		return nil, sortErr
	}
	return evaluator.Array{IsRecord: true, Elements: sorted}, nil
}

func toDirection(name string, v registry.Value) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", newTypeError(fmt.Sprintf("%s: expected a direction string, got %T", name, v))
	}
	if s != "asc" && s != "desc" {
		return "", newRuntimeError(fmt.Sprintf("%s: unsupported direction %q (want \"asc\" or \"desc\")", name, s))
	}
	return s, nil
}

// compareValues orders two sort keys: numeric-comparable values
// (float64/int/bool/nil) compare numerically; strings compare
// lexicographically. Comparing across those two families (or any
// other type) is a formula-level TypeError -- there's no defined
// ordering between them.
func compareValues(a, b evaluator.Value) (int, error) {
	if af, ok := sortableFloat(a); ok {
		if bf, ok := sortableFloat(b); ok {
			switch {
			case af < bf:
				return -1, nil
			case af > bf:
				return 1, nil
			default:
				return 0, nil
			}
		}
	}
	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok {
			switch {
			case as < bs:
				return -1, nil
			case as > bs:
				return 1, nil
			default:
				return 0, nil
			}
		}
	}
	return 0, newTypeError(fmt.Sprintf("cannot compare %T and %T for sorting", a, b))
}

func sortableFloat(v evaluator.Value) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case bool:
		if n {
			return 1, true
		}
		return 0, true
	case nil:
		return 0, true
	default:
		return 0, false
	}
}
