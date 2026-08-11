package date

import (
	"testing"
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Audits DATE_ADD's amount argument against rfc.md §10's coercion
// rows, the same rules arithmetic (internal/evaluator) already
// applies to any other numeric-context argument.
func TestCoercionAudit_DateAdd_AmountCoercion(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		amount registry.Value
		want   time.Time
	}{
		{"bool amount", true, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"null amount", nil, start},
		{"numeric string amount", "3", time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)},
		{"int amount", 5, time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := (dateAddFunction{}).Evaluate([]registry.Value{start, tt.amount, "days"})
			if err != nil {
				t.Fatalf("Evaluate() error = %v, want %v", err, tt.want)
			}
			if !got.(time.Time).Equal(tt.want) {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
