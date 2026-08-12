package evaluator

import "testing"

func TestContext_Lookup_ResolvesRecordSliceIntoArray(t *testing.T) {
	ctx := NewContext(map[string]any{
		"orders": []map[string]any{
			{"sku": "A1", "qty": 5.0},
			{"sku": "B2", "qty": 3.0},
		},
	})

	v, ok := ctx.Lookup("orders")
	if !ok {
		t.Fatal("Lookup() ok = false, want true")
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

	v, ok := ctx.Lookup("nums")
	if !ok {
		t.Fatal("Lookup() ok = false, want true")
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

	v, ok := ctx.Lookup("price")
	if !ok || v != 10.0 {
		t.Errorf("Lookup(price) = (%v, %v), want (10, true)", v, ok)
	}
	v, ok = ctx.Lookup("name")
	if !ok || v != "widget" {
		t.Errorf("Lookup(name) = (%v, %v), want (widget, true)", v, ok)
	}
}
