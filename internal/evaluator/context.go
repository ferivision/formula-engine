package evaluator

import "github.com/ferivision/formula-engine/internal/registry"

// Context holds the field data a formula can reference by name during
// evaluation. A nil Context (or one built from a nil map) has no
// fields -- every identifier lookup misses.
type Context struct {
	data map[string]any
}

func NewContext(data map[string]any) *Context {
	return &Context{data: data}
}

// Lookup returns the value stored under name and whether it was
// present at all. A present-but-nil value (ok == true, v == nil) is
// distinct from an absent field (ok == false) -- see rfc.md §10 vs.
// FR-6.
//
// A []map[string]any or []any value is resolved into an Array
// (feature-0002 rfc.md §2) rather than passed through raw. A
// conversion failure (e.g. a mixed-type []any, see NewArray) falls
// through to the raw value for now rather than being surfaced here --
// propagating that error properly is F2-FASE-2.2's job.
func (c *Context) Lookup(name string) (registry.Value, bool) {
	if c == nil {
		return nil, false
	}
	v, ok := c.data[name]
	if !ok {
		return nil, false
	}

	switch v.(type) {
	case []map[string]any, []any:
		if arr, err := NewArray(v); err == nil {
			return arr, true
		}
	}
	return v, true
}
