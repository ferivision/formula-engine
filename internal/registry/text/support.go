package text

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

// toText applies rfc.md §10's text-context rule: a missing (nil)
// value becomes an empty string, a string passes through unchanged.
// Non-string, non-nil values have no defined text-context rule yet
// and are a formula-level TypeError until a future phase's table
// update.
func toText(name string, v registry.Value) (string, error) {
	if v == nil {
		return "", nil
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", &apperror.FormulaError{
		Code:    apperror.ErrTypeMismatch,
		Message: fmt.Sprintf("%s: expected a string, got %T", name, v),
	}
}
