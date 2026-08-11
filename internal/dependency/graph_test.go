package dependency

import (
	"sort"
	"testing"

	"github.com/ferivision/formula-engine/internal/parser"
)

func source(t *testing.T, name, expr string) FormulaSource {
	t.Helper()
	tokens, err := parser.Tokenize(expr)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", expr, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", expr, err)
	}
	return FormulaSource{Name: name, Node: node}
}

func TestBuildGraph_ChainedFormulas(t *testing.T) {
	formulas := []FormulaSource{
		source(t, "subtotal", "price * quantity"),
		source(t, "discount", "subtotal * 0.1"),
		source(t, "total", "subtotal - discount"),
	}
	dataFields := map[string]bool{"price": true, "quantity": true}

	graph, undefined := BuildGraph(formulas, dataFields)

	if len(undefined) != 0 {
		t.Fatalf("undefined = %v, want none", undefined)
	}

	assertDeps(t, graph, "subtotal", []string{})
	assertDeps(t, graph, "discount", []string{"subtotal"})
	assertDeps(t, graph, "total", []string{"subtotal", "discount"})
}

func TestBuildGraph_UndefinedReferenceIsFlagged(t *testing.T) {
	formulas := []FormulaSource{
		source(t, "total", "subtotal + unknown_field"),
	}
	dataFields := map[string]bool{"subtotal": false} // subtotal is neither a data field nor another formula here

	graph, undefined := BuildGraph(formulas, dataFields)

	if len(graph.Edges["total"]) != 0 {
		t.Errorf("Edges[total] = %v, want no formula dependencies", graph.Edges["total"])
	}
	got := undefined["total"]
	sort.Strings(got)
	want := []string{"subtotal", "unknown_field"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("undefined[total] = %v, want %v", got, want)
	}
}

func TestBuildGraph_PlainDataFieldIsNotAnEdge(t *testing.T) {
	formulas := []FormulaSource{
		source(t, "total", "price * 2"),
	}
	dataFields := map[string]bool{"price": true}

	graph, undefined := BuildGraph(formulas, dataFields)

	if len(graph.Edges["total"]) != 0 {
		t.Errorf("Edges[total] = %v, want no formula dependencies (price is a data field)", graph.Edges["total"])
	}
	if len(undefined) != 0 {
		t.Errorf("undefined = %v, want none", undefined)
	}
}

func assertDeps(t *testing.T, graph *Graph, name string, want []string) {
	t.Helper()
	got := append([]string{}, graph.Edges[name]...)
	sort.Strings(got)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("Edges[%s] = %v, want %v", name, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Edges[%s] = %v, want %v", name, got, want)
		}
	}
}
