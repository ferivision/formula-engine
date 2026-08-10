package math

import (
	"fmt"
	"strconv"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func checkArgCount(name string, minArgs, maxArgs int, args []registry.Value) error {
	if len(args) < minArgs || (maxArgs != -1 && len(args) > maxArgs) {
		return &apperror.FormulaError{
			Code:    apperror.ErrRuntime,
			Message: fmt.Sprintf("%s: wrong number of arguments (got %d)", name, len(args)),
		}
	}
	return nil
}

// toFloat64 mirrors rfc.md §10's arithmetic coercion rules -- the
// same ones internal/evaluator's coerceForArithmetic applies -- so
// that e.g. MAX(x, 1) and "x + 1" agree on what x coerces to. Numbers
// pass through, bools become 1/0, nil becomes 0, and numeric strings
// are parsed; int is additionally accepted as a defensive alternate
// representation of "a number" for callers whose data map holds a
// plain Go int rather than float64 (not itself a §10 table row).
func toFloat64(name string, v registry.Value) (float64, error) {
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
			return 0, &apperror.FormulaError{
				Code:    apperror.ErrTypeMismatch,
				Message: fmt.Sprintf("%s: cannot parse %q as a number", name, n),
			}
		}
		return f, nil
	default:
		return 0, &apperror.FormulaError{
			Code:    apperror.ErrTypeMismatch,
			Message: fmt.Sprintf("%s: expected a number, got %T", name, v),
		}
	}
}
