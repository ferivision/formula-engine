package apperror

import "testing"

func TestFormulaError_ErrorWithoutPosition(t *testing.T) {
	err := &FormulaError{Code: ErrUndefinedReference, Message: "field \"total\" is not defined"}

	got := err.Error()
	want := `undefined_reference: field "total" is not defined`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFormulaError_ErrorWithPosition(t *testing.T) {
	err := &FormulaError{Code: ErrSyntax, Message: "unexpected token", Line: 1, Column: 7}

	got := err.Error()
	want := "syntax_error: unexpected token (line 1, column 7)"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFormulaError_CodeConstants(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{ErrSyntax, "syntax_error"},
		{ErrCircularReference, "circular_dependency_detected"},
		{ErrUndefinedReference, "undefined_reference"},
		{ErrTypeMismatch, "type_mismatch"},
		{ErrRuntime, "runtime_error"},
	}

	for _, tt := range tests {
		if string(tt.code) != tt.want {
			t.Errorf("code = %q, want %q", tt.code, tt.want)
		}
	}
}

func TestFormulaError_ImplementsError(t *testing.T) {
	var _ error = &FormulaError{}
}
