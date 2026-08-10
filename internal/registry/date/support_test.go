package date

import (
	"errors"
	"testing"
	"time"

	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/registry"
)

func TestToTime_RealTime(t *testing.T) {
	want := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := toTime("F", want)
	if err != nil || !got.Equal(want) {
		t.Errorf("toTime() = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestToTime_NonTimeErrors(t *testing.T) {
	_, err := toTime("F", "2024-01-01")
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("toTime() error = %v, want ErrTypeMismatch", err)
	}
}

func TestCheckArgCount_TooMany(t *testing.T) {
	err := checkArgCount("F", 0, 0, []registry.Value{1.0})
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrRuntime {
		t.Errorf("checkArgCount() error = %v, want ErrRuntime", err)
	}
}
