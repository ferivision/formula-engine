package aggregate

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestToRecordArray_RejectsScalarArray(t *testing.T) {
	scalarArr, err := evaluator.NewArray([]any{1.0, 2.0})
	if err != nil {
		t.Fatalf("setup: NewArray() error = %v", err)
	}
	_, err = toRecordArray("F", scalarArr)
	if err == nil {
		t.Fatal("toRecordArray() error = nil, want type error for scalar array")
	}
}

func TestToRecordArray_AllowsEmptyArrayRegardlessOfIsRecord(t *testing.T) {
	emptyArr, err := evaluator.NewArray([]any{})
	if err != nil {
		t.Fatalf("setup: NewArray() error = %v", err)
	}
	if _, err := toRecordArray("F", emptyArr); err != nil {
		t.Errorf("toRecordArray() error = %v, want nil for empty array", err)
	}
}

func TestIsNullElement(t *testing.T) {
	if !isNullElement(nil) {
		t.Error("isNullElement(nil) = false, want true")
	}
	if !isNullElement(evaluator.Record(nil)) {
		t.Error("isNullElement(Record(nil)) = false, want true")
	}
	if isNullElement(evaluator.Record{"a": 1.0}) {
		t.Error("isNullElement(non-nil Record) = true, want false")
	}
	if isNullElement(5.0) {
		t.Error("isNullElement(5.0) = true, want false")
	}
}

func TestFieldValue_MissingKeyIsNil(t *testing.T) {
	rec := evaluator.Record{"a": 1.0}
	if fieldValue(rec, "missing") != nil {
		t.Errorf("fieldValue() = %v, want nil for missing key", fieldValue(rec, "missing"))
	}
	if fieldValue(rec, "a") != 1.0 {
		t.Errorf("fieldValue() = %v, want 1", fieldValue(rec, "a"))
	}
}

func TestValuesEqual_MixedNumberAndNumericString(t *testing.T) {
	eq, err := valuesEqual(5.0, "5")
	if err != nil || !eq {
		t.Errorf("valuesEqual() = (%v, %v), want (true, nil)", eq, err)
	}
}

func TestValuesEqual_NonNumericStringErrors(t *testing.T) {
	_, err := valuesEqual(5.0, "abc")
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("valuesEqual() error = %v, want ErrTypeMismatch", err)
	}
}

func TestCheckArgCount_TooFew(t *testing.T) {
	err := checkArgCount("F", 3, 3, []registry.Value{1.0})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("checkArgCount() error = %v, want ErrRuntime", err)
	}
}
