package parser

type NodeType int

const (
	NodeNumber NodeType = iota
	NodeString
	NodeIdentifier
	NodeBinary
)

// Node is one node of the formula's abstract syntax tree. For
// NodeNumber/NodeString/NodeIdentifier, Value holds the literal text
// or name. For NodeBinary, Operator/Left/Right describe the operation.
type Node struct {
	Type     NodeType
	Value    string
	Operator TokenType
	Left     *Node
	Right    *Node
}
