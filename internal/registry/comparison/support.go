package comparison

import (
	"fmt"
	"strconv"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func newRuntimeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrRuntime, Message: message}
}

func newTypeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrTypeMismatch, Message: message}
}

func checkArgCount(name string, minArgs, maxArgs int, args []registry.Value) error {
	if len(args) < minArgs || (maxArgs != -1 && len(args) > maxArgs) {
		return newRuntimeError(fmt.Sprintf("%s: wrong number of arguments (got %d)", name, len(args)))
	}
	return nil
}

// toFloat64 mirrors rfc.md §10's arithmetic coercion rules: numbers
// pass through, bools become 1/0, nil becomes 0, and numeric strings
// are parsed. Used for numeric-context comparisons (BETWEEN, and
// EQUALS whenever its operands aren't both the same non-numeric
// type). int is additionally accepted as a defensive alternate
// representation of "a number" (not itself a §10 row), kept
// consistent with the same allowance in internal/evaluator and
// internal/registry/math (FASE-13.1).
func toFloat64(v registry.Value) (float64, error) {
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
