package parser

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

func TestTokenize_UnterminatedString(t *testing.T) {
	_, err := Tokenize(`"abc`)
	if err == nil {
		t.Fatal("Tokenize() error = nil, want unterminated string error")
	}

	var fe *apperror.FormulaError
	if !errors.As(err, &fe) {
		t.Fatalf("error = %v, want *apperror.FormulaError", err)
	}
	if fe.Code != apperror.ErrSyntax {
		t.Errorf("Code = %v, want %v", fe.Code, apperror.ErrSyntax)
	}
	if fe.Line != 1 || fe.Column != 1 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 1)", fe.Line, fe.Column)
	}
}

func TestTokenize_UnrecognizedCharacter(t *testing.T) {
	_, err := Tokenize("1 @ 2")
	if err == nil {
		t.Fatal("Tokenize() error = nil, want unrecognized character error")
	}

	var fe *apperror.FormulaError
	if !errors.As(err, &fe) {
		t.Fatalf("error = %v, want *apperror.FormulaError", err)
	}
	if fe.Code != apperror.ErrSyntax {
		t.Errorf("Code = %v, want %v", fe.Code, apperror.ErrSyntax)
	}
	if fe.Line != 1 || fe.Column != 3 {
		t.Errorf("position = (line %d, col %d), want (line 1, col 3)", fe.Line, fe.Column)
	}
}
