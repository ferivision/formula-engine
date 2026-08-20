package lookup

import (
	"fmt"
	"strconv"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func newRuntimeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrRuntime, Message: message}
}

func newTypeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrTypeMismatch, Message: message}
}

// newNotFoundError reports a missing lookup key. rfc.md §10: the
// dedicated ErrLookupNotFound, formula-level like the interim
// ErrRuntime it replaced (F2-FASE-7.1 migration).
func newNotFoundError(name string, key registry.Value) error {
	return &apperror.FormulaError{
		Code:    apperror.ErrLookupNotFound,
		Message: fmt.Sprintf("%s: key %v not found", name, key),
	}
}

func checkArgCount(name string, minArgs, maxArgs int, args []registry.Value) error {
	if len(args) < minArgs || (maxArgs != -1 && len(args) > maxArgs) {
		return newRuntimeError(fmt.Sprintf("%s: wrong number of arguments (got %d)", name, len(args)))
	}
	return nil
}

// toArray requires v to be an Array, scalar or Record -- INDEX works
// on either, unlike the rest of this category which needs Records
// for field access.
func toArray(name string, v registry.Value) (evaluator.Array, error) {
	arr, ok := v.(evaluator.Array)
	if !ok {
		return evaluator.Array{}, newTypeError(fmt.Sprintf("%s: expected an array, got %T", name, v))
	}
	return arr, nil
}

// toRecordArray requires v to be an Array of Records -- this category
// reads named fields per element, which only a Record supports. An
// empty array is let through regardless of its IsRecord flag, since
// that flag is meaningless with no elements to have inferred it from.
func toRecordArray(name string, v registry.Value) (evaluator.Array, error) {
	arr, ok := v.(evaluator.Array)
	if !ok {
		return evaluator.Array{}, newTypeError(fmt.Sprintf("%s: expected an array, got %T", name, v))
	}
	if len(arr.Elements) == 0 {
		return arr, nil
	}
	if !arr.IsRecord {
		return evaluator.Array{}, newTypeError(name + ": expected an array of records, got a scalar array")
	}
	return arr, nil
}

func toFieldName(name string, v registry.Value) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", newTypeError(fmt.Sprintf("%s: expected a field name string, got %T", name, v))
	}
	return s, nil
}

// isNullElement reports whether an Array element is a blank slot --
// either a genuinely nil interface or a nil Record -- per rfc.md §9.
func isNullElement(v evaluator.Value) bool {
	if v == nil {
		return true
	}
	rec, ok := v.(evaluator.Record)
	return ok && rec == nil
}

// fieldValue reads field from rec, treating a missing key as nil.
func fieldValue(rec evaluator.Record, field string) evaluator.Value {
	return rec[field]
}

// valuesEqual mirrors internal/registry/comparison's EQUALS
// semantics for key matching: same-typed strings/bools compare
// directly; otherwise both sides are coerced numerically per
// rfc.md (0001) §10.
func valuesEqual(a, b evaluator.Value) (bool, error) {
	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok {
			return as == bs, nil
		}
	}
	if ab, ok := a.(bool); ok {
		if bb, ok := b.(bool); ok {
			return ab == bb, nil
		}
	}
	af, err := toFloat64(a)
	if err != nil {
		return false, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return false, err
	}
	return af == bf, nil
}

func toFloat64(v evaluator.Value) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case bool:
		if n {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return 0, newTypeError(fmt.Sprintf("cannot parse %q as a number", n))
		}
		return f, nil
	default:
		return 0, newTypeError(fmt.Sprintf("expected a number, got %T", v))
	}
}
