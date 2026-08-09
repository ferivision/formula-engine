package math

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestCheckArgCount_WithinRange(t *testing.T) {
	if err := checkArgCount("F", 1, 3, []registry.Value{1, 2}); err != nil {
		t.Errorf("checkArgCount() error = %v, want nil", err)
	}
}

func TestCheckArgCount_Unbounded(t *testing.T) {
	if err := checkArgCount("F", 1, -1, []registry.Value{1, 2, 3, 4, 5}); err != nil {
		t.Errorf("checkArgCount() error = %v, want nil", err)
	}
}

func TestCheckArgCount_TooFew(t *testing.T) {
	err := checkArgCount("F", 2, -1, []registry.Value{1})
	assertRuntimeError(t, err)
}

func TestCheckArgCount_TooMany(t *testing.T) {
	err := checkArgCount("F", 1, 2, []registry.Value{1, 2, 3})
	assertRuntimeError(t, err)
}

func TestToFloat64_Float(t *testing.T) {
	got, err := toFloat64("F", 3.5)
	if err != nil || got != 3.5 {
		t.Errorf("toFloat64() = (%v, %v), want (3.5, nil)", got, err)
	}
}

func TestToFloat64_Int(t *testing.T) {
	got, err := toFloat64("F", 3)
	if err != nil || got != 3.0 {
		t.Errorf("toFloat64() = (%v, %v), want (3.0, nil)", got, err)
	}
}

func TestToFloat64_NonNumeric(t *testing.T) {
	_, err := toFloat64("F", "not a number")
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toFloat64() error = %v, want ErrTypeMismatch", err)
	}
}

func assertRuntimeError(t *testing.T, err error) {
	t.Helper()
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Fatalf("error = %v, want ErrRuntime", err)
	}
}
