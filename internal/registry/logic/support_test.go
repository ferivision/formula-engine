package logic

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestToBool_True(t *testing.T) {
	got, err := toBool("F", true)
	if err != nil || got != true {
		t.Errorf("toBool() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestToBool_False(t *testing.T) {
	got, err := toBool("F", false)
	if err != nil || got != false {
		t.Errorf("toBool() = (%v, %v), want (false, nil)", got, err)
	}
}

func TestToBool_NonBool(t *testing.T) {
	_, err := toBool("F", 1.0)
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toBool() error = %v, want ErrTypeMismatch", err)
	}
}

func TestCheckArgCount_TooFew(t *testing.T) {
	err := checkArgCount("F", 2, -1, []registry.Value{true})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("checkArgCount() error = %v, want ErrRuntime", err)
	}
}
