package evaluator

import (
	"strconv"

	"github.com/ferivision/formula-engine/internal/parser"
	"github.com/ferivision/formula-engine/internal/registry"
)

// Evaluate walks an AST produced by parser.Parse and computes its
// value, resolving identifier nodes against ctx.
func Evaluate(node *parser.Node, ctx *Context) (registry.Value, error) {
	switch node.Type {
	case parser.NodeNumber:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, newRuntimeError("invalid number literal: " + node.Value)
		}
		return v, nil
	case parser.NodeString:
		return node.Value, nil
	case parser.NodeBinary:
		return evalBinary(node, ctx)
	case parser.NodeIdentifier:
		v, ok := ctx.Lookup(node.Value)
		if !ok {
			return nil, newUndefinedReferenceError(node.Value)
		}
		return v, nil
	case parser.NodeFunctionCall:
		return evalFunctionCall(node, ctx)
	default:
		return nil, newRuntimeError("unknown node type")
	}
}

func evalBinary(node *parser.Node, ctx *Context) (registry.Value, error) {
	leftVal, err := Evaluate(node.Left, ctx)
	if err != nil {
		return nil, err
	}
	rightVal, err := Evaluate(node.Right, ctx)
	if err != nil {
		return nil, err
	}
	left, err := toFloat64(leftVal)
	if err != nil {
		return nil, err
	}
	right, err := toFloat64(rightVal)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case parser.TokenPlus:
		return left + right, nil
	case parser.TokenMinus:
		return left - right, nil
	case parser.TokenStar:
		return left * right, nil
	case parser.TokenSlash:
		if right == 0 {
			return nil, newRuntimeError("division by zero")
		}
		return left / right, nil
	default:
		return nil, newRuntimeError("unknown operator")
	}
}

func evalFunctionCall(node *parser.Node, ctx *Context) (registry.Value, error) {
	fn, err := registry.Lookup(node.Value)
	if err != nil {
		return nil, err
	}

	args := make([]registry.Value, len(node.Args))
	for i, argNode := range node.Args {
		v, err := Evaluate(argNode, ctx)
		if err != nil {
			return nil, err
		}
		args[i] = v
	}

	return fn.Evaluate(args)
}

func toFloat64(v registry.Value) (float64, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, newTypeError("expected a number")
	}
	return f, nil
}
