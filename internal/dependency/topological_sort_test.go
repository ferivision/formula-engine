package dependency

import "testing"

func indexOf(order []string, name string) int {
	for i, n := range order {
		if n == name {
			return i
		}
	}
	return -1
}

func TestTopologicalSort_ChainRegardlessOfInputOrder(t *testing.T) {
	edges := map[string][]string{
		"total":    {"subtotal", "discount"},
		"discount": {"subtotal"},
		"subtotal": nil,
	}

	for _, nodes := range [][]string{
		{"subtotal", "discount", "total"},
		{"total", "discount", "subtotal"},
		{"discount", "total", "subtotal"},
	} {
		graph := &Graph{Nodes: nodes, Edges: edges}
		order := TopologicalSort(graph)

		if len(order) != 3 {
			t.Fatalf("input order %v: TopologicalSort() = %v, want all 3 nodes", nodes, order)
		}
		if indexOf(order, "subtotal") >= indexOf(order, "discount") {
			t.Errorf("input order %v: order = %v, want subtotal before discount", nodes, order)
		}
		if indexOf(order, "discount") >= indexOf(order, "total") {
			t.Errorf("input order %v: order = %v, want discount before total", nodes, order)
		}
	}
}

func TestTopologicalSort_IndependentNodesAllIncluded(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"a", "b", "c"},
		Edges: map[string][]string{"a": nil, "b": nil, "c": nil},
	}

	order := TopologicalSort(graph)
	if len(order) != 3 {
		t.Fatalf("TopologicalSort() = %v, want all 3 independent nodes", order)
	}
}

// A cycle can never be fully resolved by Kahn's algorithm -- the
// nodes involved never reach in-degree 0, so they're simply absent
// from the result. This ticket only requires that boundary behave
// safely (no panic, no infinite loop); rejecting the call outright
// with a CircularReferenceError is Phase 8's job (FASE-8.1/8.2).
func TestTopologicalSort_CycleIsOmittedNotPanicking(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"a", "b"},
		Edges: map[string][]string{"a": {"b"}, "b": {"a"}},
	}

	order := TopologicalSort(graph)
	if len(order) != 0 {
		t.Errorf("TopologicalSort() = %v, want empty (both nodes are in the cycle)", order)
	}
}

func TestTopologicalSort_CycleDoesNotBlockUnrelatedNodes(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"a", "b", "independent"},
		Edges: map[string][]string{"a": {"b"}, "b": {"a"}, "independent": nil},
	}

	order := TopologicalSort(graph)
	if len(order) != 1 || order[0] != "independent" {
		t.Errorf("TopologicalSort() = %v, want [independent] only", order)
	}
}
