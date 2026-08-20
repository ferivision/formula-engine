package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

func TestNewArray_RecordSlice(t *testing.T) {
	input := []map[string]any{
		{"sku": "A1", "qty": 5.0},
		{"sku": "B2", "qty": 3.0},
	}
	arr, err := NewArray(input)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true")
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("len(Elements) = %d, want 2", len(arr.Elements))
	}
	rec, ok := arr.Elements[0].(Record)
	if !ok {
		t.Fatalf("Elements[0] = %T, want Record", arr.Elements[0])
	}
	if rec["sku"] != "A1" || rec["qty"] != 5.0 {
		t.Errorf("Elements[0] = %v, want {sku: A1, qty: 5}", rec)
	}
	rec1, _ := arr.Elements[1].(Record)
	if rec1["sku"] != "B2" {
		t.Errorf("Elements[1][\"sku\"] = %v, want B2", rec1["sku"])
	}
}

func TestNewArray_ScalarSlice(t *testing.T) {
	input := []any{1.0, 2.0, 3.0}
	arr, err := NewArray(input)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if arr.IsRecord {
		t.Error("IsRecord = true, want false")
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("len(Elements) = %d, want 3", len(arr.Elements))
	}
	if arr.Elements[0] != 1.0 || arr.Elements[1] != 2.0 || arr.Elements[2] != 3.0 {
		t.Errorf("Elements = %v, want [1 2 3]", arr.Elements)
	}
}

func TestNewArray_EmptyScalarSlice(t *testing.T) {
	arr, err := NewArray([]any{})
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if len(arr.Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(arr.Elements))
	}
}

func TestNewArray_EmptyRecordSlice(t *testing.T) {
	arr, err := NewArray([]map[string]any{})
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true for an empty []map[string]any")
	}
	if len(arr.Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(arr.Elements))
	}
}

// A []any where every element happens to be map[string]any is the
// natural shape JSON decoding produces (json.Unmarshal into `any`
// gives []any of map[string]any, never []map[string]any directly) --
// so this must be recognized as a Record array too, not just the
// concretely-typed []map[string]any case.
func TestNewArray_AllRecordsViaAnySlice(t *testing.T) {
	input := []any{
		map[string]any{"sku": "A1", "qty": 5.0},
		map[string]any{"sku": "B2", "qty": 3.0},
	}
	arr, err := NewArray(input)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true")
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("len(Elements) = %d, want 2", len(arr.Elements))
	}
	rec, ok := arr.Elements[0].(Record)
	if !ok {
		t.Fatalf("Elements[0] = %T, want Record", arr.Elements[0])
	}
	if rec["sku"] != "A1" {
		t.Errorf("Elements[0][\"sku\"] = %v, want A1", rec["sku"])
	}
}

// rfc.md §2: mixing scalars and Records is a formula-level TypeError
// at construction time.
func TestNewArray_MixedRecordAndScalarErrors(t *testing.T) {
	input := []any{
		map[string]any{"sku": "A1"},
		5.0,
	}
	_, err := NewArray(input)
	if err == nil {
		t.Fatal("NewArray() error = nil, want type error for mixed Record/scalar array")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrArrayTypeMismatch {
		t.Errorf("error = %v, want ErrArrayTypeMismatch", err)
	}
}

// rfc.md §9: "null element inside an Array passed to an aggregate ->
// Skipped." A nil slot must not itself count as "a scalar" for
// mixed-type detection -- otherwise an all-Records array with one
// blank slot would be wrongly rejected as mixing Records and
// scalars, when it's really just a Record array with a hole in it.
func TestNewArray_NilElementDoesNotTriggerMixedTypeRejection(t *testing.T) {
	input := []any{
		map[string]any{"sku": "A1"},
		nil,
		map[string]any{"sku": "B2"},
	}
	arr, err := NewArray(input)
	if err != nil {
		t.Fatalf("NewArray() error = %v, want nil element tolerated", err)
	}
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true")
	}
	if len(arr.Elements) != 3 || arr.Elements[1] != nil {
		t.Errorf("Elements = %v, want middle element nil", arr.Elements)
	}
}

func TestNewArray_NilElementAmongScalarsIsFine(t *testing.T) {
	input := []any{1.0, nil, 2.0}
	arr, err := NewArray(input)
	if err != nil {
		t.Fatalf("NewArray() error = %v", err)
	}
	if arr.IsRecord {
		t.Error("IsRecord = true, want false")
	}
	if len(arr.Elements) != 3 || arr.Elements[1] != nil {
		t.Errorf("Elements = %v, want middle element nil", arr.Elements)
	}
}

func TestNewArray_UnsupportedTypeErrors(t *testing.T) {
	_, err := NewArray("not an array")
	if err == nil {
		t.Fatal("NewArray() error = nil, want type error")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrTypeMismatch {
		t.Errorf("error = %v, want ErrTypeMismatch", err)
	}
}
