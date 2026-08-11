package parser

type NodeType int

const (
	NodeNumber NodeType = iota
	NodeString
	NodeIdentifier
	NodeBinary
	NodeFunctionCall
)

// Node is one node of the formula's abstract syntax tree. For
// NodeNumber/NodeString/NodeIdentifier, Value holds the literal text
// or name. For NodeBinary, Operator/Left/Right describe the operation.
// For NodeFunctionCall, Value holds the function name and Args holds
// its (possibly zero) arguments in order.
type Node struct {
	Type     NodeType
	Value    string
	Operator TokenType
	Left     *Node
	Right    *Node
	Args     []*Node
}
