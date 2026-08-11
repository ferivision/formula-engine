package text

import (
	"strings"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(trimFunction{})
}

type trimFunction struct{}

func (trimFunction) Name() string { return "TRIM" }
func (trimFunction) MinArgs() int { return 1 }
func (trimFunction) MaxArgs() int { return 1 }

func (f trimFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	_, err := toText(f.Name(), args[0])
	return err
}

func (f trimFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	s, _ := toText(f.Name(), args[0])
	return strings.TrimSpace(s), nil
}
