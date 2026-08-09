package dependency

import "github.com/ferivision/formula-engine/internal/parser"

// FormulaSource pairs a formula name with its parsed AST. It mirrors
// the top-level package's FormulaInput without importing it, avoiding
// an import cycle (the top-level package will import dependency, not
// the other way around).
type FormulaSource struct {
	Name string
	Node *parser.Node
}

// Graph is the dependency graph across one batch of formulas. Nodes
// holds every formula name in original input order; Edges maps a
// formula name to the names of other formulas in the batch it
// depends on (not plain data fields).
type Graph struct {
	Nodes []string
	Edges map[string][]string
}

// BuildGraph builds the dependency graph for formulas, given the set
// of plain data field names available in the call. A reference that
// is neither another formula's name nor a known data field is
// reported in the returned map, keyed by the referencing formula's
// name -- this feeds FR-6's undefined-reference detection.
func BuildGraph(formulas []FormulaSource, dataFields map[string]bool) (*Graph, map[string][]string) {
	formulaNames := make(map[string]bool, len(formulas))
	for _, f := range formulas {
		formulaNames[f.Name] = true
	}

	graph := &Graph{
		Nodes: make([]string, 0, len(formulas)),
		Edges: make(map[string][]string, len(formulas)),
	}
	undefined := make(map[string][]string)

	for _, f := range formulas {
		graph.Nodes = append(graph.Nodes, f.Name)
		graph.Edges[f.Name] = nil

		for _, ref := range ExtractReferences(f.Node) {
			switch {
			case formulaNames[ref]:
				graph.Edges[f.Name] = append(graph.Edges[f.Name], ref)
			case dataFields[ref]:
				// Plain data field: resolved directly, not a graph edge.
			default:
				undefined[f.Name] = append(undefined[f.Name], ref)
			}
		}
	}

	return graph, undefined
}
