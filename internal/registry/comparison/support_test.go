package comparison

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestToFloat64_Number(t *testing.T) {
	got, err := toFloat64(5.0)
	if err != nil || got != 5.0 {
		t.Errorf("toFloat64() = (%v, %v), want (5, nil)", got, err)
	}
}

func TestToFloat64_Bool(t *testing.T) {
	got, err := toFloat64(true)
	if err != nil || got != 1.0 {
		t.Errorf("toFloat64() = (%v, %v), want (1, nil)", got, err)
	}
}

func TestToFloat64_Nil(t *testing.T) {
	got, err := toFloat64(nil)
	if err != nil || got != 0.0 {
		t.Errorf("toFloat64() = (%v, %v), want (0, nil)", got, err)
	}
}

func TestToFloat64_NumericString(t *testing.T) {
	got, err := toFloat64("5")
	if err != nil || got != 5.0 {
		t.Errorf("toFloat64() = (%v, %v), want (5, nil)", got, err)
	}
}

func TestToFloat64_NonNumericString(t *testing.T) {
	_, err := toFloat64("abc")
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toFloat64() error = %v, want ErrTypeMismatch", err)
	}
}

func TestCheckArgCount_TooFew(t *testing.T) {
	err := checkArgCount("F", 2, 2, []registry.Value{1.0})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("checkArgCount() error = %v, want ErrRuntime", err)
	}
}
