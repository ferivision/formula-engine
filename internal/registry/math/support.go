package math

import (
	"fmt"

	"github.com/Ferivision/formula-engine/internal/apperror"
	"github.com/Ferivision/formula-engine/internal/registry"
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

func toFloat64(name string, v registry.Value) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	default:
		return 0, &apperror.FormulaError{
			Code:    apperror.ErrTypeMismatch,
			Message: fmt.Sprintf("%s: expected a number, got %T", name, v),
		}
	}
}
