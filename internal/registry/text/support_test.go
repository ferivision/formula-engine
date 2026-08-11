package text

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestToText_StringPassesThrough(t *testing.T) {
	got, err := toText("F", "hello")
	if err != nil || got != "hello" {
		t.Errorf("toText() = (%q, %v), want (\"hello\", nil)", got, err)
	}
}

func TestToText_NullBecomesEmptyString(t *testing.T) {
	got, err := toText("F", nil)
	if err != nil || got != "" {
		t.Errorf("toText() = (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestToText_NonStringErrors(t *testing.T) {
	_, err := toText("F", 1.0)
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toText() error = %v, want ErrTypeMismatch", err)
	}
}

func TestCheckArgCount_TooFew(t *testing.T) {
	err := checkArgCount("F", 2, -1, []registry.Value{"x"})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("checkArgCount() error = %v, want ErrRuntime", err)
	}
}
