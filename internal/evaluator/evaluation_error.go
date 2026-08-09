package evaluator

import "github.com/ferivision/formula-engine/internal/apperror"

func newRuntimeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrRuntime, Message: message}
}

func newTypeError(message string) error {
	return &apperror.FormulaError{Code: apperror.ErrTypeMismatch, Message: message}
}

func newUndefinedReferenceError(name string) error {
	return &apperror.FormulaError{
		Code:    apperror.ErrUndefinedReference,
		Message: "field \"" + name + "\" is not defined",
	}
}
