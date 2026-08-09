package dependency

const (
	unvisited = iota
	visiting
	done
)

type frame struct {
	node    string
	leaving bool
}

// HasCycle reports whether graph contains a cycle, using an iterative
// depth-first search (explicit stack, three-color visited tracking)
// -- never recursion, per rfc.md §5 / CLAUDE.md. Correctness holds
// regardless of cycle length (PRD use case 5), since the stack is
// heap-allocated, not bounded by Go's call stack.
func HasCycle(graph *Graph) bool {
	state := make(map[string]int, len(graph.Nodes))

	for _, start := range graph.Nodes {
		if state[start] != unvisited {
			continue
		}

		stack := []frame{{node: start}}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if top.leaving {
				state[top.node] = done
				continue
			}
			if state[top.node] != unvisited {
				continue // already visiting or done via another path
			}
			state[top.node] = visiting
			stack = append(stack, frame{node: top.node, leaving: true})

			for _, dep := range graph.Edges[top.node] {
				switch state[dep] {
				case visiting:
					return true // back edge -> cycle
				case unvisited:
					stack = append(stack, frame{node: dep})
				}
			}
		}
	}

	return false
}
