package lookup

import (
	"fmt"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(indexFunction{})
}

type indexFunction struct{}

func (indexFunction) Name() string { return "INDEX" }
func (indexFunction) MinArgs() int { return 2 }
func (indexFunction) MaxArgs() int { return 2 }

func (f indexFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toArray(f.Name(), args[0]); err != nil {
		return err
	}
	_, err := toFloat64(args[1])
	return err
}

// Evaluate returns the element at the given 1-based position, from
// either a scalar or a Record array. An out-of-range position is a
// plain formula-level ErrRuntime -- rfc.md §4 explicitly says to
// reuse this existing code rather than add a new one, since there's
// no call-level implication.
func (f indexFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	arr, _ := toArray(f.Name(), args[0])
	pos, _ := toFloat64(args[1])

	i := int(pos) - 1
	if i < 0 || i >= len(arr.Elements) {
		return nil, newRuntimeError(fmt.Sprintf("%s: position %v out of range (array has %d elements)", f.Name(), pos, len(arr.Elements)))
	}
	return arr.Elements[i], nil
}
