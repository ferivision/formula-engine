package evaluator

import "github.com/ferivision/formula-engine/internal/registry"

// Value mirrors registry.Value, so Array/Record and evaluator code
// can use the short name the way rfc.md-0002's snippets do.
type Value = registry.Value

// Array is an ordered collection: either of scalar Values, or of
// Records (see below). A single Array is one or the other, not
// mixed -- mixing is a formula-level TypeError at construction time
// (F2-FASE-1.2), not something functions need to guard against
// individually.
type Array struct {
	Elements []Value
	IsRecord bool
}

// Record is a flexible key-value map, matching how top-level `data`
// already works -- map[string]any, not a fixed schema.
type Record map[string]Value

// NewArray converts a Go slice into an Array: []map[string]any
// becomes a Record array, []any becomes a scalar array. Anything
// else is a formula-level TypeError.
func NewArray(v any) (Array, error) {
	switch s := v.(type) {
	case []map[string]any:
		elements := make([]Value, len(s))
		for i, m := range s {
			elements[i] = Record(m)
		}
		return Array{Elements: elements, IsRecord: true}, nil
	case []any:
		elements := make([]Value, len(s))
		copy(elements, s)
		return Array{Elements: elements, IsRecord: false}, nil
	default:
		return Array{}, newTypeError("expected an array ([]any or []map[string]any)")
	}
}
