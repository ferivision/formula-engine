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

	// string + string concatenates (rfc.md §10); every other operand
	// combination -- including number + string -- goes through
	// numeric coercion below.
	if node.Operator == parser.TokenPlus {
		if leftStr, ok := leftVal.(string); ok {
			if rightStr, ok := rightVal.(string); ok {
				return leftStr + rightStr, nil
			}
		}
	}

	left, err := coerceForArithmetic(leftVal)
	if err != nil {
		return nil, err
	}
	right, err := coerceForArithmetic(rightVal)
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
	// IF requires short-circuit evaluation (only the taken branch
	// runs) and so can't go through the generic eager-evaluate-all-
	// args path below -- see internal/registry/logic.ifFunction.
	if node.Value == "IF" {
		return evalIf(node, ctx)
	}

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

func evalIf(node *parser.Node, ctx *Context) (registry.Value, error) {
	fn, err := registry.Lookup("IF")
	if err != nil {
		return nil, err
	}
	if len(node.Args) < fn.MinArgs() || len(node.Args) > fn.MaxArgs() {
		return nil, newRuntimeError("IF: wrong number of arguments (got " + strconv.Itoa(len(node.Args)) + ", want 3)")
	}

	condVal, err := Evaluate(node.Args[0], ctx)
	if err != nil {
		return nil, err
	}
	cond, err := coerceForBool(condVal)
	if err != nil {
		return nil, err
	}
	if cond {
		return Evaluate(node.Args[1], ctx)
	}
	return Evaluate(node.Args[2], ctx)
}
