package text

import (
	"strings"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(concatFunction{})
}

type concatFunction struct{}

func (concatFunction) Name() string { return "CONCAT" }
func (concatFunction) MinArgs() int { return 1 }
func (concatFunction) MaxArgs() int { return -1 }

func (f concatFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	for _, a := range args {
		if _, err := toText(f.Name(), a); err != nil {
			return err
		}
	}
	return nil
}

func (f concatFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	var sb strings.Builder
	for _, a := range args {
		s, _ := toText(f.Name(), a)
		sb.WriteString(s)
	}
	return sb.String(), nil
}
