package formulaengine

import (
	"fmt"
	"sync"
	"testing"
)

// NFR-3: Evaluate must be safe to call concurrently from multiple
// goroutines, with no interference between calls. Each goroutine here
// uses distinct inputs so a shared-state bug would show up as a
// wrong result, not just a crash; run with `make test-race` to also
// catch data races the assertions alone wouldn't reveal.
func TestEvaluate_ConcurrentCallsAreIndependent(t *testing.T) {
	const goroutines = 200

	var wg sync.WaitGroup
	errs := make(chan string, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			price := float64(i + 1)
			quantity := 2.0
			results, err := Evaluate([]FormulaInput{
				{Name: "subtotal", Expression: "price * quantity"},
				{Name: "total", Expression: "MAX(subtotal - 1, 0)"},
			}, map[string]any{"price": price, "quantity": quantity})
			if err != nil {
				errs <- fmt.Sprintf("goroutine %d: Evaluate() error = %v", i, err)
				return
			}

			wantSubtotal := price * quantity
			if results["subtotal"].Value != wantSubtotal {
				errs <- fmt.Sprintf("goroutine %d: subtotal = %v, want %v", i, results["subtotal"].Value, wantSubtotal)
			}

			wantTotal := wantSubtotal - 1
			if wantTotal < 0 {
				wantTotal = 0
			}
			if results["total"].Value != wantTotal {
				errs <- fmt.Sprintf("goroutine %d: total = %v, want %v", i, results["total"].Value, wantTotal)
			}
		}(i)
	}

	wg.Wait()
	close(errs)
	for e := range errs {
		t.Error(e)
	}
}
