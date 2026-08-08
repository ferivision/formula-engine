package registry

import "testing"

type fakeFunction struct{}

func (fakeFunction) Name() string                         { return "FAKE" }
func (fakeFunction) MinArgs() int                         { return 1 }
func (fakeFunction) MaxArgs() int                         { return 1 }
func (fakeFunction) ValidateArgTypes(args []Value) error  { return nil }
func (fakeFunction) Evaluate(args []Value) (Value, error) { return args[0], nil }

func TestRegisterAndLookup(t *testing.T) {
	Register(fakeFunction{})

	fn, err := Lookup("FAKE")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	if fn.Name() != "FAKE" {
		t.Errorf("Name() = %q, want FAKE", fn.Name())
	}
	got, err := fn.Evaluate([]Value{42})
	if err != nil || got != 42 {
		t.Errorf("Evaluate() = (%v, %v), want (42, nil)", got, err)
	}
}

func TestLookup_UnknownFunction(t *testing.T) {
	_, err := Lookup("DOES_NOT_EXIST")
	if err == nil {
		t.Fatal("Lookup() error = nil, want error for unknown function")
	}
}
