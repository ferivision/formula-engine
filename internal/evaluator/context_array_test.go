package evaluator

import (
	"errors"
	"testing"

	"github.com/ferivision/formula-engine/internal/apperror"
)

func TestContext_Lookup_ResolvesRecordSliceIntoArray(t *testing.T) {
	ctx := NewContext(map[string]any{
		"orders": []map[string]any{
			{"sku": "A1", "qty": 5.0},
			{"sku": "B2", "qty": 3.0},
		},
	})

	v, err := ctx.Lookup("orders")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	arr, isArray := v.(Array)
	if !isArray {
		t.Fatalf("Lookup() = %T, want Array", v)
	}
	if !arr.IsRecord {
		t.Error("IsRecord = false, want true")
	}
	if len(arr.Elements) != 2 {
		t.Errorf("len(Elements) = %d, want 2", len(arr.Elements))
	}
}

func TestContext_Lookup_ResolvesScalarSliceIntoArray(t *testing.T) {
	ctx := NewContext(map[string]any{
		"nums": []any{1.0, 2.0, 3.0},
	})

	v, err := ctx.Lookup("nums")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	arr, isArray := v.(Array)
	if !isArray {
		t.Fatalf("Lookup() = %T, want Array", v)
	}
	if arr.IsRecord {
		t.Error("IsRecord = true, want false")
	}
	if len(arr.Elements) != 3 {
		t.Errorf("len(Elements) = %d, want 3", len(arr.Elements))
	}
}

func TestContext_Lookup_ScalarFieldUnaffected(t *testing.T) {
	ctx := NewContext(map[string]any{"price": 10.0, "name": "widget"})

	v, err := ctx.Lookup("price")
	if err != nil || v != 10.0 {
		t.Errorf("Lookup(price) = (%v, %v), want (10, nil)", v, err)
	}
	v, err = ctx.Lookup("name")
	if err != nil || v != "widget" {
		t.Errorf("Lookup(name) = (%v, %v), want (widget, nil)", v, err)
	}
}

// rfc.md §9: a mixed-type Array is a TypeError "at the point the
// Array is constructed... not deferred to first use" -- so a
// malformed field must fail as soon as it's looked up, not silently
// produce a wrong/empty Array.
func TestContext_Lookup_SurfacesArrayConstructionError(t *testing.T) {
	ctx := NewContext(map[string]any{
		"orders": []any{map[string]any{"sku": "A1"}, 5.0},
	})

	_, err := ctx.Lookup("orders")
	if err == nil {
		t.Fatal("Lookup() error = nil, want a construction error for mixed-type array")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrArrayTypeMismatch {
		t.Errorf("error = %v, want ErrArrayTypeMismatch", err)
	}
}

// A field absent entirely is a different situation from a malformed
// array -- must still be distinguishable so evaluator.go can report
// ErrUndefinedReference only for the genuinely-absent case.
func TestContext_Lookup_AbsentFieldIsDistinctFromConstructionError(t *testing.T) {
	ctx := NewContext(map[string]any{"price": 10.0})

	_, err := ctx.Lookup("missing")
	if err == nil {
		t.Fatal("Lookup() error = nil, want a not-found error")
	}
	if !errors.Is(err, errFieldNotFound) {
		t.Errorf("error = %v, want errFieldNotFound", err)
	}
}

// An identifier referencing a previously-computed formula result that
// is itself an Array (e.g. from FILTER, once Phase 4 exists) must
// resolve to that Array unchanged -- NewArray only understands raw Go
// slices, so re-running it on an already-converted Array would be
// wrong (and can't even match its type switch).
func TestContext_Lookup_AlreadyComputedArrayPassesThroughUnchanged(t *testing.T) {
	precomputed := Array{IsRecord: false, Elements: []Value{1.0, 2.0}}
	ctx := NewContext(map[string]any{"filtered": precomputed})

	v, err := ctx.Lookup("filtered")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	arr, ok := v.(Array)
	if !ok {
		t.Fatalf("Lookup() = %T, want Array", v)
	}
	if len(arr.Elements) != 2 || arr.Elements[0] != 1.0 {
		t.Errorf("Elements = %v, want unchanged [1 2]", arr.Elements)
	}
}

// Full-pipeline regression checks through Evaluate, not just Lookup
// directly -- confirms evaluator.go's NodeIdentifier case tells a
// malformed array apart from a genuinely undefined field.
func TestEvaluate_MixedArrayFieldIsTypeErrorNotUndefinedReference(t *testing.T) {
	node := parseHelper(t, "orders")
	_, err := Evaluate(node, NewContext(map[string]any{
		"orders": []any{map[string]any{"sku": "A1"}, 5.0},
	}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want a type-mismatch error")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrArrayTypeMismatch {
		t.Errorf("error = %v, want ErrArrayTypeMismatch (not ErrUndefinedReference)", err)
	}
}

func TestEvaluate_GenuinelyUndefinedFieldStillReportsUndefinedReference(t *testing.T) {
	node := parseHelper(t, "missing")
	_, err := Evaluate(node, NewContext(map[string]any{"price": 10.0}))
	if err == nil {
		t.Fatal("Evaluate() error = nil, want undefined-reference error")
	}
	var fe *apperror.FormulaError
	if !errors.As(err, &fe) || fe.Code != apperror.ErrUndefinedReference {
		t.Errorf("error = %v, want ErrUndefinedReference", err)
	}
}
