package date

import (
	"fmt"
	"time"

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

func toTime(name string, v registry.Value) (time.Time, error) {
	t, ok := v.(time.Time)
	if !ok {
		return time.Time{}, newTypeError(fmt.Sprintf("%s: expected a date, got %T", name, v))
	}
	return t, nil
}

func toFloat64(name string, v registry.Value) (float64, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, newTypeError(fmt.Sprintf("%s: expected a number, got %T", name, v))
	}
	return f, nil
}
