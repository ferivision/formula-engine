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
