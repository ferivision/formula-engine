package math

import (
	"testing"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Audits every math function against rfc.md §10's coercion rows,
// the same rules arithmetic (internal/evaluator) already applies.
// A gap here means MAX(x, 1) and "x + 1" would disagree on what x
// coerces to, which is exactly what FASE-13.1 exists to catch.
func TestCoercionAudit_MathFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   registry.Function
		args []registry.Value
		want registry.Value
	}{
		{"MAX bool", maxFunction{}, []registry.Value{true, 0.0}, 1.0},
		{"MAX null", maxFunction{}, []registry.Value{nil, 3.0}, 3.0},
		{"MAX numeric string", maxFunction{}, []registry.Value{"5", 3.0}, 5.0},
		{"MAX int", maxFunction{}, []registry.Value{5, 3.0}, 5.0},

		{"MIN bool", minFunction{}, []registry.Value{false, 3.0}, 0.0},
		{"MIN null", minFunction{}, []registry.Value{nil, 3.0}, 0.0},
		{"MIN numeric string", minFunction{}, []registry.Value{"2", 3.0}, 2.0},
		{"MIN int", minFunction{}, []registry.Value{2, 3.0}, 2.0},

		{"SUM bool", sumFunction{}, []registry.Value{true, 1.0}, 2.0},
		{"SUM null", sumFunction{}, []registry.Value{nil, 1.0}, 1.0},
		{"SUM numeric string", sumFunction{}, []registry.Value{"2", 1.0}, 3.0},
		{"SUM int", sumFunction{}, []registry.Value{2, 1.0}, 3.0},

		{"AVG bool", avgFunction{}, []registry.Value{true, 1.0}, 1.0},
		{"AVG null", avgFunction{}, []registry.Value{nil, 4.0}, 2.0},
		{"AVG numeric string", avgFunction{}, []registry.Value{"3", 1.0}, 2.0},
		{"AVG int", avgFunction{}, []registry.Value{3, 1.0}, 2.0},

		{"ROUND bool digits", roundFunction{}, []registry.Value{3.14159, true}, 3.1},
		{"ROUND null digits", roundFunction{}, []registry.Value{3.7, nil}, 4.0},
		{"ROUND numeric string value", roundFunction{}, []registry.Value{"3.6"}, 4.0},
		{"ROUND int value", roundFunction{}, []registry.Value{3, 0.0}, 3.0},

		{"FLOOR bool", floorFunction{}, []registry.Value{true}, 1.0},
		{"FLOOR null", floorFunction{}, []registry.Value{nil}, 0.0},
		{"FLOOR numeric string", floorFunction{}, []registry.Value{"3.7"}, 3.0},
		{"FLOOR int", floorFunction{}, []registry.Value{3}, 3.0},

		{"CEIL bool", ceilFunction{}, []registry.Value{false}, 0.0},
		{"CEIL null", ceilFunction{}, []registry.Value{nil}, 0.0},
		{"CEIL numeric string", ceilFunction{}, []registry.Value{"3.2"}, 4.0},
		{"CEIL int", ceilFunction{}, []registry.Value{3}, 3.0},

		{"ABS bool", absFunction{}, []registry.Value{true}, 1.0},
		{"ABS null", absFunction{}, []registry.Value{nil}, 0.0},
		{"ABS numeric string", absFunction{}, []registry.Value{"-5"}, 5.0},
		{"ABS int", absFunction{}, []registry.Value{-5}, 5.0},
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

func TestCoercionAudit_MathFunctions_NonNumericStringErrors(t *testing.T) {
	fns := []registry.Function{
		maxFunction{}, minFunction{}, sumFunction{}, avgFunction{},
		roundFunction{}, floorFunction{}, ceilFunction{}, absFunction{},
	}
	for _, fn := range fns {
		t.Run(fn.Name(), func(t *testing.T) {
			args := make([]registry.Value, fn.MinArgs())
			for i := range args {
				args[i] = "not a number"
			}
			if _, err := fn.Evaluate(args); err == nil {
				t.Errorf("Evaluate() error = nil, want type error for non-numeric string")
			}
		})
	}
}
