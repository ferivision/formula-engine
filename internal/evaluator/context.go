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
func (c *Context) Lookup(name string) (registry.Value, bool) {
	if c == nil {
		return nil, false
	}
	v, ok := c.data[name]
	return v, ok
}
