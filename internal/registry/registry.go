package registry

import "github.com/ferivision/formula-engine/internal/apperror"

var functions = make(map[string]Function)

// Register adds fn to the registry under fn.Name(), overwriting any
// existing registration under that name. Intended to be called from
// each built-in's init() function, before any Evaluate call runs.
func Register(fn Function) {
	functions[fn.Name()] = fn
}

// Lookup finds a registered function by name.
func Lookup(name string) (Function, error) {
	fn, ok := functions[name]
	if !ok {
		return nil, &apperror.FormulaError{
			Code:    apperror.ErrUndefinedReference,
			Message: "unknown function: " + name,
		}
	}
	return fn, nil
}
