package parser

import (
	"errors"
	"testing"

	"github.com/Ferivision/formula-engine/internal/apperror"
)

func parseErr(t *testing.T, input string) *apperror.FormulaError {
	t.Helper()
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", input, err)
	}
	_, err = Parse(tokens)
	if err == nil {
		t.Fatalf("Parse(%q) error = nil, want a syntax error", input)
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) {
		t.Fatalf("Parse(%q) error = %v, want *apperror.FormulaError", input, err)
	}
	if fe.Code != apperror.ErrSyntax {
		t.Fatalf("Parse(%q) code = %v, want %v", input, fe.Code, apperror.ErrSyntax)
	}
	return fe
}

func TestParse_UnbalancedParenMissingClose(t *testing.T) {
	fe := parseErr(t, "(1 + 2")
	if fe.Line != 1 || fe.Column != 7 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 7) at end of input", fe.Line, fe.Column)
	}
}

func TestParse_UnbalancedParenExtraClose(t *testing.T) {
	fe := parseErr(t, "1 + 2)")
	if fe.Line != 1 || fe.Column != 6 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 6) at the stray ')'", fe.Line, fe.Column)
	}
}

func TestParse_TrailingOperator(t *testing.T) {
	fe := parseErr(t, "1 +")
	if fe.Line != 1 || fe.Column != 4 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 4) at end of input", fe.Line, fe.Column)
	}
}

func TestParse_FunctionCallTrailingComma(t *testing.T) {
	fe := parseErr(t, "MAX(1,)")
	if fe.Line != 1 || fe.Column != 7 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 7) at the stray ')'", fe.Line, fe.Column)
	}
}

func TestParse_UnclosedFunctionCall(t *testing.T) {
	fe := parseErr(t, "MAX(1, 2")
	if fe.Line != 1 || fe.Column != 9 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 9) at end of input", fe.Line, fe.Column)
	}
}

func TestParse_EmptyInput(t *testing.T) {
	fe := parseErr(t, "")
	if fe.Line != 1 || fe.Column != 1 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 1)", fe.Line, fe.Column)
	}
}

func TestParse_UnexpectedLeadingOperator(t *testing.T) {
	fe := parseErr(t, "* 2")
	if fe.Line != 1 || fe.Column != 1 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 1)", fe.Line, fe.Column)
	}
}
