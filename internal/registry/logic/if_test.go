package logic

import "testing"

func TestIf_Name(t *testing.T) {
	if (ifFunction{}).Name() != "IF" {
		t.Errorf("Name() = %q, want IF", (ifFunction{}).Name())
	}
}

func TestIf_ArgBounds(t *testing.T) {
	f := ifFunction{}
	if f.MinArgs() != 3 || f.MaxArgs() != 3 {
		t.Errorf("bounds = (%d, %d), want (3, 3)", f.MinArgs(), f.MaxArgs())
	}
}

func TestIf_EvaluateIsNeverMeantToBeCalled(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Evaluate() did not panic; internal/evaluator must be the only caller of IF's evaluation")
		}
	}()
	_, _ = (ifFunction{}).Evaluate(nil)
}
