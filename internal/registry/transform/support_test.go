package transform

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/evaluator"
)

func TestToArray_RejectsNonArray(t *testing.T) {
	_, err := toArray("F", 5.0)
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toArray() error = %v, want ErrTypeMismatch", err)
	}
}

func TestToRecordArray_RejectsScalarArray(t *testing.T) {
	arr, _ := evaluator.NewArray([]any{1.0})
	if _, err := toRecordArray("F", arr); err == nil {
		t.Fatal("toRecordArray() error = nil, want type error for scalar array")
	}
}

func TestToRecordArray_AllowsEmptyArray(t *testing.T) {
	arr, _ := evaluator.NewArray([]any{})
	if _, err := toRecordArray("F", arr); err != nil {
		t.Errorf("toRecordArray() error = %v, want nil for empty array", err)
	}
}

func TestIsNullElement(t *testing.T) {
	if !isNullElement(nil) {
		t.Error("isNullElement(nil) = false, want true")
	}
	if isNullElement(5.0) {
		t.Error("isNullElement(5.0) = true, want false")
	}
}

func TestValuesEqual_MixedNumberAndNumericString(t *testing.T) {
	eq, err := valuesEqual(5.0, "5")
	if err != nil || !eq {
		t.Errorf("valuesEqual() = (%v, %v), want (true, nil)", eq, err)
	}
}
