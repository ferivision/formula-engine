package evaluator

import (
	"errors"

	"github.com/ferivision/formula-engine/internal/registry"
)

// errFieldNotFound is Lookup's sentinel for "name isn't in the data
// map at all" -- distinct from a value being present but failing
// Array construction, so callers (evaluator.go) can tell FR-6's
// undefined-reference case apart from a formula-level TypeError.
var errFieldNotFound = errors.New("field not found")

// Context holds the field data a formula can reference by name during
// evaluation. A nil Context (or one built from a nil map) has no
// fields -- every identifier lookup misses.
type Context struct {
	data map[string]any
}

func NewContext(data map[string]any) *Context {
	return &Context{data: data}
}

// Lookup returns the value stored under name, or errFieldNotFound if
// name isn't present at all. A present-but-nil value (err == nil,
// v == nil) is distinct from an absent field (err == errFieldNotFound)
// -- see rfc.md §10 vs. FR-6.
//
// A []map[string]any or []any value is resolved into an Array
// (feature-0002 rfc.md §2); a value that's already an Array (e.g. a
// prior formula's computed result) passes through unchanged, since
// NewArray only understands raw Go slices. A malformed array (e.g. a
// mixed-type []any) returns NewArray's construction error directly,
// per rfc.md §9's "not deferred to first use" rule -- not swallowed,
// and not confusable with errFieldNotFound.
func (c *Context) Lookup(name string) (registry.Value, error) {
	if c == nil {
		return nil, errFieldNotFound
	}
	v, ok := c.data[name]
	if !ok {
		return nil, errFieldNotFound
	}

	if _, isArray := v.(Array); isArray {
		return v, nil
	}

	switch v.(type) {
	case []map[string]any, []any:
		return NewArray(v)
	}
	return v, nil
}
