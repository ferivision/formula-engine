package dependency

import "github.com/ferivision/formula-engine/internal/parser"

// ExtractReferences returns the distinct field/formula names node
// references, found via an explicit-stack walk of the AST -- never
// recursive, per rfc.md §5 / CLAUDE.md. The function name in a
// NodeFunctionCall is not itself a reference; only its arguments are
// walked.
func ExtractReferences(node *parser.Node) []string {
	if node == nil {
		return nil
	}

	seen := make(map[string]bool)
	var refs []string
	stack := []*parser.Node{node}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil {
			continue
		}

		switch n.Type {
		case parser.NodeIdentifier:
			if !seen[n.Value] {
				seen[n.Value] = true
				refs = append(refs, n.Value)
			}
		case parser.NodeBinary:
			stack = append(stack, n.Left, n.Right)
		case parser.NodeFunctionCall:
			stack = append(stack, n.Args...)
		}
	}

	return refs
}
