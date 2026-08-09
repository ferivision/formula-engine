package evaluator

import (
	"strconv"

	"github.com/ferivision/formula-engine/internal/registry"
)

// coerceForArithmetic converts v into a float64 for use in +, -, *, /,
// per rfc.md §10: numbers pass through, bools become 1/0, a missing
// (nil) value becomes 0, numeric strings are parsed, and anything
// else is a formula-level TypeError.
func coerceForArithmetic(v registry.Value) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
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
			return 0, newTypeError("cannot parse \"" + n + "\" as a number")
		}
		return f, nil
	default:
		return 0, newTypeError("expected a number")
	}
}

// coerceForText applies rfc.md §10's text-context rule: a missing
// (nil) value becomes an empty string. Non-string, non-nil values
// have no defined text-context rule yet (see rfc.md §10) and are a
// formula-level TypeError until a future phase's table update.
func coerceForText(v registry.Value) (string, error) {
	if v == nil {
		return "", nil
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", newTypeError("expected a string")
}
