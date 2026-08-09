package logic

import (
	"fmt"

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

// toBool requires a real bool -- rfc.md §10 has no defined coercion
// from number/string/null into a boolean, so anything else falls
// through to the table's catch-all: a formula-level TypeError.
func toBool(name string, v registry.Value) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, &apperror.FormulaError{
			Code:    apperror.ErrTypeMismatch,
			Message: fmt.Sprintf("%s: expected a boolean, got %T", name, v),
		}
	}
	return b, nil
}
