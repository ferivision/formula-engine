package dependency

// TopologicalSort returns graph's nodes in a valid dependency-
// respecting order (a formula's dependencies before the formula
// itself), using iterative Kahn's algorithm -- an explicit queue,
// never recursion, per rfc.md §5 / CLAUDE.md.
//
// If graph contains a cycle, the nodes involved in it never reach
// in-degree 0 and are simply absent from the result; this function
// does not itself detect or reject cycles -- see FASE-8.1/8.2.
func TopologicalSort(graph *Graph) []string {
	inDegree := make(map[string]int, len(graph.Nodes))
	dependents := make(map[string][]string, len(graph.Nodes))

	for _, n := range graph.Nodes {
		inDegree[n] = 0
	}
	for name, deps := range graph.Edges {
		inDegree[name] = len(deps)
		for _, dep := range deps {
			dependents[dep] = append(dependents[dep], name)
		}
	}

	queue := make([]string, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}

	order := make([]string, 0, len(graph.Nodes))
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)

		for _, dependent := range dependents[n] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	return order
}
