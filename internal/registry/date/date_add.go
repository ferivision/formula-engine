package date

import "github.com/ferivision/formula-engine/internal/registry"

func init() {
	registry.Register(dateAddFunction{})
}

type dateAddFunction struct{}

func (dateAddFunction) Name() string { return "DATE_ADD" }
func (dateAddFunction) MinArgs() int { return 3 }
func (dateAddFunction) MaxArgs() int { return 3 }

func (f dateAddFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toTime(f.Name(), args[0]); err != nil {
		return err
	}
	if _, err := toFloat64(f.Name(), args[1]); err != nil {
		return err
	}
	_, err := toUnit(f.Name(), args[2])
	return err
}

// Evaluate adds amount of unit ("days", "months", or "years") to
// date, via time.Time.AddDate. When a month/year doesn't have the
// original day (e.g. adding a month to Jan 31), Go normalizes into
// the following month rather than clamping -- that's AddDate's
// documented behavior, not a bug here.
func (f dateAddFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	t, _ := toTime(f.Name(), args[0])
	amount, _ := toFloat64(f.Name(), args[1])
	unit, _ := toUnit(f.Name(), args[2])

	n := int(amount)
	switch unit {
	case "days":
		return t.AddDate(0, 0, n), nil
	case "months":
		return t.AddDate(0, n, 0), nil
	case "years":
		return t.AddDate(n, 0, 0), nil
	}
	return nil, newRuntimeError(f.Name() + ": unreachable unit " + unit)
}

func toUnit(name string, v registry.Value) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", newTypeError(name + ": expected a unit string (\"days\", \"months\", or \"years\")")
	}
	switch s {
	case "days", "months", "years":
		return s, nil
	default:
		return "", newRuntimeError(name + ": unsupported unit " + s + " (want \"days\", \"months\", or \"years\")")
	}
}
