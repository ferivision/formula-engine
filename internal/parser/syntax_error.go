package parser

import "github.com/ferivision/formula-engine/internal/apperror"

func newSyntaxError(message string, line, column int) error {
	return &apperror.FormulaError{
		Code:    apperror.ErrSyntax,
		Message: message,
		Line:    line,
		Column:  column,
	}
}
