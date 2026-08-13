package formulaengine

import (
	"github.com/ferivision/formula-engine/internal/apperror"
	"github.com/ferivision/formula-engine/internal/dependency"
	"github.com/ferivision/formula-engine/internal/evaluator"
	"github.com/ferivision/formula-engine/internal/parser"

	_ "github.com/ferivision/formula-engine/internal/registry/aggregate"  // registers SUMIF, COUNTIF
	_ "github.com/ferivision/formula-engine/internal/registry/comparison" // registers EQUALS, BETWEEN
	_ "github.com/ferivision/formula-engine/internal/registry/date"       // registers NOW, DATE_ADD
	_ "github.com/ferivision/formula-engine/internal/registry/logic"      // registers AND, OR, NOT
	_ "github.com/ferivision/formula-engine/internal/registry/math"       // registers MAX, MIN, SUM, AVG, ROUND, FLOOR, CEIL, ABS
	_ "github.com/ferivision/formula-engine/internal/registry/text"       // registers CONCAT, UPPER, LOWER, TRIM, LENGTH
)

// FormulaInput is one named formula to evaluate.
type FormulaInput struct {
	Name       string // identifies this formula in the output map
	Expression string // e.g. "MAX(subtotal - discount, 0)"
}

// Result is the outcome of evaluating a single formula.
type Result struct {
	Value any   // the computed value, nil if Err is set
	Err   error // per-formula error, nil on success
}

// Evaluate computes one or more formulas against the given data and
// returns a result per formula name, or an error if the call as a
// whole cannot proceed (e.g. a circular reference across the inputs).
//
// Per rfc.md §7/§11: a circular or undefined reference is a
// call-level error -- returned here as the second value, with a nil
// results map and no partial evaluation. A syntax error in one
// formula's expression is formula-level -- it's attached to that
// formula's Result.Err, and the rest of the batch still proceeds.
func Evaluate(formulas []FormulaInput, data map[string]any) (map[string]Result, error) {
	dataFields := make(map[string]bool, len(data))
	for k := range data {
		dataFields[k] = true
	}

	results := make(map[string]Result, len(formulas))
	nodes := make(map[string]*parser.Node, len(formulas))
	sources := make([]dependency.FormulaSource, 0, len(formulas))

	for _, f := range formulas {
		tokens, err := parser.Tokenize(f.Expression)
		if err != nil {
			results[f.Name] = Result{Err: err}
			continue
		}
		node, err := parser.Parse(tokens)
		if err != nil {
			results[f.Name] = Result{Err: err}
			continue
		}
		nodes[f.Name] = node
		sources = append(sources, dependency.FormulaSource{Name: f.Name, Node: node})
	}

	graph, undefined := dependency.BuildGraph(sources, dataFields)
	if err := firstUndefinedError(sources, undefined); err != nil {
		return nil, err
	}
	if dependency.HasCycle(graph) {
		return nil, &apperror.FormulaError{
			Code:    apperror.ErrCircularReference,
			Message: "circular reference detected among the given formulas",
		}
	}

	computed := make(map[string]any, len(sources))
	for _, name := range dependency.TopologicalSort(graph) {
		if err := firstFailedDependency(name, graph, results); err != nil {
			results[name] = Result{Err: err}
			continue
		}

		ctxData := make(map[string]any, len(data)+len(computed))
		for k, v := range data {
			ctxData[k] = v
		}
		for k, v := range computed {
			ctxData[k] = v
		}

		value, err := evaluator.Evaluate(nodes[name], evaluator.NewContext(ctxData))
		if err != nil {
			results[name] = Result{Err: err}
			continue
		}
		results[name] = Result{Value: value}
		computed[name] = value
	}

	return results, nil
}

// firstFailedDependency reports whether any formula name directly
// depends on (per graph.Edges) has already failed, so name can
// inherit that failure with a clear "dependency failed" message
// instead of being evaluated with the dependency silently missing
// from context (where it would otherwise surface as a misleading
// ErrUndefinedReference -- the dependency IS defined, it just failed
// at runtime). Because dependencies are evaluated in topological
// order before name, an inherited failure here already reflects any
// failure the dependency itself inherited, so this naturally cascades
// transitively without extra work.
func firstFailedDependency(name string, graph *dependency.Graph, results map[string]Result) error {
	for _, dep := range graph.Edges[name] {
		if r, ok := results[dep]; ok && r.Err != nil {
			return &apperror.FormulaError{
				Code:    apperror.ErrRuntime,
				Message: "depends on formula \"" + dep + "\", which failed: " + r.Err.Error(),
			}
		}
	}
	return nil
}

// firstUndefinedError picks the first undefined reference in sources'
// original order, so the reported error is deterministic (NFR-2)
// rather than dependent on Go's randomized map iteration order.
func firstUndefinedError(sources []dependency.FormulaSource, undefined map[string][]string) error {
	for _, s := range sources {
		refs, ok := undefined[s.Name]
		if !ok || len(refs) == 0 {
			continue
		}
		return &apperror.FormulaError{
			Code:    apperror.ErrUndefinedReference,
			Message: "formula \"" + s.Name + "\" references undefined field/formula \"" + refs[0] + "\"",
		}
	}
	return nil
}
