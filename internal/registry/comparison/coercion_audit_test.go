package comparison

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Audits EQUALS and BETWEEN against rfc.md §10's coercion rows, the
// same rules arithmetic (internal/evaluator) already applies.
func TestCoercionAudit_ComparisonFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   registry.Function
		args []registry.Value
		want registry.Value
	}{
		{"EQUALS bool vs number", equalsFunction{}, []registry.Value{true, 1.0}, true},
		{"EQUALS null vs zero", equalsFunction{}, []registry.Value{nil, 0.0}, true},
		{"EQUALS numeric string vs number", equalsFunction{}, []registry.Value{"5", 5.0}, true},
		{"EQUALS int vs float", equalsFunction{}, []registry.Value{5, 5.0}, true},

		{"BETWEEN bool value", betweenFunction{}, []registry.Value{true, 0.0, 1.0}, true},
		{"BETWEEN null value", betweenFunction{}, []registry.Value{nil, 0.0, 1.0}, true},
		{"BETWEEN numeric string bound", betweenFunction{}, []registry.Value{5.0, "1", "10"}, true},
		{"BETWEEN int value", betweenFunction{}, []registry.Value{5, 1.0, 10.0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn.Evaluate(tt.args)
			if err != nil {
				t.Fatalf("Evaluate(%v) error = %v, want %v", tt.args, err, tt.want)
			}
			if got != tt.want {
				t.Errorf("Evaluate(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
