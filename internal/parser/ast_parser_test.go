package parser

import "testing"

func mustParse(t *testing.T, input string) *Node {
	t.Helper()
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", input, err)
	}
	node, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", input, err)
	}
	return node
}

func TestParse_NumberLiteral(t *testing.T) {
	node := mustParse(t, "42")
	if node.Type != NodeNumber || node.Value != "42" {
		t.Errorf("node = %+v, want NodeNumber(42)", node)
	}
}

func TestParse_StringLiteral(t *testing.T) {
	node := mustParse(t, `"hello"`)
	if node.Type != NodeString || node.Value != "hello" {
		t.Errorf("node = %+v, want NodeString(hello)", node)
	}
}

func TestParse_Identifier(t *testing.T) {
	node := mustParse(t, "subtotal")
	if node.Type != NodeIdentifier || node.Value != "subtotal" {
		t.Errorf("node = %+v, want NodeIdentifier(subtotal)", node)
	}
}

func TestParse_MultiplicationBindsTighterThanAddition(t *testing.T) {
	// 1 + 2 * 3 must parse as 1 + (2 * 3), not (1 + 2) * 3.
	node := mustParse(t, "1 + 2 * 3")

	if node.Type != NodeBinary || node.Operator != TokenPlus {
		t.Fatalf("root = %+v, want top-level +", node)
	}
	if node.Left.Type != NodeNumber || node.Left.Value != "1" {
		t.Errorf("left = %+v, want NodeNumber(1)", node.Left)
	}
	if node.Right.Type != NodeBinary || node.Right.Operator != TokenStar {
		t.Fatalf("right = %+v, want nested *", node.Right)
	}
	if node.Right.Left.Value != "2" || node.Right.Right.Value != "3" {
		t.Errorf("right operands = %+v / %+v, want 2 and 3", node.Right.Left, node.Right.Right)
	}
}

func TestParse_ParenthesesOverridePrecedence(t *testing.T) {
	// (1 + 2) * 3 must parse as (1 + 2) * 3, not 1 + (2 * 3).
	node := mustParse(t, "(1 + 2) * 3")

	if node.Type != NodeBinary || node.Operator != TokenStar {
		t.Fatalf("root = %+v, want top-level *", node)
	}
	if node.Left.Type != NodeBinary || node.Left.Operator != TokenPlus {
		t.Fatalf("left = %+v, want nested +", node.Left)
	}
	if node.Right.Value != "3" {
		t.Errorf("right = %+v, want NodeNumber(3)", node.Right)
	}
}

func TestParse_SubtractionIsLeftAssociative(t *testing.T) {
	// 1 - 2 - 3 must parse as (1 - 2) - 3, not 1 - (2 - 3).
	node := mustParse(t, "1 - 2 - 3")

	if node.Type != NodeBinary || node.Operator != TokenMinus {
		t.Fatalf("root = %+v, want top-level -", node)
	}
	if node.Left.Type != NodeBinary || node.Left.Operator != TokenMinus {
		t.Fatalf("left = %+v, want nested (1 - 2)", node.Left)
	}
	if node.Left.Left.Value != "1" || node.Left.Right.Value != "2" {
		t.Errorf("nested operands = %+v / %+v, want 1 and 2", node.Left.Left, node.Left.Right)
	}
	if node.Right.Value != "3" {
		t.Errorf("right = %+v, want NodeNumber(3)", node.Right)
	}
}
