package apperror

import "fmt"

type ErrorCode string

const (
	ErrSyntax             ErrorCode = "syntax_error"
	ErrCircularReference  ErrorCode = "circular_dependency_detected"
	ErrUndefinedReference ErrorCode = "undefined_reference"
	ErrTypeMismatch       ErrorCode = "type_mismatch"
	ErrRuntime            ErrorCode = "runtime_error"
)

type FormulaError struct {
	Code    ErrorCode
	Message string
	Line    int
	Column  int
}

func (e *FormulaError) Error() string {
	if e.Line != 0 || e.Column != 0 {
		return fmt.Sprintf("%s: %s (line %d, column %d)", e.Code, e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
