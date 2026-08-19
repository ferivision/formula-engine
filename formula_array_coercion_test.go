package formulaengine

import "testing"

// rfc.md (0002) §9: accessing a non-existent key on a Record resolves
// to null, not an error, consistent with feature-0001's existing
// null-handling per context (arithmetic -> 0, text -> passthrough as
// nil for a raw lookup result).
func TestCoercion_MissingRecordFieldTreatedAsNullInSumif(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "total", Expression: `SUMIF(orders, "status", "shipped", "nonexistent_field")`},
	}, map[string]any{
		"orders": []map[string]any{{"status": "shipped", "qty": 5.0}},
	})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["total"].Err != nil {
		t.Fatalf("total.Err = %v, want nil (missing field -> null -> 0, not an error)", results["total"].Err)
	}
	if results["total"].Value != 0.0 {
		t.Errorf("total.Value = %v, want 0", results["total"].Value)
	}
}

func TestCoercion_MissingRecordFieldReturnedAsNilFromVlookup(t *testing.T) {
	results, err := Evaluate([]FormulaInput{
		{Name: "result", Expression: `VLOOKUP("A1", table, "sku", "nonexistent_field")`},
	}, map[string]any{
		"table": []map[string]any{{"sku": "A1"}},
	})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if results["result"].Err != nil {
		t.Fatalf("result.Err = %v, want nil (missing field -> null, not an error)", results["result"].Err)
	}
	if results["result"].Value != nil {
		t.Errorf("result.Value = %v, want nil", results["result"].Value)
	}
}
