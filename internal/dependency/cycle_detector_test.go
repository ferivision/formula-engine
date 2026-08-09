package dependency

import (
	"fmt"
	"testing"
)

func TestHasCycle_DirectCycle(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"a", "b"},
		Edges: map[string][]string{"a": {"b"}, "b": {"a"}},
	}
	if !HasCycle(graph) {
		t.Error("HasCycle() = false, want true for a direct A<->B cycle")
	}
}

func TestHasCycle_AcyclicGraphReportsNoCycle(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"subtotal", "discount", "total"},
		Edges: map[string][]string{
			"subtotal": nil,
			"discount": {"subtotal"},
			"total":    {"subtotal", "discount"},
		},
	}
	if HasCycle(graph) {
		t.Error("HasCycle() = true, want false for an acyclic dependency chain")
	}
}

func TestHasCycle_DisjointAcyclicComponents(t *testing.T) {
	graph := &Graph{
		Nodes: []string{"a", "b", "c"},
		Edges: map[string][]string{"a": nil, "b": {"a"}, "c": nil},
	}
	if HasCycle(graph) {
		t.Error("HasCycle() = true, want false")
	}
}

// PRD use case 5 requires correctness "no matter how long the cycle
// is" -- this proves it at 1000+ nodes without recursion-depth
// concerns, since the detector is iterative (rfc.md §5 / CLAUDE.md).
func TestHasCycle_LongIndirectCycle(t *testing.T) {
	const n = 1000
	nodes := make([]string, n)
	edges := make(map[string][]string, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("n%d", i)
		next := fmt.Sprintf("n%d", (i+1)%n) // n999 -> n0 closes the cycle
		nodes[i] = name
		edges[name] = []string{next}
	}
	graph := &Graph{Nodes: nodes, Edges: edges}

	if !HasCycle(graph) {
		t.Error("HasCycle() = false, want true for a 1000-node indirect cycle")
	}
}

func TestHasCycle_LongAcyclicChain(t *testing.T) {
	const n = 1000
	nodes := make([]string, n)
	edges := make(map[string][]string, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("n%d", i)
		nodes[i] = name
		if i > 0 {
			edges[name] = []string{fmt.Sprintf("n%d", i-1)}
		}
	}
	graph := &Graph{Nodes: nodes, Edges: edges}

	if HasCycle(graph) {
		t.Error("HasCycle() = true, want false for a 1000-node acyclic chain")
	}
}
