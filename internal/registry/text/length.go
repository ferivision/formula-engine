package text

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(lengthFunction{})
}

type lengthFunction struct{}

func (lengthFunction) Name() string { return "LENGTH" }
func (lengthFunction) MinArgs() int { return 1 }
func (lengthFunction) MaxArgs() int { return 1 }

func (f lengthFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toText(f.Name(), args[0])
	return err
}

// Evaluate counts runes, not bytes, so multi-byte characters (e.g.
// "é") count as one character each.
func (f lengthFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	s, _ := toText(f.Name(), args[0])
	return float64(len([]rune(s))), nil
}
