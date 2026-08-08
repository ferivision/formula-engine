package registry

// Value is a computed value flowing into and out of a function call:
// a float64, string, bool, or nil (missing/null), per rfc.md §10.
type Value = any

// Function is the contract every built-in function implements. Each
// built-in lives in its own file and registers itself via Register in
// an init() function -- see rfc.md §6.
type Function interface {
	Name() string
	MinArgs() int
	MaxArgs() int // -1 = unbounded
	ValidateArgTypes(args []Value) error
	Evaluate(args []Value) (Value, error)
}
