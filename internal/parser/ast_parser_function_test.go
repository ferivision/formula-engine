package parser

import "testing"

func TestParse_FunctionCallNoArgs(t *testing.T) {
	node := mustParse(t, "NOW()")

	if node.Type != NodeFunctionCall || node.Value != "NOW" {
		t.Fatalf("node = %+v, want NodeFunctionCall(NOW)", node)
	}
	if len(node.Args) != 0 {
		t.Errorf("Args = %v, want empty", node.Args)
	}
}

func TestParse_FunctionCallOneArg(t *testing.T) {
	node := mustParse(t, "ABS(1)")

	if node.Type != NodeFunctionCall || node.Value != "ABS" {
		t.Fatalf("node = %+v, want NodeFunctionCall(ABS)", node)
	}
	if len(node.Args) != 1 || node.Args[0].Value != "1" {
		t.Fatalf("Args = %+v, want [NodeNumber(1)]", node.Args)
	}
}

func TestParse_FunctionCallManyArgs(t *testing.T) {
	node := mustParse(t, "SUM(1, 2, 3)")

	if node.Type != NodeFunctionCall || node.Value != "SUM" {
		t.Fatalf("node = %+v, want NodeFunctionCall(SUM)", node)
	}
	if len(node.Args) != 3 {
		t.Fatalf("Args = %+v, want 3 args", node.Args)
	}
	for i, want := range []string{"1", "2", "3"} {
		if node.Args[i].Value != want {
			t.Errorf("Args[%d] = %+v, want NodeNumber(%s)", i, node.Args[i], want)
		}
	}
}

func TestParse_NestedFunctionCalls(t *testing.T) {
	node := mustParse(t, "MAX(1, MIN(2, 3))")

	if node.Type != NodeFunctionCall || node.Value != "MAX" {
		t.Fatalf("node = %+v, want NodeFunctionCall(MAX)", node)
	}
	if len(node.Args) != 2 {
		t.Fatalf("Args = %+v, want 2 args", node.Args)
	}
	if node.Args[0].Type != NodeNumber || node.Args[0].Value != "1" {
		t.Errorf("Args[0] = %+v, want NodeNumber(1)", node.Args[0])
	}
	inner := node.Args[1]
	if inner.Type != NodeFunctionCall || inner.Value != "MIN" || len(inner.Args) != 2 {
		t.Fatalf("Args[1] = %+v, want NodeFunctionCall(MIN) with 2 args", inner)
	}
	if inner.Args[0].Value != "2" || inner.Args[1].Value != "3" {
		t.Errorf("MIN args = %+v / %+v, want 2 and 3", inner.Args[0], inner.Args[1])
	}
}

func TestParse_FunctionCallWithinArithmetic(t *testing.T) {
	// MAX(1, 2) * 3 must parse as (MAX(1, 2)) * 3.
	node := mustParse(t, "MAX(1, 2) * 3")

	if node.Type != NodeBinary || node.Operator != TokenStar {
		t.Fatalf("root = %+v, want top-level *", node)
	}
	if node.Left.Type != NodeFunctionCall || node.Left.Value != "MAX" {
		t.Fatalf("left = %+v, want NodeFunctionCall(MAX)", node.Left)
	}
	if node.Right.Value != "3" {
		t.Errorf("right = %+v, want NodeNumber(3)", node.Right)
	}
}
