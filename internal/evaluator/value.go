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
// becomes a Record array, []any becomes either a Record or a scalar
// array depending on its elements (see newArrayFromAnySlice --
// []any of all map[string]any is what json.Unmarshal into `any`
// actually produces for a JSON array of objects, not
// []map[string]any directly, so this case matters in practice, not
// just in theory). Anything else -- including a []any that mixes
// Records and scalars -- is a formula-level TypeError.
func NewArray(v any) (Array, error) {
	switch s := v.(type) {
	case []map[string]any:
		elements := make([]Value, len(s))
		for i, m := range s {
			elements[i] = Record(m)
		}
		return Array{Elements: elements, IsRecord: true}, nil
	case []any:
		return newArrayFromAnySlice(s)
	default:
		return Array{}, newTypeError("expected an array ([]any or []map[string]any)")
	}
}

func newArrayFromAnySlice(s []any) (Array, error) {
	elements := make([]Value, len(s))
	hasRecord, hasScalar := false, false

	for i, e := range s {
		// A nil slot (rfc.md §9's "null element") is neutral -- it
		// belongs equally well in a Record array or a scalar array,
		// so it must not itself trigger the mixed-type rejection
		// below.
		if e == nil {
			continue
		}
		if m, ok := e.(map[string]any); ok {
			hasRecord = true
			elements[i] = Record(m)
			continue
		}
		hasScalar = true
		elements[i] = e
	}

	if hasRecord && hasScalar {
		return Array{}, newTypeError("array mixes records and scalar values")
	}

	return Array{Elements: elements, IsRecord: hasRecord}, nil
}
