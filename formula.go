package formulaengine

import "github.com/Ferivision/formula-engine/internal/apperror"

// FormulaInput is one named formula to evaluate.
type FormulaInput struct {
	Name       string // identifies this formula in the output map
	Expression string // e.g. "MAX(subtotal - discount, 0)"
}

// Result is the outcome of evaluating a single formula.
type Result struct {
	Value any   // the computed value, nil if Err is set
	Err   error // per-formula error, nil on success
}

// Evaluate computes one or more formulas against the given data and
// returns a result per formula name, or an error if the call as a
// whole cannot proceed (e.g. a circular reference across the inputs).
func Evaluate(formulas []FormulaInput, data map[string]any) (map[string]Result, error) {
	return nil, &apperror.FormulaError{Code: apperror.ErrRuntime, Message: "not implemented"}
}
