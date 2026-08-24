package formulaengine

import (
	"fmt"
	"testing"
)

// buildChain returns n formulas, each depending on the previous one
// (f0 depends on the "base" data field; f1 depends on f0; and so
// on), in reverse order -- so the benchmark also exercises
// TopologicalSort actually doing work, not just confirming an
// already-sorted input.
func buildChain(n int) []FormulaInput {
	formulas := make([]FormulaInput, n)
	formulas[0] = FormulaInput{Name: "f0", Expression: "base + 1"}
	for i := 1; i < n; i++ {
		formulas[i] = FormulaInput{
			Name:       fmt.Sprintf("f%d", i),
			Expression: fmt.Sprintf("f%d + 1", i-1),
		}
	}
	// Reverse so formulas are listed dependent-before-dependency.
	for i, j := 0, len(formulas)-1; i < j; i, j = i+1, j-1 {
		formulas[i], formulas[j] = formulas[j], formulas[i]
	}
	return formulas
}

func benchmarkChain(b *testing.B, n int) {
	formulas := buildChain(n)
	data := map[string]any{"base": 1.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Evaluate(formulas, data); err != nil {
			b.Fatalf("Evaluate() error = %v", err)
		}
	}
}

func BenchmarkEvaluate_Chain10(b *testing.B)    { benchmarkChain(b, 10) }
func BenchmarkEvaluate_Chain100(b *testing.B)   { benchmarkChain(b, 100) }
func BenchmarkEvaluate_Chain1000(b *testing.B)  { benchmarkChain(b, 1000) }
func BenchmarkEvaluate_Chain10000(b *testing.B) { benchmarkChain(b, 10000) }

// buildOrders returns n order records with an alternating status, so
// a status-filtering benchmark matches roughly half of them.
func buildOrders(n int) []map[string]any {
	orders := make([]map[string]any, n)
	for i := 0; i < n; i++ {
		status := "pending"
		if i%2 == 0 {
			status = "shipped"
		}
		orders[i] = map[string]any{
			"sku":    fmt.Sprintf("SKU-%d", i),
			"status": status,
			"qty":    float64(n - i),
		}
	}
	return orders
}

func benchmarkArrayCall(b *testing.B, n int, expression string) {
	formulas := []FormulaInput{{Name: "result", Expression: expression}}
	data := map[string]any{"orders": buildOrders(n)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results, err := Evaluate(formulas, data)
		if err != nil {
			b.Fatalf("Evaluate() error = %v", err)
		}
		if results["result"].Err != nil {
			b.Fatalf("result.Err = %v", results["result"].Err)
		}
	}
}

// SUMIF: one O(n) scan over the Array, representative of the
// conditional-aggregate category (rfc.md §14 Phase 3).
func benchmarkSumif(b *testing.B, n int) {
	benchmarkArrayCall(b, n, `SUMIF(orders, "status", "shipped", "qty")`)
}

func BenchmarkSumif_Array10(b *testing.B)    { benchmarkSumif(b, 10) }
func BenchmarkSumif_Array100(b *testing.B)   { benchmarkSumif(b, 100) }
func BenchmarkSumif_Array1000(b *testing.B)  { benchmarkSumif(b, 1000) }
func BenchmarkSumif_Array10000(b *testing.B) { benchmarkSumif(b, 10000) }

// SORT: O(n log n) via sort.SliceStable, representative of the
// transformation category (rfc.md §14 Phase 4) -- the one function
// in this feature with worse-than-linear complexity.
func benchmarkSort(b *testing.B, n int) {
	benchmarkArrayCall(b, n, `SORT(orders, "qty", "asc")`)
}

func BenchmarkSort_Array10(b *testing.B)    { benchmarkSort(b, 10) }
func BenchmarkSort_Array100(b *testing.B)   { benchmarkSort(b, 100) }
func BenchmarkSort_Array1000(b *testing.B)  { benchmarkSort(b, 1000) }
func BenchmarkSort_Array10000(b *testing.B) { benchmarkSort(b, 10000) }

// VLOOKUP: worst-case O(n) scan (the key is the last element, forcing
// a full pass), representative of the lookup category (rfc.md §14
// Phase 5).
func benchmarkVlookup(b *testing.B, n int) {
	benchmarkArrayCall(b, n, fmt.Sprintf(`VLOOKUP("SKU-%d", orders, "sku", "qty")`, n-1))
}

func BenchmarkVlookup_Array10(b *testing.B)    { benchmarkVlookup(b, 10) }
func BenchmarkVlookup_Array100(b *testing.B)   { benchmarkVlookup(b, 100) }
func BenchmarkVlookup_Array1000(b *testing.B)  { benchmarkVlookup(b, 1000) }
func BenchmarkVlookup_Array10000(b *testing.B) { benchmarkVlookup(b, 10000) }
